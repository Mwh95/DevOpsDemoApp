package http

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// KeycloakJWKSVerifier verifies JWT access tokens using Keycloak's JWKS endpoint.
type KeycloakJWKSVerifier struct {
	issuer    string
	jwksURL   string
	keys      map[string]*rsa.PublicKey
	keysMu    sync.RWMutex
	client    *http.Client
	lastFetch time.Time
}

// NewKeycloakJWKSVerifier creates a verifier that fetches JWKS from the given issuer (e.g. https://keycloak/login/realms/myrealm).
// It derives the JWKS URL from issuer by appending /.well-known/openid-configuration and then using jwks_uri.
func NewKeycloakJWKSVerifier(issuer string, client *http.Client) (*KeycloakJWKSVerifier, error) {
	if client == nil {
		client = http.DefaultClient
	}
	v := &KeycloakJWKSVerifier{issuer: strings.TrimSuffix(issuer, "/"), client: client, keys: make(map[string]*rsa.PublicKey)}
	if err := v.refreshKeys(); err != nil {
		return nil, err
	}
	return v, nil
}

func (v *KeycloakJWKSVerifier) refreshKeys() error {
	discoveryURL := strings.TrimSuffix(v.issuer, "/") + "/.well-known/openid-configuration"
	resp, err := v.client.Get(discoveryURL)
	if err != nil {
		return fmt.Errorf("fetch oidc config: %w", err)
	}
	defer resp.Body.Close()
	var discovery struct {
		JWKSURI string `json:"jwks_uri"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return fmt.Errorf("decode oidc config: %w", err)
	}
	if discovery.JWKSURI == "" {
		return fmt.Errorf("missing jwks_uri in discovery")
	}
	v.jwksURL = discovery.JWKSURI
	resp2, err := v.client.Get(v.jwksURL)
	if err != nil {
		return fmt.Errorf("fetch jwks: %w", err)
	}
	defer resp2.Body.Close()
	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("decode jwks: %w", err)
	}
	v.keysMu.Lock()
	defer v.keysMu.Unlock()
	v.keys = make(map[string]*rsa.PublicKey)
	for _, k := range jwks.Keys {
		if k.Kty != "RSA" || (k.Use != "" && k.Use != "sig") {
			continue
		}
		n, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			continue
		}
		var e int
		if len(eBytes) >= 4 {
			e = int(binary.BigEndian.Uint32(eBytes))
		} else {
			for _, b := range eBytes {
				e = e<<8 | int(b)
			}
		}
		v.keys[k.Kid] = &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: e}
	}
	v.lastFetch = time.Now()
	return nil
}

// VerifyAndExtract verifies the JWT and returns the subject (user ID). Token should be the raw Bearer token string.
func (v *KeycloakJWKSVerifier) VerifyAndExtract(ctx context.Context, tokenString string) (subject string, err error) {
	v.keysMu.RLock()
	if time.Since(v.lastFetch) > 5*time.Minute {
		v.keysMu.RUnlock()
		v.keysMu.Lock()
		if time.Since(v.lastFetch) > 5*time.Minute {
			_ = v.refreshKeys()
		}
		v.keysMu.Unlock()
		v.keysMu.RLock()
	}
	defer v.keysMu.RUnlock()

	tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in header")
		}
		pub, ok := v.keys[kid]
		if !ok {
			return nil, fmt.Errorf("unknown key id %q", kid)
		}
		return pub, nil
	})
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}
	if !tok.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	iss, _ := claims["iss"].(string)
	if iss != v.issuer {
		return "", fmt.Errorf("invalid issuer %q", iss)
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", fmt.Errorf("missing sub")
	}
	return sub, nil
}

// TokenVerifier verifies a bearer token and returns the subject (user ID). Used for testing with a fake.
type TokenVerifier interface {
	VerifyAndExtract(ctx context.Context, token string) (subject string, err error)
}

// RequireAuth returns middleware that validates Bearer token and sets user ID in context.
func RequireAuth(v TokenVerifier) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"error":"missing or invalid authorization"}`, http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			sub, err := v.VerifyAndExtract(r.Context(), token)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := WithUserID(r.Context(), sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth (method) delegates to the package-level RequireAuth function.
func (v *KeycloakJWKSVerifier) RequireAuth(next http.Handler) http.Handler {
	return RequireAuth(v)(next)
}
