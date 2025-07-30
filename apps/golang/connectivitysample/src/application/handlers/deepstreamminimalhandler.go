package handlers

import (
	services "connectivitysample/src/application/services"
	entities "connectivitysample/src/domain/entities"
	"connectivitysample/src/domain/enums"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// Handles all requests coming to the '/deepstreamminimal' endpoint.
type DeepstreamMinimalHandler struct {
	queryStringService    *services.QueryStringService
	TopicRestService      *services.TopicRestService
	commandLineParameters *entities.CommandLineParameters
	fileReader            *services.FileReader
}

func NewDeepstreamMinimalHandler(queryStringService *services.QueryStringService,
	TopicRestService *services.TopicRestService, commandLineParameters *entities.CommandLineParameters) *DeepstreamMinimalHandler {
	return &DeepstreamMinimalHandler{
		queryStringService:    queryStringService,
		TopicRestService:      TopicRestService,
		commandLineParameters: commandLineParameters,
		fileReader:            services.NewFileReader("dsmin", "json"),
	}
}

// Renders 'deepstreamminimal-metadata-page.html' when '/deepstreamminimal' endpoint gets requested. (commonly from MC)
// This page:
// 1 - Allows a MC user to start or stop sending Deepstream Minimal metadata to a certain camera.
// 2 - Can be loaded in MC> Recording Servers (node) -> Select a camera -> select 'properties' tab -> select 'Processing server' tab -> select 'senddeepstreamminimal' topic from the treeview.

func (dsh *DeepstreamMinimalHandler) Handle(w http.ResponseWriter, r *http.Request) {

	// The path can include a device id and stream id. If so, we extract it.
	// URL to load on the processing-server tab from a certain cameraID: https://${EXTERNAL_HOSTNAME}:7443/metadata/senddeepstreamminimal/camID/streamIDcamID
	queryStringContext, err := dsh.queryStringService.ExtractTopicNameCameraAndStreamIDsFromPath("/deepstreamminimal/(.+)/(.+?)/.*/(.+)", r, w)
	if err != nil {
		// The path didn't include device id and stream id.
		// Since we didn't implement a 'configuration' page for this topic, we just return.
		return
	}

	topicName := queryStringContext.TopicName
	cameraID := queryStringContext.CameraID
	streamID := queryStringContext.StreamID

	isMetadataBeingSent := dsh.TopicRestService.IsDataBeingSent(cameraID)
	metadataSendingCurrentStatus := startStatus
	if isMetadataBeingSent {
		metadataSendingCurrentStatus = stopStatus
	}

	// return the deepstreamminimal-metadata-page.html
	path := "templates/deepstreamminimal-metadata-page.html"
	template, err := template.ParseFS(templateFS, path)
	if err != nil {
		log.Printf("Error parsing template file %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// data to be passed to the template
	pageData := struct {
		CameraID   string
		StreamID   string
		TopicName  string
		Status     string
		AppUrlPath string
	}{
		CameraID:   cameraID,
		StreamID:   streamID,
		TopicName:  topicName,
		Status:     metadataSendingCurrentStatus,
		AppUrlPath: dsh.commandLineParameters.AppUrlPath(),
	}
	// write to the response
	template.Execute(w, pageData)
}

// Handles the requests coming to the '/deepstreamminimal/processing/senddeepstreamminimal' endpoint
// This happens when the user clicks the 'start'-'stop' button in the MC
// As a result, we stop sending metadata (if we were sending it) or start sending metadata (if we were not)
func (dsh *DeepstreamMinimalHandler) ProcessingHandle(w http.ResponseWriter, r *http.Request) {
	// Request body coming in
	var eventData struct {
		CameraID  string `json:"cameraId"`
		StreamID  string `json:"streamId"`
		TopicName string `json:"topicName"`
	}
	err := json.NewDecoder(r.Body).Decode(&eventData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Response body coming out
	var eventDataResponse struct {
		EventStatus string `json:"EventStatus"`
	}

	isMetadataBeingSent := dsh.TopicRestService.IsDataBeingSent(eventData.CameraID)
	if isMetadataBeingSent {
		dsh.TopicRestService.StopSendingData(eventData.CameraID)
		eventDataResponse.EventStatus = startStatus
	} else {

		// Get the Deepstream Minimal metadata content
		jsons, err := dsh.fileReader.ReadMultipleFiles()
		if err != nil {
			log.Printf("Error getting the Deepstream Minimal metadata content: %s", err)
			return
		}
		log.Printf("Sending metadata")
		dsh.TopicRestService.SendDataAsync(eventData.CameraID, eventData.StreamID, eventData.TopicName, enums.Metadata, "json", jsons)
		eventDataResponse.EventStatus = stopStatus
	}

	// Respond to the client (Button Start or Stop text in MC)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventDataResponse)
}
