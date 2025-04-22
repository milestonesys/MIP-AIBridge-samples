package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"
)

// Service used to send & stop analytic events into the VMS.
type AnalyticEventService struct {
	eventsBeingSent *sync.Map // In-memory map to keep track of analytic events being sent (which are bound to certain cameraIDs)
	graphqlService  *GraphqlService
}

func NewAnalyticEventService(eventsBeingSent *sync.Map, graphqlService *GraphqlService) *AnalyticEventService {
	return &AnalyticEventService{
		eventsBeingSent: eventsBeingSent,
		graphqlService:  graphqlService,
	}
}

// Public methods

// Check if the given cameraID is being used for sending analytic events.
func (aes *AnalyticEventService) IsEventBeingSent(cameraID string) bool {

	_, found := aes.eventsBeingSent.Load(cameraID)

	return found
}

// Start sending analytic events for a given cameraID every 1 second (fire & forget)
func (aes *AnalyticEventService) SendEventAsync(cameraID string, topicName string) {
	// Add or Update the cameraID in the map
	aes.eventsBeingSent.Store(cameraID, "running")
	log.Printf("Sending an event related to the cameraID %s", cameraID)

	go func() {

		// Analytical event will be sent every 1 seconds.
		ticker := time.NewTicker(1 * time.Second)

		topicRestUrl, err := aes.graphqlService.GetRestEventTopicEndpoint(context.Background(), topicName)
		if err != nil {
			log.Printf("Error getting the topic rest url: %s", err)
			return
		}

		for stopSendingEvent := false; !stopSendingEvent; {
			// wait for the next tick
			<-ticker.C

			// Check if the event is still being sent
			found := aes.IsEventBeingSent(cameraID)
			if !found {
				stopSendingEvent = true
				return
			}

			// Get the analytical event content
			analyticalEventContent, err := getAnalyticalEventContent(cameraID)
			if err != nil {
				log.Printf("Error getting the analytical event content: %s", err)
				return
			}

			// Publish event into Milestone AI Bridge.
			publishAnalyticalEvent(topicRestUrl, analyticalEventContent)
		}
	}()
}

// Stop sending analytic events related to a certain cameraID
func (aes *AnalyticEventService) StopSendingEvent(cameraID string) {

	log.Printf("Stopping sending an event related to the cameraID %s", cameraID)
	aes.eventsBeingSent.Delete(cameraID)
}

// Private methods

// Publish analytic event into Milestone AI Bridge
func publishAnalyticalEvent(topicRestUrl string, analyticalEventContent string) error {

	// Publish event
	response, err := http.Post(topicRestUrl,
		"application/json",
		bytes.NewBuffer([]byte(analyticalEventContent)))

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		errorMessage := fmt.Sprintf("Error when publishing the analytical event. Response status:'%s' Body:'%s'", response.Status, string(body))
		return errors.New(errorMessage)
	}
	return nil
}

// Get a valid analytic event for a given cameraID
func getAnalyticalEventContent(cameraID string) (string, error) {

	// return the analyticEvent.json
	path := "templates/analyticEvent.json"
	template, err := template.ParseFS(templateFS, path)
	if err != nil {
		log.Printf("Error parsing template file %s: %v", path, err)
		return "", err
	}

	// data to be passed to the json template
	analyticalEventPlaceHolders := struct {
		CameraID string
	}{
		CameraID: cameraID,
	}

	// Create a buffer to capture the output
	var output bytes.Buffer

	// Execute the template and write to the buffer
	err = template.Execute(&output, analyticalEventPlaceHolders)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return "", err
	}

	// Convert the buffer to a string and return
	return output.String(), nil
}
