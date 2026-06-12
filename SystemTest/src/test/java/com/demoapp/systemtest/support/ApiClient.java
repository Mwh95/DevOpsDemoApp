package com.demoapp.systemtest.support;

import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpRequest.BodyPublishers;
import java.net.http.HttpResponse;
import java.net.http.HttpResponse.BodyHandlers;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.LinkedHashMap;
import java.util.Map;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

/** Thin REST client used by the {@code @api} step definitions and the UI clean-up hook. */
public class ApiClient {

    private final String baseUrl;
    private final String tokenEndpoint;
    private final HttpClient http;
    private final ObjectMapper mapper = new ObjectMapper();

    public ApiClient(String baseUrl, String tokenEndpoint) {
        this.baseUrl = baseUrl;
        this.tokenEndpoint = tokenEndpoint;
        this.http = HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(20)).build();
    }

    /** Fetches an access token via the test realm's direct access (password) grant. */
    public String passwordGrantToken(String username, String password) {
        Map<String, String> form = new LinkedHashMap<>();
        form.put("grant_type", "password");
        form.put("client_id", "map-app");
        form.put("username", username);
        form.put("password", password);
        form.put("scope", "openid profile");

        HttpRequest request = HttpRequest.newBuilder(URI.create(tokenEndpoint))
                .header("Content-Type", "application/x-www-form-urlencoded")
                .POST(BodyPublishers.ofString(urlEncode(form)))
                .build();
        HttpResponse<String> response = send(request);
        if (response.statusCode() != 200) {
            throw new IllegalStateException(
                    "Token request failed (" + response.statusCode() + "): " + response.body());
        }
        try {
            return mapper.readTree(response.body()).get("access_token").asText();
        } catch (Exception e) {
            throw new IllegalStateException("Could not parse token response", e);
        }
    }

    public ApiResponse get(String path, String token) {
        return exchange(builder(path, token).GET().build());
    }

    public ApiResponse getNoAuth(String path) {
        return exchange(HttpRequest.newBuilder(URI.create(baseUrl + path)).GET().build());
    }

    public ApiResponse post(String path, String token, String jsonBody) {
        return exchange(builder(path, token)
                .header("Content-Type", "application/json")
                .POST(BodyPublishers.ofString(jsonBody))
                .build());
    }

    public ApiResponse put(String path, String token, String jsonBody) {
        return exchange(builder(path, token)
                .header("Content-Type", "application/json")
                .PUT(BodyPublishers.ofString(jsonBody))
                .build());
    }

    public ApiResponse delete(String path, String token) {
        return exchange(builder(path, token).DELETE().build());
    }

    private HttpRequest.Builder builder(String path, String token) {
        HttpRequest.Builder builder = HttpRequest.newBuilder(URI.create(baseUrl + path));
        if (token != null) {
            builder.header("Authorization", "Bearer " + token);
        }
        return builder;
    }

    private ApiResponse exchange(HttpRequest request) {
        HttpResponse<String> response = send(request);
        JsonNode json = null;
        String body = response.body();
        if (body != null && !body.isBlank()) {
            try {
                json = mapper.readTree(body);
            } catch (Exception ignored) {
                // Non-JSON body (e.g. plain-text auth errors); leave json null.
            }
        }
        return new ApiResponse(response.statusCode(), body, json);
    }

    private HttpResponse<String> send(HttpRequest request) {
        try {
            return http.send(request, BodyHandlers.ofString(StandardCharsets.UTF_8));
        } catch (Exception e) {
            throw new IllegalStateException("HTTP request failed: " + request.uri(), e);
        }
    }

    private static String urlEncode(Map<String, String> form) {
        StringBuilder sb = new StringBuilder();
        for (Map.Entry<String, String> entry : form.entrySet()) {
            if (sb.length() > 0) {
                sb.append('&');
            }
            sb.append(URLEncoder.encode(entry.getKey(), StandardCharsets.UTF_8))
                    .append('=')
                    .append(URLEncoder.encode(entry.getValue(), StandardCharsets.UTF_8));
        }
        return sb.toString();
    }
}
