package nexus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"salimon/tina-core/helpers"
)

func getNexusBaseURL() string {
	return os.Getenv("NEXUS_BASE_URL")
}

func sendHttpRequest(method string, path string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s", getNexusBaseURL(), path)
	accessToken, err := helpers.GenerateNexusAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("entity", os.Getenv("ENTITY_ID"))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Request failed with status code %d", res.StatusCode))
	}
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil

}
