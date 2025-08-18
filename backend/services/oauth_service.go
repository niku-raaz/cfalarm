package services

import (
	"cfalarm/config"
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: config.GetEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  config.GetEnv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"profile",
			"email",
			"https://www.googleapis.com/auth/calendar.events",
		},
		Endpoint: google.Endpoint,
	}
}

func ExchangeCode(code string) (*oauth2.Token, error) {
	return GetGoogleOAuthConfig().Exchange(context.Background(), code)
}