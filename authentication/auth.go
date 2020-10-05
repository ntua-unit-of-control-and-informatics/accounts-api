package authentication

import (
	"context"
	"log"
	"os"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var Verifier *oidc.IDTokenVerifier
var oauth2Config oauth2.Config
var keyset oidc.KeySet

func Init() {

	log.Println("Starting oidc configuration")
	oidcProvider := os.Getenv("OIDC_PROVIDER")
	if oidcProvider == "" {
		// oidcProvider = "https://login.jaqpot.org/auth/realms/jaqpot"
		oidcProvider = "http://192.168.10.100:30100/auth/realms/jaqpot"
		// oidcProvider = "https://login.cloud.nanosolveit.eu/auth/realms/nanosolveit"
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, oidcProvider)
	if err != nil {
		panic(err)
		// handle error
	}

	if err != nil {
		panic(err)
	}

	clientID := "accounts-api"
	// clientSecret := "cbfd6e04-a51c-4982-a25b-7aaba4f30c81"

	// redirectURL := "http://localhost:8181/demo/callback"
	// Configure an OpenID Connect aware OAuth2 client.

	oauth2Config = oauth2.Config{
		ClientID: clientID,
		// ClientSecret: clientSecret,
		// RedirectURL:  redirectURL,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "openid"},
	}

	oidcConfig := &oidc.Config{
		ClientID:          clientID,
		SkipClientIDCheck: true,
		SkipIssuerCheck:   true,
	}
	keyset = oidc.NewRemoteKeySet(ctx, oidcProvider)

	Verifier = provider.Verifier(oidcConfig)

}

type OidcClaims struct {
	*jwt.Claims
	Name             string   `json:"name,omitEmpty"`
	PreferedUsername string   `json:"preferred_username"`
	GivenName        string   `json:"given_name"`
	FamilyName       string   `json:"family_name,omitEmpty"`
	Email            string   `json:"email"`
	Groups           []string `json:"groups"`
}

func GetClaims(token string) (claims OidcClaims, err error) {
	resultCl := OidcClaims{}
	tokenVer, err := Verifier.Verify(context.TODO(), token)
	erro := tokenVer.Claims(&resultCl)
	if erro != nil {
		log.Println("failed to parse Claims: ", erro.Error())
	}
	return resultCl, err
}

// create an instance of the CustomClaim
//   customClaims := CustomClaims{
// 	Claims: &jwt.Claims{
// 	 Issuer:   "issuer1",
// 	 Subject:  "subject1",
// 	 ID:       "id1",
// 	 Audience: jwt.Audience{"aud1", "aud2"},
// 	 IssuedAt:
// 	  jwt.NewNumericDate(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)),
// 	 Expiry:
// 	  jwt.NewNumericDate(time.Date(2017, 1, 1, 0, 8, 0, 0, time.UTC)),
// 	},
// 	PrivateClaim1: "val1",
// 	PrivateClaim2: []string{"val2", "val3"},
// 	AnyJSONObjectClaim: map[string]interface{}{
// 	 "name": "john",
// 	 "phones": map[string]string{
// 	  "phone1": "123",
// 	  "phone2": "456",
// 	 },
// 	},
//    }
//   // add claims to the Builder
//   builder = builder.Claims(customClaims)

// func GetClaims(token string){
// 	var userInfo oidc.UserInfo

// }

// func Verify(ctx context.Context, token string) (verifies bool, err error) {
// 	verifier.Verify(ctx, token)
// }
