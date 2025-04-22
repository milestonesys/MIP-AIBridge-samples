package handlers

import (
	"bytes"
	services "connectivitysample/src/application/services"
	entities "connectivitysample/src/domain/entities"
	"context"
	"html/template"
	"log"
	"net/http"
)

// Handles all requests coming to the '/snapshot/' endpoint.
type SnapshotHandler struct {
	tokenService          *services.TokenService
	queryStringService    *services.QueryStringService
	graphqlService        *services.GraphqlService
	commandLineParameters *entities.CommandLineParameters
}

func NewSnapshotHandler(tokenService *services.TokenService, graphqlService *services.GraphqlService,
	queryStringService *services.QueryStringService, commandLineParameters *entities.CommandLineParameters) *SnapshotHandler {
	return &SnapshotHandler{
		tokenService:          tokenService,
		graphqlService:        graphqlService,
		queryStringService:    queryStringService,
		commandLineParameters: commandLineParameters,
	}
}

// Renders 'snapshot-camera-page.html' or 'snapshot-configuration-page.html'  when '/snapshot/' endpoint gets requested (commonly from MC)
//
// 'snapshot-camera-page.html' page:
// 1 - Shows a 'live' snapshot of a given camera.
// 2 - Can be loaded in MC. Recording Servers (node) -> Select a camera -> select 'properties' tab -> select 'Processing server' tab -> select 'getsnapshot' topic from the treeview.
//
// 'snapshot-configuration-page.html' page (Optional):
// 1 - Can be used to tweak/configure settings related to the 'getsnapshot' topic.
// 2 - Can be loaded in MC: Processing Servers (node) -> Select a processing server (item) -> select 'connectivity sample' from the 'registered applications' list -> select 'getsnapshot' from the 'topics' list.
func (sh *SnapshotHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Extract token from request headers and verify the token is valid
	token, err := sh.tokenService.ExtractAndVerifyToken(r, w)
	if err != nil {
		log.Println("Error at SnapshotHandler method, during extractAndVerifyToken execution:" + err.Error())
		return
	}

	// Is a configuration Page Request?
	if r.URL.Path == "/"+sh.commandLineParameters.AppUrlPath()+"/snapshot/" {
		err := snapshotConfigurationPage(w)
		if err != nil {
			log.Println("Error at SnapshotHandler method, during snapshotConfigurationPage execution:" + err.Error())
		}
		return
	}

	// we expect the path to provide a device id and stream id
	queryStringContext, err := sh.queryStringService.ExtractCameraAndStreamIDsFromPath("/snapshot/(.+?)/(.+)", r, w)
	if err != nil {
		log.Println("Error at SnapshotHandler method, during extractDeviceAndStreamIDsFromPath execution:" + err.Error())
		return
	}
	cameraID := queryStringContext.CameraID
	streamID := queryStringContext.StreamID

	// Gets the base64 encoded snapshot image for a given device and stream
	innerURL, err := sh.graphqlService.GetSnapshot(context.Background(), cameraID, streamID, token, sh.commandLineParameters)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Error at SnapshotHandler method, during querySnapshotSignallingEndpoint execution:" + err.Error())
		return
	}

	// Return the snapshot-camera-page.html
	path := "templates/snapshot-camera-page.html"
	template, err := template.ParseFS(templateFS, path)
	if err != nil {
		log.Printf("Error parsing template file %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Data to be passed to the template
	pageData := struct {
		ImageUrl string
	}{
		ImageUrl: innerURL,
	}
	// Write to the response
	template.Execute(w, pageData)
}

// Can be used to tweak/configure settings related to the 'getsnapshot' topic.
func snapshotConfigurationPage(w http.ResponseWriter) error {
	path := "templates/snapshot-configuration-page.html"
	template, err := template.ParseFS(templateFS, path)
	if err != nil {
		log.Printf("Error parsing template file %s: %v", path, err)
		return err
	}

	// Create a buffer to capture the output
	var output bytes.Buffer

	// Execute the template and write to the buffer
	err = template.Execute(&output, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return err
	}
	// Write to the response
	template.Execute(w, nil)

	return nil
}
