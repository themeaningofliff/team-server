package webservice

import (
	"log"
	"net/http"

	"golang.org/x/oauth2"
	oauth2ClientAPI "google.golang.org/api/oauth2/v2"
)

/*
	This implementation cheats and uses a Google API to validate the token.

	An easy way to validate an ID token for debugging and low-volume use is to use the tokeninfo endpoint (https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=XYZ123). Calling this endpoint involves an additional
	network request that does most of the validation for you, but introduces some latency and the potential for network errors.
*/
func verifyTokenAPI(idToken string) (*oauth2ClientAPI.Tokeninfo, error) {
	httpClient := &http.Client{}
	oauth2Service, err := oauth2ClientAPI.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
}

// Test function to wrap up whatever verify we decide to go with.
// Takes in the JWT ID token.
func verifyToken(rawIDToken string) {
	/*
		2018/01/27 17:58:49 Token :  eyJhbGciOiJSUzI1NiIsImtpZCI6IjI2YzAxOGIyMzNmZTJlZWY0N2ZlZGJiZGQ5Mzk4MTcwZmM5YjI5ZDgifQ.eyJhenAiOiI.....

		An ID Token is a JWT (JSON Web Token), that is, a cryptographically signed Base64-encoded JSON object. Normally, it is critical that you validate an ID token before you use it,
		but since you are communicating directly with Google over an intermediary-free HTTPS channel and using your client secret to authenticate yourself to Google, you can be confident that the
		token you receive really comes from Google and is valid. If your server passes the ID token to other components of your app, it is extremely important that the other components
		validate the token before using it.

		{"iss":"accounts.google.com",
		"at_hash":"HK6E_P6Dh8Y93mRNtsDB1Q",
		"email_verified":"true",
		"sub":"10769150350006150715113082367",
		"azp":"1234987819200.apps.googleusercontent.com",
		"email":"jsmith@example.com",
		"aud":"1234987819200.apps.googleusercontent.com",
		"iat":1353601026,
		"exp":1353604926,
		"nonce": "0394852-3190485-2490358",
		"hd":"example.com" }
	*/
	log.Println("JWT ID Token: ", string(rawIDToken))

	// try the OIDC verifier.
	idToken, err := oauthVerifier.Verify(oauth2.NoContext, rawIDToken)
	if err != nil {
		// http.Error(w, "Failed to verify OIDC ID Token: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to verify OIDC ID Token: " + err.Error())
		return
	}
	pprint(idToken, "OIDC id token")

	// try the GoogleAPI Token Verifier (under the hood this requires a http call to a Google Endpoint)
	tokenInfo, error := verifyTokenAPI(rawIDToken)
	if error != nil {
		// http.Error(w, "Failed to verify OIDC ID Token: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to verify GoogleAPI ID Token: " + err.Error())
		return
	}
	pprint(tokenInfo, "GoogleAPI id token")

	/*

		{
		    "Issuer": "accounts.google.com",
		    "Audience": [
		        "379625204959-4t2js39veijsiopjog6e2rtfruo0qrb3.apps.googleusercontent.com"
		    ],
		    "Subject": "100682826382643775970",
		    "Expiry": "2018-01-27T22:26:09-05:00",
		    "IssuedAt": "2018-01-27T21:26:09-05:00",
		    "Nonce": "",
		    "AccessTokenHash": "rldghiB4Walnp9odvCAjUQ"
		}
		2018/01/27 21:26:09 GoogleAPI id token
		 {
		    "audience": "379625204959-4t2js39veijsiopjog6e2rtfruo0qrb3.apps.googleusercontent.com",
		    "email": "omarhafezau@gmail.com",
		    "expires_in": 3599,
		    "issued_to": "379625204959-4t2js39veijsiopjog6e2rtfruo0qrb3.apps.googleusercontent.com",
		    "user_id": "100682826382643775970",
		    "verified_email": true
		}

	*/
}
