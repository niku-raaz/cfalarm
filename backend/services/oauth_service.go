package services

import (
	"context"
	"cfalarm/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig = &oauth2.Config{
	ClientID:     config.GetEnv("GOOGLE_CLIENT_ID"),
	ClientSecret: config.GetEnv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  config.GetEnv("GOOGLE_REDIRECT_URL"),
	Scopes: []string{
	    "profile",
	    "email",
	    "https://www.googleapis.com/auth/calendar.events",
    },
	Endpoint:     google.Endpoint,
}

func ExchangeCode(code string) (*oauth2.Token, error) {
	return GoogleOAuthConfig.Exchange(context.Background(), code)
}
