package com.demoapp.systemtest.support;

import com.demoapp.systemtest.env.TestEnvironment;
import com.fasterxml.jackson.databind.JsonNode;
import com.microsoft.playwright.BrowserContext;
import com.microsoft.playwright.Page;

/**
 * Per-scenario shared state. A fresh instance is created by PicoContainer for every scenario and
 * injected into the hooks and step definitions that declare it as a constructor argument.
 */
public class World {

    public final TestEnvironment env = TestEnvironment.get();
    public final ApiClient api = new ApiClient(env.baseUrl(), env.tokenEndpoint());

    // UI state
    public BrowserContext browserContext;
    public Page page;

    // API state
    public String token;
    public ApiResponse lastResponse;
    public JsonNode createdMarker;

    public String baseUrl() {
        return env.baseUrl();
    }

    public String createdMarkerId() {
        if (createdMarker == null || createdMarker.get("id") == null) {
            throw new IllegalStateException("No marker has been created in this scenario yet");
        }
        return createdMarker.get("id").asText();
    }
}
