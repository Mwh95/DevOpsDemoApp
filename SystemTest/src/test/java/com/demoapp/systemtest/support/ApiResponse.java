package com.demoapp.systemtest.support;

import java.util.List;
import java.util.Map;

import com.fasterxml.jackson.databind.JavaType;

/**
 * HTTP response wrapper. The raw {@code body} is mapped into model objects on demand via the shared
 * {@link Json#MAPPER}; mapping failures throw with the status and body for context rather than being
 * swallowed.
 */
public record ApiResponse(int status, String body) {

    /** Maps the body to a single object of the given type. */
    public <T> T as(Class<T> type) {
        try {
            return Json.MAPPER.readValue(body, type);
        } catch (Exception e) {
            throw new IllegalStateException(
                    "Could not map response body to " + type.getSimpleName()
                            + " (status " + status + "): " + body, e);
        }
    }

    /** Maps the body to a list of the given element type. */
    public <T> List<T> asListOf(Class<T> elementType) {
        try {
            JavaType listType = Json.MAPPER.getTypeFactory()
                    .constructCollectionType(List.class, elementType);
            return Json.MAPPER.readValue(body, listType);
        } catch (Exception e) {
            throw new IllegalStateException(
                    "Could not map response body to List<" + elementType.getSimpleName()
                            + "> (status " + status + "): " + body, e);
        }
    }

    /** Reads a single top-level field as text, for ad-hoc responses without a dedicated model. */
    public String fieldText(String name) {
        Object value = as(Map.class).get(name);
        return value == null ? null : String.valueOf(value);
    }
}
