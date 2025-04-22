# Register files

The register files contain a GraphQL register mutation input. You can read about this mutation at the `<AI Bridge External Hostname>:4000` (Which is accessible only when Milestone AI Bridge is running in debug mode).

In the registration input a list of apps can be provided. Each app can register it self (you can see example in the code on how this sample register it self - check the `graphqlrepository.go` file for more info). 

## Basic fields

Each app must have an `id`, `url`, `name`, `version` and a `description`. 

```graphql
id: "28a6bc9a-0833-46c6-958e-19da4ee6d9e5"
url: "https://${EXTERNAL_HOSTNAME}:${APP_WEBSERVER_PORT}/${APP_URL_PATH}/"
name: "Connectivity sample"
version:"1.0.0"
description: "Sample IVA to showcase how IVAs work."
```

## Manufacturer field

Optionally, the app `manufacturer` info can be provided. 

As a Milestone Technology Partner, you need to enable the MIP SDK features related to the Installed
Integration Insights initiative that allows us to gather datapoints about your integrations.

During the verification process, a partner is asked to provide their manufacturer name to be added to the white list. 

The manufacturer information is used to identify if the manufacturer is in our white list. The manufacturer white list contains all the manufacturers we are allowed to store integrations data from. This means whenever an application registers is self at Milestone AI Bridge and the manufacturer data was provided we check if the manufacturer is on our white list to store its integration data or not.


```graphql
manufacturer: {
  name: "Sample Manufacturer"
}
```

For more information, refer to the [How to guide - "Manufacturer Name & Integration Name in Installed Integration Insights"](https://content.milestonesys.com/media/?mediaId=2926DE15-79F9-4B5A-91248030F2DF36CC)

## Application topics

This sample app only exposes event topics, therefore in the register file there are only one list of eventTopics. However, some other apps might have eventTopics, videoTopics and metadataTopics (Check the `<AI Bridge External Hostname>:4000` for more info).

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