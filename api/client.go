package api

import (
	"io"
	"net/http"
	"strings"
)

type Client interface {
	GetAvailableProfiles() (Profiles, error)
	CreateProject() (Project, error)
	UploadFiles(project Project, files []FileUpload) error
	StartEdit(project Project, params StartEditParams) error
	GetEditStatus(project Project) (editStatus, error)
	StartExport(project Project) error
	GetExportStatus(project Project) (exportStatus, error)
	GetExportDownloadLinks(project Project) ([]ExportDownload, error)
}

type ClientOption func(*client)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *client) {
		c.baseURL = baseURL
	}
}

func WithAuthToken(authToken string) ClientOption {
	return func(c *client) {
		c.authToken = authToken
	}
}

func WithHTTPClient(http http.Client) ClientOption {
	return func(c *client) {
		c.http = http
	}
}

func NewClient(opts ...ClientOption) Client {
	c := &client{
		baseURL: "https://api-beta.imagen-ai.com/v1",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type client struct {
	baseURL   string
	authToken string
	http      http.Client
}

func (t *client) urlForPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return t.baseURL + path
	}
	return t.baseURL + "/" + path
}

func (t *client) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, t.urlForPath(path), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", t.authToken)
	return req, nil
}

func (t *client) newGetRequest(path string) (*http.Request, error) {
	return t.newRequest(http.MethodGet, path, http.NoBody)
}

func (t *client) newPostRequest(path string, body io.Reader) (*http.Request, error) {
	return t.newRequest(http.MethodPost, path, body)
}
