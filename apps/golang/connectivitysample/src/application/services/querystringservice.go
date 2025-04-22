package services

import (
	"errors"
	"net/http"
	"regexp"
)

// Service used to handle incoming query string data
type QueryStringService struct {
}

func NewQueryStringService() *QueryStringService {
	return &QueryStringService{}
}

// Extracts the cameraID and streamID from the request's query string
func (hs *QueryStringService) ExtractCameraAndStreamIDsFromPath(pattern string, r *http.Request, w http.ResponseWriter) (struct {
	CameraID string
	StreamID string
}, error) {
	reg := regexp.MustCompile(pattern)
	res := reg.FindStringSubmatch(r.URL.Path)
	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		return struct {
			CameraID string
			StreamID string
		}{}, errors.New("error at ExtractCameraAndStreamIDsFromPath method, during extraction of pattern " + pattern + "used to extract the DeviceID and StreamID")
	}
	return struct {
		CameraID string
		StreamID string
	}{
		CameraID: res[1],
		StreamID: res[2],
	}, nil
}

// Extracts the topicName, deviceID and streamID from the request's query string
func (hs *QueryStringService) ExtractTopicNameCameraAndStreamIDsFromPath(pattern string, r *http.Request, w http.ResponseWriter) (struct {
	TopicName string
	CameraID  string
	StreamID  string
}, error) {
	reg := regexp.MustCompile(pattern)
	res := reg.FindStringSubmatch(r.URL.Path)
	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		return struct {
			TopicName string
			CameraID  string
			StreamID  string
		}{}, errors.New("error at ExtractTopicNameDeviceAndStreamIDsFromPath method, during extraction of pattern " + pattern + " used to extract the TopicName, DeviceID, and StreamID")
	}
	return struct {
		TopicName string
		CameraID  string
		StreamID  string
	}{
		TopicName: res[1],
		CameraID:  res[2],
		StreamID:  res[3],
	}, nil
}
