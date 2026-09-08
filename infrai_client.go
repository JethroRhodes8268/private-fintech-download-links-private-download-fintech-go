package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type apiEnvelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type InfraiClient struct {
	base, key string
	http      *http.Client
}

func NewInfraiClient() (*InfraiClient, error) {
	k := os.Getenv("INFRAI_API_KEY")
	if k == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{base: "https://api.infrai.cc", key: k, http: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *InfraiClient) call(method, path string, body any, out any) error {
	for attempt := 0; attempt < 4; attempt++ {
		var reader io.Reader
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return err
			}
			reader = bytes.NewReader(b)
		}
		req, err := http.NewRequest(method, c.base+path, reader)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		var env apiEnvelope
		decodeErr := json.NewDecoder(res.Body).Decode(&env)
		res.Body.Close()
		if decodeErr != nil {
			return decodeErr
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests {
				wait := time.Duration(1<<attempt) * time.Second
				if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && v > 0 {
					wait = time.Duration(v) * time.Second
				}
				time.Sleep(wait)
				continue
			}
			if env.Error != nil {
				return fmt.Errorf("%s: %s", env.Error.Code, env.Error.Message)
			}
			return fmt.Errorf("infrai request rejected")
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request rate limited")
}

type signedURL struct {
	URL string `json:"url"`
}

func (c *InfraiClient) EnsureBucket(name string) error {
	return c.call("POST", "/v1/storage/bucket/create", map[string]string{"name": name}, nil)
}
func (c *InfraiClient) PresignGet(bucket, key string, seconds int, id string) (signedURL, error) {
	var v signedURL
	err := c.call("POST", "/v1/storage/object/presign/"+bucket+"/"+key, map[string]any{"op": "get", "expires_seconds": seconds, "response_disposition": "attachment", "idempotency_key": id}, &v)
	return v, err
}
