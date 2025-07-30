package services

import (
	entities "connectivitysample/src/domain/entities"
	"connectivitysample/src/infrastructure/repositories"
	"context"
)

// Service used to send GraphQL 'queries'+'mutations' to retrieve + modify data hosted on AI Bridge.
type GraphqlService struct {
	requestUrl string

	graphqlRepository *repositories.GraphqlRepository
}

func NewGraphqlService(requestUrl string) *GraphqlService {
	return &GraphqlService{
		requestUrl:        requestUrl,
		graphqlRepository: repositories.NewGraphqlRepository(),
	}
}

/* --- GraphQL queries ----------------------------------------------------- */

// Gets the list of identity providers registered in Milestone AI Bridge
func (gs *GraphqlService) GetIdentityProviderRegisteredIssuers(ctx context.Context) ([]string, error) {
	return gs.graphqlRepository.GetIdentityProviderRegisteredIssuers(ctx, gs.requestUrl)
}

// Gets the list of VMS ids registered in Milestone AI Bridge
func (gs *GraphqlService) GetVmsIds(ctx context.Context) ([]string, error) {
	return gs.graphqlRepository.GetVmsIds(ctx, gs.requestUrl)
}

// Gets the base64 encoded snapshot image for a given device and stream
func (gs *GraphqlService) GetSnapshot(ctx context.Context, deviceID, streamID, token string, commandLineParameters *entities.CommandLineParameters) (string, error) {
	return gs.graphqlRepository.GetSnapshot(ctx, gs.requestUrl, deviceID, streamID, token, commandLineParameters)
}

// Returns the REST endpoint for a given event topic name
func (gs *GraphqlService) GetRestEventTopicEndpoint(ctx context.Context, topicName string) (string, error) {
	return gs.graphqlRepository.GetRestEventTopicEndpoint(ctx, gs.requestUrl, topicName)
}

// Returns the REST endpoint for a given metadata topic name
func (gs *GraphqlService) GetRestMetadataTopicEndpoint(ctx context.Context, topicName string) (string, error) {
	return gs.graphqlRepository.GetRestMetadataTopicEndpoint(ctx, gs.requestUrl, topicName)
}

/* --- GraphQL mutations ----------------------------------------------------- */

// Register the app in AIBridge to all connected endpoints
func (gs *GraphqlService) RegisterConnectivitySample(ctx context.Context, populatedRegistrationFileContent, vmsId string) error {
	return gs.graphqlRepository.RegisterConnectivitySample(ctx, gs.requestUrl, populatedRegistrationFileContent, vmsId)
}
