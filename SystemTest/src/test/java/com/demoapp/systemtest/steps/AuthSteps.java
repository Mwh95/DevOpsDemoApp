package com.demoapp.systemtest.steps;

import static com.microsoft.playwright.assertions.PlaywrightAssertions.assertThat;

import com.demoapp.systemtest.support.MapApp;
import com.demoapp.systemtest.support.World;

import io.cucumber.java.en.Given;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;

public class AuthSteps {

    private final World world;

    public AuthSteps(World world) {
        this.world = world;
    }

    private MapApp app() {
        return new MapApp(world.page);
    }

    @Given("I open the Map Markers app")
    public void iOpenTheApp() {
        app().open();
    }

    @Given("I am signed in as {string} with password {string}")
    public void iAmSignedIn(String username, String password) {
        MapApp app = app();
        app.open();
        app.login(username, password);
    }

    @When("I sign in as {string} with password {string}")
    public void iSignIn(String username, String password) {
        app().login(username, password);
    }

    @When("I sign out")
    public void iSignOut() {
        app().signOut();
    }

    @Then("I am prompted to sign in")
    public void iAmPromptedToSignIn() {
        assertThat(app().signInPrompt()).isVisible();
    }

    @Then("I see the Map Markers application shell")
    public void iSeeTheAppShell() {
        assertThat(app().appHeader()).isVisible();
    }
}
