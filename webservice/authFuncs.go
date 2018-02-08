package webservice

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/context"
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
		log.Println("Failed to verify GoogleAPI ID Token: " + err.Error())
		return nil, err
	}

	return tokenInfo, nil
}

// Test function to wrap up whatever verify we decide to go with.
// Takes in the JWT ID token.
func verifyToken(rawIDToken string) {
	/*
		2018/01/27 17:58:49 Token :  eyJhbGciOiJSUzI1NiIsImtpZCI6IjI2YzAxOGIyMzNmZTJlZWY0N2ZlZGJiZGQ5Mzk4MTcwZmM5YjI5ZDgifQ.eyJhenAiOiI.....

		// https://stackoverflow.com/questions/8311836/how-to-identify-a-google-oauth2-user/13016081#13016081
		An ID Token is a JWT (JSON Web Token) part of OpenID Connect, that is, a cryptographically signed Base64-encoded JSON object. Normally, it is critical that you validate an ID token before you use it,
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
	/*
		2018/01/28 11:10:22 OIDC id token
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
	*/
	pprint(idToken, "OIDC id token")

	// get all the "claims" - additional info - from the token.
	var newClaims = struct {
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{new(json.RawMessage)}

	// get just the email related claims.
	// var claims struct {
	// 	Email         string `json:"email"`
	// 	EmailVerified bool   `json:"email_verified"`
	// }
	// if err := idToken.Claims(&claims); err != nil {
	// 	log.Println("Failed to get claims off OIDC ID Token: " + err.Error())
	// 	return
	// }

	if err := idToken.Claims(&newClaims.IDTokenClaims); err != nil {
		log.Println("Failed to get claims off OIDC ID Token: " + err.Error())
		return
	}

	/*
		2018/01/31 07:54:33 IDTokenClaims:
		{
			"azp": "379625204959-4t2js39veijsiopjog6e2rtfruo0qrb3.apps.googleusercontent.com",
			"aud": "379625204959-4t2js39veijsiopjog6e2rtfruo0qrb3.apps.googleusercontent.com",
			"sub": "100682826382643775970",
			"email": "omarhafezau@gmail.com",
			"email_verified": true,
			"at_hash": "4CJLvsaQNF-MiNwRRxZnHQ",
			"exp": 1517406869,
			"iss": "accounts.google.com",
			"iat": 1517403269
		}
	*/
	pprint(newClaims.IDTokenClaims, "IDTokenClaims: ")

	// try the GoogleAPI Token Verifier (under the hood this requires a http call to a Google Endpoint)
	tokenInfo, error := verifyTokenAPI(rawIDToken)
	if error != nil {
		// http.Error(w, "Failed to verify OIDC ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	/*
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
	pprint(tokenInfo, "GoogleAPI id token")

}

// ValidateHandler wraps a standard http handler in an auth check to make sure they have a valid token.
func ValidateHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				// try the GoogleAPI Token Verifier (under the hood this requires a http call to a Google Endpoint)
				// according to https://developers.google.com/identity/sign-in/web/backend-auth, there is nothing wrong with using the id token to auth the user.
				tokenInfo, err := verifyTokenAPI(bearerToken[1])
				if err != nil {
					log.Println("Invalid authorization token: " + err.Error())
					http.Error(w, "Invalid authorization token", http.StatusInternalServerError)
					return
				}

				// set the necessary details on the context for use in the next handler.
				context.Set(req, "user_id", tokenInfo.UserId)               // "user_id": "100682826382643775970",
				context.Set(req, "email", tokenInfo.Email)                  // "email": "omarhafezau@gmail.com",
				context.Set(req, "verified_email", tokenInfo.VerifiedEmail) // "verified_email": true

				next(w, req)
				// token, error := jwt.Parse(bearerToken[1],
				// 	func(token *jwt.Token) (interface{}, error) {
				// 		pprint(token, "Bearer Token")
				// 		/*
				// 					"Name": "RS256",
				// 					"Hash": 5
				// 				},
				// 				"Header": {
				// 					"alg": "RS256",
				// 					"kid": "26c018b233fe2eef47fedbbdd9398170fc9b29d8"
				// 				},
				// 				"Claims": {
				// 					"at_hash": "UFbG3Mxnpb98V0e7BW54lg",
				// 					"aud": "379625204959-4t2js39veijsiopjog6e2rtfruo0qrb3.apps.googleusercontent.com",
				// 					"azp": "379625204959-4t2js39veijsiopjog6e2rtfruo0qrb3.apps.googleusercontent.com",
				// 					"email": "omarhafezau@gmail.com",
				// 					"email_verified": true,
				// 					"exp": 1517406187,
				// 					"iat": 1517402587,
				// 					"iss": "accounts.google.com",
				// 					"sub": "100682826382643775970"
				// 				},
				// 				"Signature": "",
				// 				"Valid": false
				// 			}
				// 		*/

				// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// 			return nil, fmt.Errorf("There was an error")
				// 		}

				// 		return []byte("secret"), nil
				// 	})
				// if error != nil {
				// 	json.NewEncoder(w).Encode(Exception{Message: error.Error()})
				// 	return
				// }

				// if token.Valid {
				// 	context.Set(req, "decoded", token.Claims)
				// 	next(w, req)
				// } else {
				// 	json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				// }
			}
		} else {
			// TODO: should probably re-direct to login page.
			LoginAgain(w, req)
			// json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}
