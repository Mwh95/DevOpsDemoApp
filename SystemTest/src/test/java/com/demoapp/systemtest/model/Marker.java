package com.demoapp.systemtest.model;

/** Marker resource as returned by the Map API ({@code /api/markers}). */
public record Marker(
        String id,
        String userId,
        double latitude,
        double longitude,
        String label,
        String note,
        String createdAt,
        String updatedAt) {
}
