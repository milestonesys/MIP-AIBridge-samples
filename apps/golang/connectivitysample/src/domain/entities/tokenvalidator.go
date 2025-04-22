package entities

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
)

// Validates whether an oauth token is valid or not
type TokenValidator struct {
	bearerToken string
	tokenClaims TokenClaims
}

// Constructor
func NewTokenValidator(bearerToken string, tokenClaims TokenClaims) *TokenValidator {

	return &TokenValidator{
		bearerToken: bearerToken,
		tokenClaims: tokenClaims,
	}
}

// Given a token and tokenClaims will verify the bearerToken has the required issuers and clientID
func (tv *TokenValidator) IsValid(ctx context.Context) error {
	// Verify client id and issuer of token
	if err := tv.verifyTokenClaims(); err != nil {
		return errors.New("error at verifyToken method, during IsValid execution:" + err.Error())
	}

	// Verify signature of token using public key from token issuer
	if err := tv.verifyTokenSignature(ctx, tv.tokenClaims.GetTokenIssuer()); err != nil {
		return errors.New("error at verifyToken method, during verifyTokenSignature execution:" + err.Error())
	}

	// Success
	return nil
}

// Parse the issuers and clientID from the token and validate they match the required ones from the tokenClaim
func (tv *TokenValidator) verifyTokenClaims() error {
	// Parse token and verify specified claims
	parser := jwt.Parser{}
	if _, _, err := parser.ParseUnverified(tv.bearerToken, tv.tokenClaims); err != nil {
		return errors.New("error at verifyTokenClaims method, during ParseUnverified execution:" + err.Error())
	}
	// Verify the token has the required issuers and clientID defined by the tokenClaim
	return tv.tokenClaims.Valid()
}

// Verifies the token has been signed by the provider
func (tv *TokenValidator) verifyTokenSignature(ctx context.Context, issuer string) error {
	// Initialize provider instance of OpenID server
	client := &http.Client{Timeout: 10 * time.Second}
	clictx := oidc.ClientContext(ctx, client)
	provider, err := oidc.NewProvider(clictx, issuer)
	if err != nil {
		return errors.New("error at verifyTokenSignature method, during NewProvider execution:" + err.Error())
	}

	// Verify signature of token using public key from provider
	verifier := provider.Verifier(&oidc.Config{
		SkipClientIDCheck:    true,
		SkipExpiryCheck:      false,
		SkipIssuerCheck:      false,
		SupportedSigningAlgs: []string{"RS256"},
	})
	_, err = verifier.Verify(ctx, tv.bearerToken)
	if err != nil {
		return errors.New("error at verifyTokenSignature method, during Verify execution:" + err.Error())
	}
	return nil
}
