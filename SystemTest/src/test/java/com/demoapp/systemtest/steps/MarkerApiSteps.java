package com.demoapp.systemtest.steps;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import com.demoapp.systemtest.support.ApiResponse;
import com.demoapp.systemtest.support.World;
import com.fasterxml.jackson.databind.JsonNode;

import io.cucumber.java.en.Given;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;

public class MarkerApiSteps {

    private static final String DEFAULT_PASSWORD = "Test1234!";

    private final World world;

    public MarkerApiSteps(World world) {
        this.world = world;
    }

    @Given("I have an access token for {string}")
    public void iHaveAnAccessTokenFor(String username) {
        world.token = world.api.passwordGrantToken(username, DEFAULT_PASSWORD);
    }

    @When("I GET {string}")
    public void iGet(String path) {
        world.lastResponse = world.api.get(path, world.token);
    }

    @When("I GET {string} without a token")
    public void iGetWithoutToken(String path) {
        world.lastResponse = world.api.getNoAuth(path);
    }

    @When("I create a marker with label {string} and note {string}")
    public void iCreateAMarker(String label, String note) {
        String body = String.format(
                "{\"latitude\":52.52,\"longitude\":13.405,\"label\":%s,\"note\":%s}",
                quote(label), quote(note));
        world.lastResponse = world.api.post("/api/markers", world.token, body);
        if (world.lastResponse.status() == 201) {
            world.createdMarker = world.lastResponse.json();
        }
    }

    @When("I read the created marker")
    public void iReadTheCreatedMarker() {
        world.lastResponse = world.api.get("/api/markers/" + world.createdMarkerId(), world.token);
    }

    @When("I update the created marker label to {string}")
    public void iUpdateTheCreatedMarkerLabel(String label) {
        String body = String.format("{\"label\":%s}", quote(label));
        world.lastResponse = world.api.put("/api/markers/" + world.createdMarkerId(), world.token, body);
        if (world.lastResponse.status() == 200) {
            world.createdMarker = world.lastResponse.json();
        }
    }

    @When("I list the markers")
    public void iListTheMarkers() {
        world.lastResponse = world.api.get("/api/markers", world.token);
    }

    @When("I delete the created marker")
    public void iDeleteTheCreatedMarker() {
        world.lastResponse = world.api.delete("/api/markers/" + world.createdMarkerId(), world.token);
    }

    @When("{string} tries to read that marker")
    public void otherUserTriesToReadThatMarker(String username) {
        String otherToken = world.api.passwordGrantToken(username, DEFAULT_PASSWORD);
        world.lastResponse = world.api.get("/api/markers/" + world.createdMarkerId(), otherToken);
    }

    @Then("the response status is {int}")
    public void theResponseStatusIs(int expected) {
        assertEquals(expected, world.lastResponse.status(),
                "Unexpected status. Body: " + world.lastResponse.body());
    }

    @Then("the response field {string} equals {string}")
    public void theResponseFieldEquals(String field, String value) {
        assertEquals(value, world.lastResponse.fieldText(field),
                "Unexpected value for field '" + field + "'. Body: " + world.lastResponse.body());
    }

    @Then("the created marker has label {string}")
    public void theCreatedMarkerHasLabel(String label) {
        assertNotNull(world.createdMarker, "No created marker available");
        assertEquals(label, world.createdMarker.get("label").asText());
    }

    @Then("the marker list contains {string}")
    public void theMarkerListContains(String label) {
        JsonNode list = world.lastResponse.json();
        assertNotNull(list, "Marker list response was not JSON: " + world.lastResponse.body());
        boolean found = false;
        for (JsonNode marker : list) {
            if (label.equals(marker.path("label").asText())) {
                found = true;
                break;
            }
        }
        assertTrue(found, "Marker list did not contain label '" + label + "': " + world.lastResponse.body());
    }

    private static String quote(String value) {
        return "\"" + value.replace("\\", "\\\\").replace("\"", "\\\"") + "\"";
    }
}
