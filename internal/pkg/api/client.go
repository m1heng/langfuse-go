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
}

func New(config *ClientConfig) *Client {
	if config == nil {
		config = &ClientConfig{
			LangfuseHost: os.Getenv("LANGFUSE_HOST"),
			PublicKey:    os.Getenv("LANGFUSE_PUBLIC_KEY"),
			SecretKey:    os.Getenv("LANGFUSE_SECRET_KEY"),
		}
	}

	restClient := restclientgo.New(config.LangfuseHost)
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
