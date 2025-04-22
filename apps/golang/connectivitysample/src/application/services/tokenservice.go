package services

import (
	"connectivitysample/src/common"
	entities "connectivitysample/src/domain/entities"
	"connectivitysample/src/infrastructure/repositories"
	"context"
	"errors"
	"net/http"
	"strings"
)

// Service used to deal with Oauth tokens
type TokenService struct {
	commandLineParameters *entities.CommandLineParameters
	requestUrl            string
	graphqlRepository     *repositories.GraphqlRepository
}

func NewTokenService(graphqlRepository *repositories.GraphqlRepository, requestUrl string, commandLineParameters *entities.CommandLineParameters) *TokenService {
	return &TokenService{
		commandLineParameters: commandLineParameters,
		requestUrl:            requestUrl,
		graphqlRepository:     graphqlRepository,
	}
}

// Extract the bearer token from the http request and validates if it's a valid token (according to the VMS' IDP)
func (ts *TokenService) ExtractAndVerifyToken(r *http.Request, w http.ResponseWriter) (string, error) {
	bearerToken := ""
	// Extract token from request headers
	auth, found := r.Header["Authorization"]

	if ts.commandLineParameters.EnforceOauth() {
		if !found || len(auth) != 1 || !strings.HasPrefix(auth[0], "Bearer ") {
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			return "", errors.New("error at extractAndVerifyToken method, Bearer token has not been passed")
		}

		// Lookup registered issuers
		registeredIssuers, err := ts.graphqlRepository.GetIdentityProviderRegisteredIssuers(context.Background(), ts.requestUrl)
		if err != nil {
			return "", errors.New("error at extractAndVerifyToken method, during RegisteredIssuers execution:" + err.Error())
		}

		// Verify that token is valid
		claims := entities.NewTokenClaims(registeredIssuers, common.ManagementClient_Oauth_ClientID)
		bearerToken = auth[0][7:]

		tokenValidator := entities.NewTokenValidator(bearerToken, claims)
		if err := tokenValidator.IsValid(context.Background()); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return "", errors.New("error at extractAndVerifyToken method, during verifyToken execution: " + err.Error())
		}
	}
	return bearerToken, nil
}
