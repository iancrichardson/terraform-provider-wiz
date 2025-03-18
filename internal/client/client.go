package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/machinebox/graphql"
)

// Config holds the configuration for the Wiz API client
type Config struct {
	ClientID     string
	ClientSecret string
	APIURL       string
	AuthURL      string
}

// Client is the Wiz API client
type Client struct {
	config        *Config
	graphqlClient *graphql.Client
	token         string
	tokenExpiry   time.Time
}

type accessToken struct {
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
}

// NewClient creates a new Wiz API client
func NewClient(config *Config) (*Client, error) {
	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("client_id and client_secret are required")
	}

	if config.APIURL == "" {
		config.APIURL = "https://api.eu1.demo.wiz.io/graphql"
	}

	if config.AuthURL == "" {
		config.AuthURL = "https://auth.demo.wiz.io/oauth/token"
	}

	graphqlClient := graphql.NewClient(config.APIURL)

	return &Client{
		config:        config,
		graphqlClient: graphqlClient,
	}, nil
}

// authenticate gets a new access token
func (c *Client) authenticate() error {
	// Check if token is still valid
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return nil
	}

	authData := url.Values{}
	authData.Set("grant_type", "client_credentials")
	authData.Set("audience", "wiz-api")
	authData.Set("client_id", c.config.ClientID)
	authData.Set("client_secret", c.config.ClientSecret)

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, c.config.AuthURL, strings.NewReader(authData.Encode()))
	if err != nil {
		return fmt.Errorf("error creating authentication request: %w", err)
	}

	req.Header.Add("Encoding", "UTF-8")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error authenticating: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error authenticating, status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading authentication response: %w", err)
	}

	var at accessToken
	if err := json.Unmarshal(bodyBytes, &at); err != nil {
		return fmt.Errorf("error parsing authentication response: %w", err)
	}

	c.token = at.Token
	c.tokenExpiry = time.Now().Add(time.Duration(at.Expires) * time.Second)

	return nil
}

// RunQuery executes a GraphQL query
func (c *Client) RunQuery(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	if err := c.authenticate(); err != nil {
		return err
	}

	req := graphql.NewRequest(query)

	// Set variables
	for k, v := range variables {
		req.Var(k, v)
	}

	// Set auth header
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Run the query
	if err := c.graphqlClient.Run(ctx, req, response); err != nil {
		return err
	}

	return nil
}
