package entities

import (
	"errors"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Extends the jwt.Claims interface.
// Any TokenClaim should implement the jwt.Claims interface and expose the token ClientID and Issuer
type TokenClaims interface {
	// TokenClaims interface is jwt.Claims with some extended methods
	// The token will be parsed to obtain the ClientID and TokenIssuer to validate the token against the required fields
	jwt.Claims

	// Get the token client id
	GetClientID() string
	// Get the token issuer
	GetTokenIssuer() string
}

// Concrete implementation of the TokenClaims interface
type tokenClaims struct {
	// Json fields to parsed
	ClientID    string `json:"client_id"`
	TokenIssuer string `json:"iss"`

	// Mandatory fields.

	// Used to select the proper issuer to validate the token.
	registeredIssuers []string
	// Expected ClientID that should be set in the incoming oauth token.
	requiredClientID string
}

// Constructor
func NewTokenClaims(registeredIssuers []string, requiredClientID string) TokenClaims {
	return &tokenClaims{
		registeredIssuers: registeredIssuers,
		requiredClientID:  requiredClientID,
	}
}

func (tc *tokenClaims) Valid() error {
	// Verify that client id of token is as required
	if tc.ClientID != tc.requiredClientID {
		return errors.New("token is not issued to a trusted client")
	}

	// Parse token issuer url
	tkIssuerURL, err := url.Parse(tc.TokenIssuer)
	if err != nil {
		return err
	}

	// Verify that token issuer URL is absolute
	if !tkIssuerURL.IsAbs() {
		return errors.New("token issuer url is not absolute")
	}

	// Verify that token issuer URL resolves to one of the registered issuers
	for _, rgIssuer := range tc.registeredIssuers {
		// Parse issuer url
		rgIssuerURL, err := url.Parse(rgIssuer)
		if err != nil {
			continue
		}

		// Compare by hostname
		if strings.EqualFold(tkIssuerURL.Hostname(), rgIssuerURL.Hostname()) {
			return nil
		}
	}

	// If we get here, the token issuer did not match one of the registered issuers
	return errors.New("token issuer not recognized")
}

func (tc *tokenClaims) GetClientID() string {
	return tc.ClientID
}

func (tc *tokenClaims) GetTokenIssuer() string {
	return tc.TokenIssuer
}
