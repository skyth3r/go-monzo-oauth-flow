package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func accounts(c *MonzoClient) error {
	path := "accounts"
	requestURL := fmt.Sprintf("https://%s/%s", apiHostname, path)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	if rsp.Body == nil {
		return fmt.Errorf("response body is empty")
	}

	rspJson := map[string]interface{}{}
	err = json.NewDecoder(rsp.Body).Decode(&rspJson)
	if err != nil {
		return err
	}

	accounts, ok := rspJson["accounts"].([]interface{})
	if !ok {
		return fmt.Errorf("cannot find accounts in response")
	}

	for _, account := range accounts {
		fmt.Println(account)
	}

	return nil
}
