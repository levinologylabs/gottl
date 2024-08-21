package providers

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleConfig struct {
	RedirectURL  string `json:"redirect_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (gc GoogleConfig) Scopes() []string {
	return []string{"openid profile email"}
}

func (gc GoogleConfig) Endpoint() oauth2.Endpoint {
	return google.Endpoint
}

func (gc GoogleConfig) OathConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     gc.ClientID,
		ClientSecret: gc.ClientSecret,
		RedirectURL:  gc.RedirectURL,
		Scopes:       gc.Scopes(),
		Endpoint:     gc.Endpoint(),
	}
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	PictureURL    string `json:"picture"`
}
