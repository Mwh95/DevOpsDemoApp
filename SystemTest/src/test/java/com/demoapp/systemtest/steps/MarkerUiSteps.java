package com.demoapp.systemtest.steps;

import static com.microsoft.playwright.assertions.PlaywrightAssertions.assertThat;

import com.demoapp.systemtest.support.MapApp;
import com.demoapp.systemtest.support.World;

import io.cucumber.java.en.Given;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;

public class MarkerUiSteps {

    private final World world;

    public MarkerUiSteps(World world) {
        this.world = world;
    }

    private MapApp app() {
        return new MapApp(world.page);
    }

    @Then("the map is visible")
    public void theMapIsVisible() {
        assertThat(app().mapContainer()).isVisible();
    }

    @When("I enable edit mode")
    public void iEnableEditMode() {
        app().enableEditMode();
    }

    @When("I add a marker at the map center with label {string} and note {string}")
    public void iAddMarkerAtCenter(String label, String note) {
        MapApp app = app();
        app.clickAddMarkerAtCenter();
        app.fillMarkerForm(label, note);
    }

    @When("I click the map and add a marker with label {string} and note {string}")
    public void iClickMapAndAddMarker(String label, String note) {
        MapApp app = app();
        app.clickMapAt(640, 400);
        app.fillMarkerForm(label, note);
    }

    @Given("a marker labelled {string} exists on the map")
    public void aMarkerExists(String label) {
        MapApp app = app();
        app.enableEditMode();
        app.clickAddMarkerAtCenter();
        app.fillMarkerForm(label, "seed note");
        assertThat(app.markerIcons()).hasCount(1);
    }

    @When("I open the marker and change its label to {string}")
    public void iOpenAndChangeLabel(String newLabel) {
        MapApp app = app();
        app.openFirstMarkerPopup();
        app.clickPopupEdit();
        app.changeLabel(newLabel);
    }

    @When("I open the marker and delete it")
    public void iOpenAndDelete() {
        MapApp app = app();
        app.openFirstMarkerPopup();
        app.clickPopupDelete();
    }

    @Then("a marker labelled {string} is shown on the map")
    public void aMarkerLabelledIsShown(String label) {
        MapApp app = app();
        assertThat(app.markerIcons()).hasCount(1);
        app.openFirstMarkerPopup();
        assertThat(app.popupLabel()).hasText(label);
    }

    @Then("no markers are shown on the map")
    public void noMarkersAreShown() {
        assertThat(app().markerIcons()).hasCount(0);
    }
}
