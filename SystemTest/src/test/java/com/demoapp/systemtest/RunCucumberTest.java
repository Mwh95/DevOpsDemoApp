package com.demoapp.systemtest;

import static io.cucumber.junit.platform.engine.Constants.GLUE_PROPERTY_NAME;
import static io.cucumber.junit.platform.engine.Constants.PLUGIN_PROPERTY_NAME;

import org.junit.platform.suite.api.ConfigurationParameter;
import org.junit.platform.suite.api.IncludeEngines;
import org.junit.platform.suite.api.SelectClasspathResource;
import org.junit.platform.suite.api.Suite;

/** JUnit Platform suite that runs every Gherkin feature through the Cucumber engine. */
@Suite
@IncludeEngines("cucumber")
@SelectClasspathResource("features")
@ConfigurationParameter(key = GLUE_PROPERTY_NAME, value = "com.demoapp.systemtest")
@ConfigurationParameter(
        key = PLUGIN_PROPERTY_NAME,
        value = "pretty, "
                + "html:build/reports/cucumber/cucumber.html, "
                + "json:build/reports/cucumber/cucumber.json")
public class RunCucumberTest {
}
