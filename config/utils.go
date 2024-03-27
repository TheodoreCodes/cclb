package config

import (
	"encoding/json"
	"net/url"
)

type configUrl url.URL

func (rv *configUrl) UnmarshalJSON(data []byte) error {
	var rawUrl string

	err := json.Unmarshal(data, &rawUrl)
	if err != nil {
		return err
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	*rv = configUrl(*u)

	return nil
}

func (rv *configUrl) URL() *url.URL {
	u := url.URL(*rv)
	return &u
}
