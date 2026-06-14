package com.demoapp.systemtest.model;

/** The subset of the OIDC token endpoint response the tests rely on. */
public record TokenResponse(String accessToken) {
}
