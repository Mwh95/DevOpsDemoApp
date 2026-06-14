@api
Feature: Map markers REST API
  As a client of the Map API
  I want the marker endpoints to behave correctly
  So that the application data stays consistent and secure

  Scenario: Liveness endpoint reports the service is up
    When I GET "/public/health/live"
    Then the response status is 200
    And the response field "status" equals "UP"

  Scenario: Readiness endpoint reports the service is up
    When I GET "/public/health/ready"
    Then the response status is 200
    And the response field "status" equals "UP"

  Scenario: The markers endpoint rejects unauthenticated requests
    When I GET "/api/markers" without a token
    Then the response status is 401

  Scenario: A marker can be created, read, updated, listed and deleted
    Given I have an access token for "testuser"
    When I create a marker with label "API Home" and note "from api"
    Then the response status is 201
    And the created marker has label "API Home"
    When I read the created marker
    Then the response status is 200
    When I update the created marker label to "API Updated"
    Then the response status is 200
    And the created marker has label "API Updated"
    When I list the markers
    Then the marker list contains "API Updated"
    When I delete the created marker
    Then the response status is 204
    When I read the created marker
    Then the response status is 404

  Scenario: Markers created by one user are not visible to another
    Given I have an access token for "testuser"
    And I create a marker with label "Private" and note "secret"
    When "otheruser" tries to read that marker
    Then the response status is 404
