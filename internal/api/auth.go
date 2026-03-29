package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alex-305/ticktui/internal/config"
	"github.com/cli/browser"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	authURL     = "https://ticktick.com/oauth/authorize"
	tokenURL    = "https://ticktick.com/oauth/token"
	scope       = "tasks:write tasks:read"
	redirectURL = "http://localhost:8080"
)

func GetAuthURL(clientID string) string {
	return fmt.Sprintf("%s?scope=%s&client_id=%s&state=state&redirect_uri=%s&response_type=code",
		authURL, scope, clientID, redirectURL)
}

func GetAccessToken(clientID, clientSecret, code string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetBasicAuth(clientID, clientSecret).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":   "authorization_code",
			"code":         code,
			"redirect_uri": redirectURL,
		}).
		Post(tokenURL)

	if err != nil {
		return "", errors.Wrap(err, "requesting access token")
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", errors.Wrap(err, "parsing response")
	}

	return result.AccessToken, nil
}

func LaunchBrowserAndSaveAuthToken() error {

	err := godotenv.Load()

	if err != nil {
		return errors.Errorf("Failed to load environment variables")
	}

	clientIDEnv := "TICKTICK_CLIENT_ID"
	clientSecretEnv := "TICKTICK_CLIENT_SECRET"
	clientID := os.Getenv(clientIDEnv)
	clientSecret := os.Getenv(clientSecretEnv)

	if clientID == "" {
		return errors.Errorf("Missing required %s environment variable", clientIDEnv)
	}

	if clientSecret == "" {
		return errors.Errorf("Missing required %s environment variable", clientSecretEnv)
	}

	server := &http.Server{Addr: ":8080"}

	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		val := r.URL.Query().Get("code")
		if val != "" {
			fmt.Fprintf(w, "Auth successful! You can return to your terminal.")
			codeChan <- val
		}
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	authURL := GetAuthURL(clientID)
	if err := browser.OpenURL(authURL); err != nil {
		return errors.Errorf("Failed to open browser")
	}

	var authCode string
	select {
	case authCode = <-codeChan:
		_ = server.Close()
	case err := <-errChan:
		return errors.Errorf("Server error: %v", err)
	}

	token, err := GetAccessToken(clientID, clientSecret, authCode)
	if err != nil {
		return errors.Errorf("Failed to access token")
	}

	err = config.SaveToken(token)
	if err != nil {
		return errors.Errorf("Failed to save token")
	}

	return nil
}
