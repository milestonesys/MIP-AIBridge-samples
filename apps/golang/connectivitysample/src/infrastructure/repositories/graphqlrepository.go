package repositories

import (
	"bytes"
	"connectivitysample/src/domain/entities"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var contentType string = "application/json"

// Fetches data through Milestone AI Bridge Graphql API.
type GraphqlRepository struct {
	client *http.Client
}

func NewGraphqlRepository() *GraphqlRepository {
	return &GraphqlRepository{
		client: &http.Client{},
	}
}

// Public methods

/* --- GraphQL queries ----------------------------------------------------- */

// Query all VMS system registered issuers connected to AIBridge
func (gr *GraphqlRepository) GetIdentityProviderRegisteredIssuers(ctx context.Context, requestUrl string) ([]string, error) {
	// Build request
	request := `{
	"query": "query { about { videoManagementSystems { idp } } }"
}`

	// Send request
	response, err := gr.sendRequest(ctx, requestUrl, request)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result struct {
		Data struct {
			About struct {
				VideoManagementSystems []struct {
					IDP string `json:"idp"`
				} `json:"videoManagementSystems"`
			} `json:"about"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var registeredIssuers []string
	for _, vms := range result.Data.About.VideoManagementSystems {
		registeredIssuers = append(registeredIssuers, vms.IDP)
	}

	// We are done
	return registeredIssuers, nil
}

// Query all VMS system connected to AIBridge
func (gr *GraphqlRepository) GetVmsIds(ctx context.Context, requestUrl string) ([]string, error) {
	// Build request
	request := `{
	"query": "query { about { videoManagementSystems { id } } }",
	"variables": null
}`

	// Send request
	response, err := gr.sendRequest(ctx, requestUrl, request)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result struct {
		Data struct {
			About struct {
				VideoManagementSystems []struct {
					ID string `json:"id"`
				} `json:"videoManagementSystems"`
			} `json:"about"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var vmsIds []string
	for _, vms := range result.Data.About.VideoManagementSystems {
		vmsIds = append(vmsIds, vms.ID)
	}
	return vmsIds, nil
}

// Query the snapshot signalling endpoint for a given deviceID and streamID
func (gr *GraphqlRepository) GetSnapshot(ctx context.Context, requestUrl, deviceID, streamID, token string, commandLineParameters *entities.CommandLineParameters) (string, error) {
	// Build request
	request := `{
		"query" : "query GetSnapshot($deviceID: ID!, $streamID: ID!, $max_width: Int!, $max_height: Int!, $token:String!) { cameras(deviceIDs: [$deviceID]) { videoStreams(streamID: $streamID) { snapshot(maxWidth: $max_width, maxHeight: $max_height, token: $token) { jpegImage } } } }",
		"variables" : { "deviceID" : "` + deviceID + `", "streamID" : "` + streamID + `", "max_width":  ` + fmt.Sprintf("%d", commandLineParameters.SnapshotConfiguration().SnapshotMaxWidth()) + `, "max_height": ` + fmt.Sprintf("%d", commandLineParameters.SnapshotConfiguration().SnapshotMaxHeight()) + `, "token": "` + token + `" }
	}`

	// Send request
	response, err := gr.sendRequest(ctx, requestUrl, request)
	if err != nil {
		return "", err
	}

	// Struct to parse response
	var result struct {
		Data struct {
			Cameras []struct {
				VideoStreams []struct {
					SnapShot struct {
						JpegImage string `json:"jpegImage"`
					} `json:"snapshot"`
				} `json:"videoStreams"`
			} `json:"cameras"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Data.Cameras) != 1 || len(result.Data.Cameras[0].VideoStreams) != 1 {
		return "", errors.New("requested device / stream was not found")
	}

	// If everything went well, we return the Snapshot in Base64
	return result.Data.Cameras[0].VideoStreams[0].SnapShot.JpegImage, nil
}

// Gets the rest event topic endpoint for a given topic name.
func (gr *GraphqlRepository) GetRestEventTopicEndpoint(ctx context.Context, requestUrl, topicName string) (string, error) {
	// Build request
	request := `{
	"query": "query Query_Event_Topics_By_Name($topicName:String!) {eventTopics(topicName:$topicName){topicAvailability{rest}}}",
	"variables" : { "topicName": "` + topicName + `"}
	}`

	// Send request
	response, err := gr.sendRequest(ctx, requestUrl, request)
	if err != nil {
		return "", err
	}

	// Struct to parse response
	var result struct {
		Data struct {
			EventTopics []struct {
				TopicAvailability struct {
					Rest string `json:"rest"`
				} `json:"topicAvailability"`
			} `json:"eventTopics"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Data.EventTopics) != 1 {
		return "", errors.New("requested topic name was not found")
	}

	// If everything went well, we return the rest endpoint for this topic.
	return result.Data.EventTopics[0].TopicAvailability.Rest, nil
}

// Gets the rest metadata topic endpoint for a given topic name.
func (gr *GraphqlRepository) GetRestMetadataTopicEndpoint(ctx context.Context, requestUrl, topicName string) (string, error) {
	// Build request
	request := `{
	"query": "query Query_Metadata_Topics_By_Name($topicName:String!) {metadataTopics(topicName:$topicName){topicAvailability{rest}}}",
	"variables" : { "topicName": "` + topicName + `"}
	}`

	// Send request
	response, err := gr.sendRequest(ctx, requestUrl, request)
	if err != nil {
		return "", err
	}

	// Struct to parse response
	var result struct {
		Data struct {
			MetadataTopics []struct {
				TopicAvailability struct {
					Rest string `json:"rest"`
				} `json:"topicAvailability"`
			} `json:"metadataTopics"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Data.MetadataTopics) != 1 {
		return "", errors.New("requested topic name was not found")
	}

	// If everything went well, we return the rest endpoint for this topic.
	return result.Data.MetadataTopics[0].TopicAvailability.Rest, nil
}

/* --- GraphQL mutations ----------------------------------------------------- */

// Register the app in AIBridge to all connected endpoints
func (gr *GraphqlRepository) RegisterConnectivitySample(ctx context.Context, requestUrl, populatedRegistrationFileContent string, vmsId string) error {

	// Build request
	requestObj := map[string]string{"query": `mutation { register(input: { id: "` + vmsId + `"` + populatedRegistrationFileContent + `}) { id } }`}
	requestVal, err := json.Marshal(requestObj)
	if err != nil {
		return err
	}

	// Send request
	response, err := gr.sendRequest(ctx, requestUrl, string(requestVal))
	if err != nil {
		return err
	}

	// Parse response and extract endpoint id
	// anonymous struct to receive the response from the server:
	var result struct {
		Data struct {
			Register struct {
				ID string `json:"id"`
			} `json:"register"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}

	// Success
	return nil
}

// Private methods

// Sends graphql request to a given url. The query is defined in the request body
func (gr *GraphqlRepository) sendRequest(ctx context.Context, requestUrl, body string) (*http.Response, error) {
	// Create request
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	// Set the content type
	request.Header.Add("Content-Type", contentType)

	// Execute request
	response, err := gr.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		responseBody, _ := io.ReadAll(response.Body)
		return nil, errors.New(response.Status + " status code received. Message: " + string(responseBody))
	}

	return response, nil
}
