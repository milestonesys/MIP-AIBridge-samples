package services

import (
	"embed"
	"log"
	"regexp"
	"strings"
	"time"
)

//go:embed templates/*
var embeddedFiles embed.FS

type FileReader struct {
	fileName   string
	fileFormat string
}

func NewFileReader(fileName string, fileFormat string) *FileReader {
	return &FileReader{
		fileName:   fileName,
		fileFormat: fileFormat,
	}
}

// Public methods

// Reads all the files that start with the specified filename and file format (for example onvif1-2.xml or onvif-frame1-2.xml)
// Having two or more files helps us being able to simulate movement with their data or create some variance
// The method returns an array with all the files' content on string format

func (fr *FileReader) ReadMultipleFiles() ([]string, error) {
	templatesDir := "templates"
	files, err := embeddedFiles.ReadDir(templatesDir)

	if err != nil {
		log.Println("Error reading files on directory "+templatesDir+": ", err)
		return nil, err
	}

	// We need to check which files we want to use from the templates folder
	// We are picking up all files that start with the filename and end with a number ('onvif1.xml' for example)
	// Follow this template to add new metadata files

	re := regexp.MustCompile(fr.fileName + `(\d+).` + fr.fileFormat)

	var stringFiles []string

	for _, file := range files {
		filename := file.Name()
		if !re.MatchString(filename) {
			continue
		}

		content, err := fr.readFile(templatesDir + "/" + filename)
		if err != nil {
			// log error message
			return nil, err
		}
		stringFiles = append(stringFiles, content)
	}

	return stringFiles, nil
}

// Reads a single file that uses the specified filename plus the file format, like analyticEvent.json

func (fr *FileReader) ReadSingleFile() (string, error) {
	content, err := fr.readFile("templates/" + fr.fileName + "." + fr.fileFormat)
	if err != nil {
		return "", err
	}

	// Convert the buffer to a string and return
	return string(content), nil
}

// Replaces the UtcTime and SourceStreamID for their current values
// In case custom files are used they should use ##CurrentTimestamp## and ##CameraStreamId##

func TreatFile(metadataContent string, sourceStreamID string) string {
	newUTC := strings.Replace(metadataContent, "##CurrentTimestamp##", string(time.Now().Format(time.RFC3339Nano)), -1)
	newStreamID := strings.Replace(newUTC, "##CameraStreamId##", sourceStreamID, -1)
	return newStreamID
}

// Replaces the UtcTime and SourceStreamID for their current values
// In case custom files are used for metadata they should use ##CurrentTimestamp## and ##CameraStreamId##

func TreatMetadataFile(metadataContent string, sourceStreamID string) string {
	newUTC := strings.Replace(metadataContent, "##CurrentTimestamp##", string(time.Now().Format(time.RFC3339Nano)), -1)
	newStreamID := strings.Replace(newUTC, "##CameraStreamId##", sourceStreamID, -1)
	return newStreamID
}

// Replaces the CameraID for its current value
// In case custom files are used for an analytic event they should use {{ .CameraID }}

func TreatEventFile(eventContent string, cameraID string) string {
	return strings.Replace(eventContent, "{{ .CameraID }}", cameraID, -1)
}

// Private methods

func (fr *FileReader) readFile(filePath string) (string, error) {
	content, err := embeddedFiles.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Convert the buffer to a string and return
	return string(content), nil
}
