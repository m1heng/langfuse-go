package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/m1heng/langfuse-go/model"
)

const (
	langfuseDefaultEndpoint = "https://cloud.langfuse.com"
)

type Client struct {
	resty *resty.Client
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

	var client *resty.Client
	if config.httpClient != nil {
		client = resty.NewWithClient(config.httpClient)
	} else {
		client = resty.New()
	}
	client.SetBaseURL(config.LangfuseHost).
		SetAuthScheme("Basic").
		SetAuthToken(basicAuth(config.PublicKey, config.SecretKey))

	return &Client{
		resty: client,
	}
}

func (c *Client) Ingestion(ctx context.Context, req *model.BatchIngestionRequest, res *model.IngestionResponse) error {
	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(res).
		Post("/api/public/ingestion")
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 300 {
		return errors.New(resp.Status())
	}

	return nil
}

// GetPromptRequest defines request for GetPrompt.
// ref: https://api.reference.langfuse.com/#get-/api/public/v2/prompts/-promptName-
func (c *Client) GetPrompt(req *model.GetPromptRequest) (*model.TextPrompt, *model.ChatPrompt, error) {
	rawJSON := map[string]interface{}{}
	errResp := map[string]interface{}{}
	r := c.resty.R().
		SetResult(&rawJSON).
		SetError(&errResp).
		SetPathParam("promptName", req.PromptName)
	if req.Version != nil {
		r.SetQueryParam("version", string(*req.Version))
	}
	if req.Label != nil {
		r.SetQueryParam("label", *req.Label)
	}

	resp, err := r.Get("/api/public/v2/prompts/{promptName}")

	if err != nil {
		return nil, nil, err
	}
	if rawJSON == nil {
		return nil, nil, errors.New("empty response")
	}

	if resp.StatusCode() != http.StatusOK {
		if errResp["error"] != nil {
			return nil, nil, errors.New(errResp["error"].(string))
		}
		return nil, nil, errors.New(resp.Status())
	}

	if rawJSON["type"] == "text" {
		textPrompt := &model.TextPrompt{}
		err = json.Unmarshal(resp.Body(), textPrompt)
		return textPrompt, nil, err
	} else if rawJSON["type"] == "chat" {
		chatPrompt := &model.ChatPrompt{}
		err = json.Unmarshal(resp.Body(), chatPrompt)
		return nil, chatPrompt, err
	}
	return nil, nil, errors.New("unknown prompt type")

}

func basicAuth(publicKey, secretKey string) string {
	auth := publicKey + ":" + secretKey
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
