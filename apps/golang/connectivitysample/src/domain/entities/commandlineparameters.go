package entities

import "sync"

// Represents the command line parameters passed to the application
type CommandLineParameters struct {
	aibWebserviceLocation   string                // Location of the AIB web-service. It may be different if running directly on host or docker-compose/K8s cluster
	snapshotConfiguration   SnapshotConfiguration // Snapshot related configuration for the application
	enforceOauth            bool                  // Whether we want to enforce that the snapshot rendered in MC is fetched using the MC logged in user. If set to false the user defined by Milestone AI Bridge user will be used instead.
	appRegistrationFilePath string                // Graphql file containing the description and topics of 'connectivity sample' app to be registered within the system
	appWebserverPort        int                   // Port used by the application web server
	appUrlPath              string                // The url path set for the current application
	tlsConfiguration        TlsConfiguration      // TLS related configuration for the application
}

var (
	singleInstance *CommandLineParameters
	once           sync.Once
)

// Singleton pattern to create a single instance of CommandLineParameters
func NewCommandLineParameters(aibWebserviceLocation string, snapshotConfiguration SnapshotConfiguration, enforceOauth bool, appRegistrationFilePath string,
	appWebserverPort int, appUrlPath string, tlsConfiguration TlsConfiguration) *CommandLineParameters {

	once.Do(func() {
		singleInstance = &CommandLineParameters{
			aibWebserviceLocation:   aibWebserviceLocation,
			snapshotConfiguration:   snapshotConfiguration,
			enforceOauth:            enforceOauth,
			appRegistrationFilePath: appRegistrationFilePath,
			appWebserverPort:        appWebserverPort,
			appUrlPath:              appUrlPath,
			tlsConfiguration:        tlsConfiguration,
		}
	})
	return singleInstance
}

func (c *CommandLineParameters) AibWebserviceLocation() string {
	return c.aibWebserviceLocation
}

func (c *CommandLineParameters) SnapshotConfiguration() SnapshotConfiguration {
	return c.snapshotConfiguration
}

func (c *CommandLineParameters) EnforceOauth() bool {
	return c.enforceOauth
}

func (c *CommandLineParameters) AppRegistrationFilePath() string {
	return c.appRegistrationFilePath
}

func (c *CommandLineParameters) AppWebserverPort() int {
	return c.appWebserverPort
}

func (c *CommandLineParameters) AppUrlPath() string {
	return c.appUrlPath
}

func (c *CommandLineParameters) TlsConfiguration() TlsConfiguration {
	return c.tlsConfiguration
}
