package com.demoapp.systemtest.hooks;

import com.demoapp.systemtest.support.ApiResponse;
import com.demoapp.systemtest.support.World;
import com.microsoft.playwright.Browser;
import com.microsoft.playwright.BrowserType;
import com.microsoft.playwright.Page;
import com.microsoft.playwright.Playwright;
import com.fasterxml.jackson.databind.JsonNode;

import io.cucumber.java.After;
import io.cucumber.java.Before;
import io.cucumber.java.Scenario;

/**
 * Lifecycle hooks for the browser. The {@link Playwright} engine and {@link Browser} are created
 * once and reused; each {@code @ui} scenario gets a fresh context and page for isolation, and the
 * test user's markers are cleared beforehand so UI assertions start from a clean slate.
 */
public class PlaywrightHooks {

    private static final String UI_USER = "testuser";
    private static final String UI_PASSWORD = "Test1234!";

    private static Playwright playwright;
    private static Browser browser;

    private final World world;

    public PlaywrightHooks(World world) {
        this.world = world;
    }

    @Before(value = "@ui", order = 100)
    public void setUpBrowser() {
        ensureBrowser();
        clearMarkers();
        world.browserContext = browser.newContext(new Browser.NewContextOptions()
                .setBaseURL(world.baseUrl())
                .setViewportSize(1280, 800)
                .setIgnoreHTTPSErrors(true));
        world.browserContext.setDefaultTimeout(30_000);
        world.page = world.browserContext.newPage();
    }

    @After(value = "@ui", order = 100)
    public void tearDownBrowser(Scenario scenario) {
        Page page = world.page;
        if (page != null) {
            if (scenario.isFailed()) {
                scenario.attach(page.screenshot(), "image/png", "failure-screenshot");
            }
        }
        if (world.browserContext != null) {
            world.browserContext.close();
        }
    }

    private void clearMarkers() {
        String token = world.api.passwordGrantToken(UI_USER, UI_PASSWORD);
        ApiResponse list = world.api.get("/api/markers", token);
        if (list.json() != null && list.json().isArray()) {
            for (JsonNode marker : list.json()) {
                world.api.delete("/api/markers/" + marker.get("id").asText(), token);
            }
        }
    }

    private static synchronized void ensureBrowser() {
        if (browser == null) {
            playwright = Playwright.create();
            browser = playwright.chromium().launch(new BrowserType.LaunchOptions().setHeadless(true));
            Runtime.getRuntime().addShutdownHook(new Thread(() -> {
                try {
                    browser.close();
                } finally {
                    playwright.close();
                }
            }));
        }
    }
}
