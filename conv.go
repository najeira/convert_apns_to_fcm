package convert_apns_to_fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	endpoint = "https://iid.googleapis.com/iid/v1:batchImport"
)

type Input struct {
	ApiKey      string   `json:"-"`
	Application string   `json:"application,omitempty"`
	Sandbox     bool     `json:"sandbox,omitempty"`
	ApnsTokens  []string `json:"apns_tokens,omitempty"`
}

type Result struct {
	ApnsToken string `json:"apns_token,omitempty"`
	Status    string `json:"status,omitempty"`
	FcmToken  string `json:"registration_token,omitempty"`
}

type output struct {
	Results []Result `json:"results,omitempty"`
}

func Convert(input Input) ([]Result, error) {
	body, err := json.Marshal(&input)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", endpoint, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key=" + input.ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%s %s", res.Status, string(resBody))
	}

	var out output
	if err := json.Unmarshal(resBody, &out); err != nil {
		return nil, err
	}
	return out.Results, nil
}
