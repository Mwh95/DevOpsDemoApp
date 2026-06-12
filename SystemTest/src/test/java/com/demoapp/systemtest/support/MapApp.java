package com.demoapp.systemtest.support;

import static com.microsoft.playwright.assertions.PlaywrightAssertions.assertThat;

import com.microsoft.playwright.Locator;
import com.microsoft.playwright.Page;
import com.microsoft.playwright.options.AriaRole;

/** Page object encapsulating interactions with the Map Markers SPA and the Keycloak login form. */
public class MapApp {

    private final Page page;

    public MapApp(Page page) {
        this.page = page;
    }

    public void open() {
        page.navigate("/");
    }

    public Locator signInPrompt() {
        return page.getByText("You need to sign in to use the map.");
    }

    public Locator appHeader() {
        return page.getByText("Map Markers");
    }

    public void clickSignIn() {
        page.getByRole(AriaRole.BUTTON, new Page.GetByRoleOptions().setName("Sign in")).click();
    }

    /** Performs the full OIDC login round-trip via the Keycloak login form. */
    public void login(String username, String password) {
        clickSignIn();
        page.locator("#username").fill(username);
        page.locator("#password").fill(password);
        page.locator("#kc-login").click();
        assertThat(appHeader()).isVisible();
    }

    /**
     * Triggers RP-initiated logout. Keycloak ends the session and (without a configured
     * post-logout redirect) lands on its own page, so we re-open the SPA afterwards; the cleared
     * session means the sign-in prompt is shown again.
     */
    public void signOut() {
        page.getByRole(AriaRole.BUTTON, new Page.GetByRoleOptions().setName("Sign out")).click();
        page.waitForURL("**/protocol/openid-connect/logout**");
        open();
    }

    public void enableEditMode() {
        // The checkbox itself sits underneath Leaflet's zoom control, so click the surrounding
        // label text (which extends clear of the control) to toggle the wrapped checkbox.
        Locator checkbox = page.getByRole(AriaRole.CHECKBOX);
        if (!checkbox.isChecked()) {
            page.locator(".map-toolbar label").click();
        }
        assertThat(checkbox).isChecked();
    }

    public void clickAddMarkerAtCenter() {
        page.getByRole(AriaRole.BUTTON,
                new Page.GetByRoleOptions().setName("Add marker at map center")).click();
    }

    /** Clicks the map at a fixed offset to place a pending marker at a distinct location. */
    public void clickMapAt(int x, int y) {
        page.locator(".leaflet-container").click(new Locator.ClickOptions()
                .setPosition(x, y));
    }

    public void fillMarkerForm(String label, String note) {
        assertThat(markerDialog()).isVisible();
        page.locator("#marker-label").fill(label);
        page.locator("#marker-note").fill(note);
        page.getByRole(AriaRole.BUTTON, new Page.GetByRoleOptions().setName("Save")).click();
    }

    public Locator markerDialog() {
        return page.getByRole(AriaRole.DIALOG);
    }

    public Locator markerIcons() {
        return page.locator(".leaflet-marker-icon");
    }

    public void openFirstMarkerPopup() {
        if (!popupLabel().isVisible()) {
            markerIcons().first().click();
        }
        assertThat(popupLabel()).isVisible();
    }

    public void changeLabel(String newLabel) {
        assertThat(markerDialog()).isVisible();
        page.locator("#marker-label").fill(newLabel);
        page.getByRole(AriaRole.BUTTON, new Page.GetByRoleOptions().setName("Save")).click();
    }

    public Locator popupLabel() {
        return page.locator(".leaflet-popup-content .popup-label");
    }

    public void clickPopupEdit() {
        popup().getByRole(AriaRole.BUTTON, new Locator.GetByRoleOptions().setName("Edit")).click();
    }

    public void clickPopupDelete() {
        popup().getByRole(AriaRole.BUTTON, new Locator.GetByRoleOptions().setName("Delete")).click();
    }

    public Locator mapContainer() {
        return page.locator(".leaflet-container");
    }

    private Locator popup() {
        return page.locator(".leaflet-popup-content");
    }
}
