# Register files

The register files contain a GraphQL register mutation input. You can read about this mutation in the reference manual at the `<AI Bridge External Hostname>:4000`.

**NOTE:**

The reference manual can only be accessed when Milestone AI Bridge is running in debug mode.

A list of apps is provided in the registration input with each app being self-registering.
An example of how the sample IVA app has been configured to be self-registering can be seen in the `RegisterConnectivitySample` method, found in the [graphqlrepository.go](../src/infrastructure/repositories/graphqlrepository.go) file

## Basic fields

Each app must have an `id`, `url`, `name`, `version`, and a `description`.

```graphql
id: "28a6bc9a-0833-46c6-958e-19da4ee6d9e5"
url: "${TLS_SCHEME}://${EXTERNAL_HOSTNAME}:${APP_WEBSERVER_PORT}/${APP_URL_PATH}"
name: "Connectivity sample"
version: "1.0.0"
description: "Sample IVA to showcase how IVAs work."
```

## Manufacturer field

This `optional` field is only relevant for Milestone Technology Partners. For more information, see to the [How to guide - "Manufacturer Name & Integration Name in Installed Integration Insights"](https://content.milestonesys.com/media/?mediaId=2926DE15-79F9-4B5A-91248030F2DF36CC).

```graphql
manufacturer: {
  name: "Sample Manufacturer"
}
```

## Application topics

This sample IVA app only exposes event topics, which is why there is only one list of eventTopics in the register file. Other apps may contain any combination of eventTopics, videoTopics and metadataTopics in their register file.
For more information, see the reference manual at the `<AI Bridge External Hostname>:4000`.

```graphql
eventTopics: [ {
  url: "${TLS_SCHEME}://${EXTERNAL_HOSTNAME}:${APP_WEBSERVER_PORT}/${APP_URL_PATH}/snapshot"
  name: "getsnapshot"
  description: "Get a Snapshot from the video"
  eventFormat: ANALYTICS_EVENT
}, {
  url: "${TLS_SCHEME}://${EXTERNAL_HOSTNAME}:${APP_WEBSERVER_PORT}/${APP_URL_PATH}/event/sendanalyticevents"
  name: "sendanalyticevents"
  description: "Send Analytic Events"
  eventFormat: ANALYTICS_EVENT
} ]
```
