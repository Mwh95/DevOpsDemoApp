package com.demoapp.systemtest.support;

import com.fasterxml.jackson.databind.JsonNode;

/** Lightweight container for an HTTP response used by the API step definitions. */
public record ApiResponse(int status, String body, JsonNode json) {

    public String fieldText(String name) {
        if (json == null || json.get(name) == null) {
            return null;
        }
        return json.get(name).asText();
    }
}
