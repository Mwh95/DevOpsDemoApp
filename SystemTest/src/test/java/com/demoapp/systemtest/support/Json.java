package com.demoapp.systemtest.support;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.fasterxml.jackson.databind.json.JsonMapper;

/**
 * Shared, pre-configured Jackson mapper for deserialising API responses into model objects.
 * Snake-case JSON (e.g. {@code user_id}, {@code access_token}) maps onto camelCase record
 * components, and unknown fields are ignored so the models can stay focused.
 */
public final class Json {

    public static final ObjectMapper MAPPER = JsonMapper.builder()
            .propertyNamingStrategy(PropertyNamingStrategies.SNAKE_CASE)
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
            .build();

    private Json() {
    }
}
