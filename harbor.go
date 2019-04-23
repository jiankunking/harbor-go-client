package harbor

import (
	"net"
	"net/http"
	"time"
)

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty" json:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PageSize int `url:"page_size,omitempty" json:"page_size,omitempty"`
}
type Client struct {
	// HTTP client used to communicate with the API.
	Client  *http.Client
	BaseURL string
	//QueryData url.Values
	BasicAuth struct{ Username, Password string }
	// Services used for talking to different parts of the Harbor API.
	Projects     *ProjectsService
	Repositories *RepositoriesService
}

func (c *Client) SetBasicAuth(username string, password string) *Client {
	c.BasicAuth = struct{ Username, Password string }{username, password}
	return c
}

func NewClient(httpClient *http.Client, baseURL, username, password string) *Client {
	return newClient(httpClient, baseURL, username, password)
}

func newClient(httpClient *http.Client, baseURL, username, password string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 120 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 30 * time.Second,
				MaxIdleConnsPerHost:   20,
			},
		}
	}
	c := &Client{Client: httpClient}

	c.SetBasicAuth(username, password)
	c.BaseURL = baseURL
	// Create all the public services.
	c.Projects = &ProjectsService{client: c}
	c.Repositories = &RepositoriesService{client: c}
	return c
}
