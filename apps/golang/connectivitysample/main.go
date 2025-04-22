// This application  implements an example web server hosting:
// The web server communication is secured with HTTPS and only the
// Management Client will be allowed to connect. Also, access to the video is authorized using the credentials of the user by which you are logged in to
// the Management Client. The URL of the web server can be registered together with either an App or a topic. Selecting such an App or topic in the
// Management Client will then render the web page hosted by the web server.
// Pages (features showcased):
// 1 - 'snapshot': a webpage showing camera's live snapshot picture.. If selecting the App or topic in relation to a specific camera, the web page will render the a snapshot image of the live feed of that camera.
// 2 - 'analytic-events': a webpage allowing to start/stop sending 'fake' analytic events related to a certain camera. If selecting the App or topic in relation to a specific camera, the web page will render a button that can allow to start/stop the sending of analytic events.

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strconv"
	"sync"

	"connectivitysample/src/application/handlers"
	"connectivitysample/src/application/services"
	"connectivitysample/src/domain/entities"
	"connectivitysample/src/infrastructure/repositories"
	runinfo "connectivitysample/src/infrastructure/runninginfo"
)

// Repositories
var graphqlRepository *repositories.GraphqlRepository

// Services
var graphqlService *services.GraphqlService
var tokenService *services.TokenService
var queryStringService *services.QueryStringService
var analyticEventService *services.AnalyticEventService

// Handlers
var homeHandler *handlers.HomeHandler
var snapshotHandler *handlers.SnapshotHandler
var eventHandler *handlers.EventHandler

func main() {
	// Note, that if you run this web server inside the kubernetes cluster,
	// then TLS termination will be done by the ingress controller.
	var tlsCertFileArg = flag.String("tls-certificate-file", "certs/tls-server/server.crt", "Path to .crt file containing certificate of server (in PEM format)")
	var tlsKeyFileArg = flag.String("tls-key-file", "certs/tls-server/server.key", "Path to .key file containing private key of server (in PEM format)")
	var tlsEnabledArg = flag.Bool("tls-enabled", false, "If true, the service will be secured using TLS; both a certificate and key must be provided if enabled")
	var aibWebserviceLocationArg = flag.String("aib-webservice-location", "localhost:4000", "Location of the AIB web-service. It may be different if running directly on host or docker-compose/K8s cluster")
	var enforceOauthArg = flag.Bool("enforce-oauth", true, "MC requests will add the logged user's oauth token. This flag signals if the app should forward that token or use Milestone AI Bridge designated user instead.")
	var snapshotMaxWidthArg = flag.Int("snapshot-max-width", 300, "Max Width for the snapshot if not specified will 300 by default")
	var snapshotMaxHeightArg = flag.Int("snapshot-max-height", 300, "Max Height for the snapshot if not specified will 300 by default")
	var appRegistrationFilePathArg = flag.String("app-registration-file-path", "config/register.graphql", "graphql file containing the description and topics of connectivity sample to be registered within the system")
	var appWebserverPortArg = flag.Int("app-webserver-port", 7443, "The port in which the app can be reached")
	var appUrlPathArg = flag.String("app-url-path", "connectivitysample", "The url path set for the current application")

	flag.Parse()

	//Set configuration from command line arguments
	tlsConfiguration, err := entities.NewTlsConfiguration(*tlsEnabledArg, *tlsCertFileArg, *tlsKeyFileArg)
	if err != nil {
		log.Println(err)
		return
	}

	snapshotConfiguration := entities.NewSnapshotConfiguration(*snapshotMaxWidthArg, *snapshotMaxHeightArg)

	commandLineParameters := entities.NewCommandLineParameters(*aibWebserviceLocationArg, snapshotConfiguration, *enforceOauthArg,
		*appRegistrationFilePathArg, *appWebserverPortArg, *appUrlPathArg, tlsConfiguration)

	// Show information to user
	runinfo.Print("connectivity-sample")

	// This is the URL of the Milestone AI Bridge GraphQL web service. Using localhost like
	// here assumes that you are running the Milestone AI Bridge in debug mode. If running in
	// production inside docker compose or a kubernetes cluster, you have to use
	// the hostname of the web service container.
	queryURL := "http://" + *aibWebserviceLocationArg + "/api/bridge/graphql"

	graphqlRepository = repositories.NewGraphqlRepository()

	graphqlService = services.NewGraphqlService(queryURL)
	tokenService = services.NewTokenService(graphqlRepository, queryURL, commandLineParameters)
	queryStringService = services.NewQueryStringService()
	analyticEventService = services.NewAnalyticEventService(&sync.Map{}, graphqlService)

	homeHandler = handlers.NewHomeHandler()
	snapshotHandler = handlers.NewSnapshotHandler(tokenService, graphqlService, queryStringService, commandLineParameters)
	eventHandler = handlers.NewEventHandler(queryStringService, analyticEventService, commandLineParameters)

	http.HandleFunc("/", homeHandler.Handle)
	http.HandleFunc("/"+commandLineParameters.AppUrlPath()+"/snapshot/", snapshotHandler.Handle)
	http.HandleFunc("/"+commandLineParameters.AppUrlPath()+"/event/", eventHandler.Handle)
	http.HandleFunc("/"+commandLineParameters.AppUrlPath()+"/event/processing/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			eventHandler.ProcessingHandle(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Registering the application against all connected VMSs")
	vmsIds, err := graphqlService.GetVmsIds(context.Background())
	if err != nil {
		log.Printf("Error querying VMS Ids %v", err)
		return
	}

	appRegistration, err := entities.NewAppRegistration(*appRegistrationFilePathArg)
	if err != nil {
		log.Printf("Error creating AppRegistration: %v", err)
		return
	}

	populatedRegistrationFileContent, err := appRegistration.GetPopulatedRegistrationFileContent()
	if err != nil {
		log.Printf("Error when mapping the registration configuration file: %v", err)
		return
	}

	for _, vmsId := range vmsIds {
		err = graphqlService.RegisterConnectivitySample(context.Background(), populatedRegistrationFileContent, vmsId)
		if err != nil {
			log.Printf("Could not register the vms ID = %s Error message: %v ", vmsId, err)
		} else {
			log.Println("Registration succeeded")
		}
	}

	// Start the web server
	if tlsConfiguration.TlsEnabled() {
		//On docker compose based installations we're exposing the web server directly to the outside world. Thus, we need to use the certificate and key files
		err = http.ListenAndServeTLS(":"+strconv.Itoa(commandLineParameters.AppWebserverPort()), tlsConfiguration.TlsCertificateFile(), tlsConfiguration.TlsKeyFile(), nil)
	} else {
		//On k8s based installation we leverage the encryption provided by the ingress controller.
		err = http.ListenAndServe(":"+strconv.Itoa(commandLineParameters.AppWebserverPort()), nil)
	}

	if err != nil {
		log.Fatal("Error while starting the webserver: ", err)
		return
	}

}
