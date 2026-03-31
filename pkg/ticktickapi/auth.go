package ticktickapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alex-305/ticktui/internal/asciiart"
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

func getOAuthCredentials() (string, string, error) {
	_ = godotenv.Load()
	id := os.Getenv("TICKTICK_CLIENT_ID")
	secret := os.Getenv("TICKTICK_CLIENT_SECRET")

	if id == "" || secret == "" {
		return "", "", errors.New("missing TICKTICK_CLIENT_ID or TICKTICK_CLIENT_SECRET")
	}
	return id, secret, nil
}

func startCallbackServer(addr string) (*http.Server, <-chan string, <-chan error) {
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintf(w, asciiart.Logo+"\n\nSuccessfully authenticated. You can now return to the comfort of your terminal :)")
			codeChan <- code
		}
	})

	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	return server, codeChan, errChan
}

func LaunchBrowserAndSaveAuthToken() error {
	clientID, clientSecret, err := getOAuthCredentials()
	if err != nil {
		return err
	}

	server, codeChan, errChan := startCallbackServer(":8080")
	defer server.Close()

	if err := browser.OpenURL(GetAuthURL(clientID)); err != nil {
		return errors.Wrap(err, "failed to open browser")
	}

	var authCode string
	select {
	case authCode = <-codeChan:
	case err := <-errChan:
		return errors.Wrap(err, "server error")
	}

	token, err := GetAccessToken(clientID, clientSecret, authCode)
	if err != nil {
		return err
	}

	return config.SaveToken(token)
}
