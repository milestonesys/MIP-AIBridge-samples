package services

import (
	"connectivitysample/src/domain/enums"
	"connectivitysample/src/infrastructure/repositories"
	"context"
	"log"
	"sync"
	"time"
)

type TopicRestService struct {
	dataBeingSent  *sync.Map
	graphqlService *GraphqlService
}

func NewTopicRestService(dataBeingSent *sync.Map, graphqlService *GraphqlService) *TopicRestService {
	return &TopicRestService{
		dataBeingSent:  dataBeingSent,
		graphqlService: graphqlService,
	}
}

// Public methods

// Check if the given cameraID is being used for sending data
func (ts *TopicRestService) IsDataBeingSent(cameraID string) bool {

	_, found := ts.dataBeingSent.Load(cameraID)

	return found
}

// Start sending data to a given cameraID every 1 second (fire & forget)
func (ts *TopicRestService) SendDataAsync(cameraID string, streamID string, topicName string, topicFormat int, fileFormat string, files []string) {
	// Add or Update the cameraID in the map
	ts.dataBeingSent.Store(cameraID, "running")
	sourceStreamID := cameraID + "/" + streamID
	log.Printf("Sending data for the SourceStreamID %s", sourceStreamID)

	go func() {

		// data will be sent every 1 seconds.
		ticker := time.NewTicker(1 * time.Second)

		topicRestUrl := ""
		var err error

		switch topicFormat {
		case enums.Event:
			topicRestUrl, err = ts.graphqlService.GetRestEventTopicEndpoint(context.Background(), topicName)
		case enums.Metadata:
			topicRestUrl, err = ts.graphqlService.GetRestMetadataTopicEndpoint(context.Background(), topicName)
		}

		if err != nil {
			log.Printf("Error getting the topic rest url: %s", err)
			return
		}

		index := 0

		for stopSendingData := false; !stopSendingData; {
			// wait for the next tick
			<-ticker.C

			// Check if the data is still being sent
			found := ts.IsDataBeingSent(cameraID)
			if !found {
				stopSendingData = true
				return
			}

			// Alternate between all the loaded files
			currentFile := ""
			switch topicFormat {
			case enums.Event:
				currentFile = TreatEventFile(files[index%len(files)], cameraID)
			case enums.Metadata:
				currentFile = TreatMetadataFile(files[index%len(files)], sourceStreamID)
			}
			index++

			// Publish data into Milestone AI Bridge.
			err := repositories.SendPostRequest(topicRestUrl, currentFile, "text/"+fileFormat)
			if err != nil {
				log.Printf("Couldn't publish the data: %v", err)
			}
		}
	}()
}

// Stop sending data related to a certain cameraID
func (ts *TopicRestService) StopSendingData(cameraID string) {

	log.Printf("Stopping sending data related to the cameraID %s", cameraID)
	ts.dataBeingSent.Delete(cameraID)
}
