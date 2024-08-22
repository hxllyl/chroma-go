package cohere

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"

	httpc "github.com/amikos-tech/chroma-go/pkg/commons/http"
)

type APIVersion string

type CohereModel string // generic type for Cohere models

func (m CohereModel) String() string {
	return string(m)
}

const (
	APIKeyEnv                    = "COHERE_API_KEY"
	DefaultBaseURL               = "https://api.cohere.ai"
	APIVersionV1      APIVersion = "v1"
	DefaultAPIVersion            = APIVersionV1
	ClientName                   = "chroma-go-client"
)

// CohereClient is a common struct for various Cohere integrations - Embeddings, Rerank etc.
type CohereClient struct {
	BaseURL       string     `validate:"required"`
	APIVersion    APIVersion `validate:"required"`
	apiKey        string     `validate:"required"`
	Client        *http.Client
	DefaultModel  CohereModel `validate:"required"`
	RetryStrategy httpc.RetryStrategy
}

func NewCohereClient(opts ...Option) (*CohereClient, error) {
	client := &CohereClient{
		Client:     http.DefaultClient,
		BaseURL:    DefaultBaseURL,
		APIVersion: DefaultAPIVersion,
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}
	validate := validator.New(validator.WithRequiredStructEnabled(), validator.WithPrivateFieldValidation())
	err := validate.Struct(client)
	if err != nil {
		return nil, err
	}
	if client.RetryStrategy == nil {
		client.RetryStrategy, err = httpc.NewSimpleRetryStrategy(httpc.WithRetryableStatusCodes(429), httpc.WithExponentialBackOff())
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func (c *CohereClient) GetAPIEndpoint(endpoint string) string {
	return strings.ReplaceAll(fmt.Sprintf("%s/%s/%s", c.BaseURL, c.APIVersion, endpoint), "^[:]//", "/")
}

// TODO GetRequest is misleading, it should be renamed to GetHTTPRequest
func (c *CohereClient) GetRequest(ctx context.Context, method string, endpoint string, content string) (*http.Request, error) {
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewBufferString(content))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", ClientName)
	httpReq.Header.Set("X-Client-Name", ClientName)
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	return httpReq, nil
}

func (c *CohereClient) DoRequest(req *http.Request) (*http.Response, error) {
	if c.RetryStrategy != nil {
		return c.RetryStrategy.DoWithRetry(c.Client, req)
	} else {
		return c.Client.Do(req)
	}
}

type Option func(p *CohereClient) error

func NoOp() Option {
	return func(p *CohereClient) error {
		return nil
	}
}

func WithBaseURL(baseURL string) Option {
	return func(p *CohereClient) error {
		p.BaseURL = strings.TrimSuffix(baseURL, "/")
		return nil
	}
}

func WithAPIKey(apiKey string) Option {
	return func(p *CohereClient) error {
		p.apiKey = apiKey
		return nil
	}
}

func WithEnvAPIKey() Option {
	return func(p *CohereClient) error {
		if apiKey := os.Getenv(APIKeyEnv); apiKey != "" {
			p.apiKey = apiKey
			return nil
		}
		return fmt.Errorf("API key env variable %s not found or does not contain a key", APIKeyEnv)
	}
}

func WithAPIVersion(version APIVersion) Option {
	return func(p *CohereClient) error {
		if version == "" {
			return fmt.Errorf("API version can't be empty")
		}
		p.APIVersion = version
		return nil
	}
}

// WithHTTPClient sets the HTTP client for the Cohere client
func WithHTTPClient(client *http.Client) Option {
	return func(p *CohereClient) error {
		if client == nil {
			return fmt.Errorf("HTTP client is nil")
		}
		p.Client = client
		return nil
	}
}

// WithDefaultModel sets the default model for the Cohere client
func WithDefaultModel(model CohereModel) Option {
	return func(p *CohereClient) error {
		if model == "" {
			return fmt.Errorf("model can't be empty")
		}
		p.DefaultModel = model
		return nil
	}
}

// WithRetryStrategy sets the retry strategy for the Cohere client
func WithRetryStrategy(retryStrategy httpc.RetryStrategy) Option {
	return func(p *CohereClient) error {
		if retryStrategy == nil {
			return fmt.Errorf("retry strategy cannot be nil")
		}
		p.RetryStrategy = retryStrategy
		return nil
	}
}
