@ui
Feature: Authentication
  As a visitor of the Map Markers app
  I want to sign in and out through Keycloak
  So that only authenticated users can manage markers

  Scenario: Unauthenticated visitor is prompted to sign in
    Given I open the Map Markers app
    Then I am prompted to sign in

  Scenario: A user signs in through Keycloak
    Given I open the Map Markers app
    When I sign in as "testuser" with password "Test1234!"
    Then I see the Map Markers application shell

  Scenario: A signed-in user signs out
    Given I am signed in as "testuser" with password "Test1234!"
    When I sign out
    Then I am prompted to sign in
