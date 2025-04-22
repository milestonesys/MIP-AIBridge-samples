package entities

import (
	"errors"
	"os"

	"connectivitysample/src/common"
	"connectivitysample/src/shared/utils"
)

const manufacturerNameEnvVar = "MANUFACTURER_NAME" //If set, it will be used as the manufacturer name when registering the app

type AppRegistration struct {
	registrationFile string // Path to the app registration file (register.graphql)
}

func NewAppRegistration(registrationFile string) (*AppRegistration, error) {
	if registrationFile == "" {
		return nil, errors.New("parameter 'registrationFile' can't be empty")
	}

	return &AppRegistration{
		registrationFile: registrationFile,
	}, nil
}

// Get the content of 'register.graphql' populated with the environment variables (and custom mappings)
func (ar *AppRegistration) GetPopulatedRegistrationFileContent() (string, error) {

	input, err := os.ReadFile(ar.registrationFile)
	if err != nil {
		return "", err
	}

	populatedRegistrationFileContent := utils.ExpandEnvWithDefault(string(input), appRegistrationMaping)

	return populatedRegistrationFileContent, nil
}

// Callback when applying the environment variables to the 'register.graphql' file
func appRegistrationMaping(key string) string {
	switch key {
	case manufacturerNameEnvVar:
		return getManufacturerName()
	case "APP_ID":
		return common.AppID
	case "APP_NAME":
		return common.AppName
	case "APP_DESCRIPTION":
		return common.AppDescription
	default:
		return os.Getenv(key)
	}
}

// The app's manufacturer (it's an OPTIONAL object when registering an app)
// Partners/Developers if you want to register your integration mind adjusting the following values:
// You can register your manufacturer name by sending an email to partner@milestone.dk
func getManufacturerName() string {
	manufacturerName := utils.GetEnv(manufacturerNameEnvVar, common.AppDefaultManufacturerName)

	// Even if the environment variable is set explicitly to 'empty' we will fallback to its default value.
	// AI Bridge will not accept an empty manufacturer name.
	if manufacturerName == "" {
		manufacturerName = common.AppDefaultManufacturerName
	}
	return manufacturerName
}
