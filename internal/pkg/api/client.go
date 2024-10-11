package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/henomis/restclientgo"
)

const (
	langfuseDefaultEndpoint = "https://cloud.langfuse.com"
)

type Client struct {
	restClient *restclientgo.RestClient
}

type ClientConfig struct {
	LangfuseHost string
	PublicKey    string
	SecretKey    string
	httpClient   *http.Client
}

func New(config *ClientConfig) *Client {
	if config == nil {
		config = &ClientConfig{}
	}

	if config.LangfuseHost == "" {
		if os.Getenv("LANGFUSE_HOST") != "" {
			config.LangfuseHost = os.Getenv("LANGFUSE_HOST")
		} else {
			config.LangfuseHost = langfuseDefaultEndpoint
		}
	}
	if config.PublicKey == "" {
		config.PublicKey = os.Getenv("LANGFUSE_PUBLIC_KEY")
	}
	if config.SecretKey == "" {
		config.SecretKey = os.Getenv("LANGFUSE_SECRET_KEY")
	}

	restClient := restclientgo.New(config.LangfuseHost)
	if config.httpClient != nil {
		restClient.SetHTTPClient(config.httpClient)
	}
	restClient.SetRequestModifier(func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", basicAuth(config.PublicKey, config.SecretKey))
		return req
	})

	return &Client{
		restClient: restClient,
	}
}

func (c *Client) Ingestion(ctx context.Context, req *Ingestion, res *IngestionResponse) error {
	return c.restClient.Post(ctx, req, res)
}

func basicAuth(publicKey, secretKey string) string {
	auth := publicKey + ":" + secretKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
