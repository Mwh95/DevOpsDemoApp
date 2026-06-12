@ui
Feature: Managing markers in the map UI
  As a signed-in user
  I want to create, edit and delete markers on the map
  So that I can label points of interest

  Background:
    Given I am signed in as "testuser" with password "Test1234!"

  Scenario: The map is shown after signing in
    Then the map is visible

  Scenario: Create a marker using the add-at-center button
    When I enable edit mode
    And I add a marker at the map center with label "Office" and note "Headquarters"
    Then a marker labelled "Office" is shown on the map

  Scenario: Create a marker by clicking the map
    When I enable edit mode
    And I click the map and add a marker with label "Park" and note "Nice spot"
    Then a marker labelled "Park" is shown on the map

  Scenario: Edit an existing marker
    Given a marker labelled "ToEdit" exists on the map
    When I open the marker and change its label to "Edited"
    Then a marker labelled "Edited" is shown on the map

  Scenario: Delete a marker
    Given a marker labelled "ToDelete" exists on the map
    When I open the marker and delete it
    Then no markers are shown on the map
