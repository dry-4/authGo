package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthProvider interface {
	GetUserInfo(token *oauth2.Token) (OAuthUserInfo, error)
}

type OAuthUserInfo struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	ID        string `json:"sub"`
	AvatarURL string `json:"picture"`
}

type GoogleOAuthProvider struct {
	Config *oauth2.Config
}

func NewGoogleOAuthProvider() *GoogleOAuthProvider {
	return &GoogleOAuthProvider{
		Config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *GoogleOAuthProvider) GetUserInfo(token *oauth2.Token) (OAuthUserInfo, error) {
	client := g.Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return OAuthUserInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OAuthUserInfo{}, err
	}

	var userInfo OAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return OAuthUserInfo{}, err
	}

	return userInfo, nil
}

func (g *GoogleOAuthProvider) GenerateAuthURL(state string) string {
	return g.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (g *GoogleOAuthProvider) ExchangeCode(code string) (*oauth2.Token, error) {
	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}
	return token, nil
}
