package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	oauthHostname = "auth.monzo.com"
	responseType  = "code"
	grantType     = "authorization_code"
	port          = ":21234"
	redirectURI   = "http://127.0.0.1" + port + "/"
	callbackURI   = redirectURI + "callback"
)

func oauth(c *MonzoClient) error {
	// Generate a randomised state to protect the client from cross-site request forgery attacks
	u := uuid.New()
	state := u.String()

	login(c.id, state)
	callbackServer(c, state)
	err := exchangeCodeForToken(c, c.callbackCode)
	if err != nil {
		return err
	}

	return nil
}

func login(id, state string) {
	params := map[string]string{
		"client_id":     id,
		"redirect_uri":  callbackURI,
		"response_type": responseType,
		"state":         state,
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	loginURL := fmt.Sprintf("https://%s/?%s", oauthHostname, values.Encode())
	fmt.Printf("Please visit the following URL to authenticate: %s\n", loginURL)

	err := exec.Command("open", loginURL).Start()
	if err != nil {
		log.Fatal(err)
	}
}

func callbackServer(c *MonzoClient, state string) {
	server := &http.Server{Addr: port}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackURL := r.URL.Query()

		if _, ok := callbackURL["code"]; !ok {
			http.Error(w, "cannot find temporary auth code in callback URL", http.StatusBadRequest)
			return
		}

		if _, ok := callbackURL["state"]; !ok {
			http.Error(w, "cannot find randomised auth state in callback URL", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(callbackURL.Get("state")) != state {
			http.Error(w, "invalid randomised auth state in callback URL", http.StatusBadRequest)
			return
		}

		c.callbackCode = callbackURL.Get("code")

		w.Write([]byte("Authentication successful! You can now close this window."))

		// Shutdown the server
		go func() {
			// Wait for 2 seconds to allow the response to be sent
			time.Sleep(2 * time.Second)

			if err := server.Shutdown(context.Background()); err != nil {
				log.Fatalf("Server Shutdown Failed:%+v", err)
			}
		}()
	})

	// Start the server
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func exchangeCodeForToken(c *MonzoClient, code string) error {
	if code == "" {
		return fmt.Errorf("no auth code provided")
	}

	params := map[string]string{
		"grant_type":    grantType,
		"client_id":     c.id,
		"client_secret": c.secret,
		"redirect_uri":  callbackURI,
		"code":          code,
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	requestURL := fmt.Sprintf("https://%s/oauth2/token", apiHostname)

	encodedValues := values.Encode()

	req, err := http.NewRequest("POST", requestURL, strings.NewReader(encodedValues))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	rspJson := map[string]interface{}{}
	err = json.NewDecoder(rsp.Body).Decode(&rspJson)
	if err != nil {
		return err
	}

	accessToken, ok := rspJson["access_token"].(string)
	if !ok {
		return fmt.Errorf("cannot find access token in response")
	}

	refreshToken, ok := rspJson["refresh_token"].(string)
	if !ok {
		return fmt.Errorf("cannot find refresh token in response")
	}

	c.accessToken = accessToken
	c.refreshToken = refreshToken

	_, ok = rspJson["user_id"].(string)
	if !ok {
		return fmt.Errorf("cannot find user ID in response")
	}

	fmt.Println("Token exchange successful 🎉")

	return nil
}
