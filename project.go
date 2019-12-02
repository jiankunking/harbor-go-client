package harbor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Project holds the details of a project.
type Project struct {
	ProjectID    int64             `json:"project_id"`
	OwnerID      int               `json:"owner_id"`
	Name         string            `json:"name"`
	CreationTime time.Time         `json:"creation_time"`
	UpdateTime   time.Time         `json:"update_time"`
	Deleted      interface{}       `json:"deleted"`
	OwnerName    string            `json:"owner_name"`
	Togglable    bool              `json:"togglable"`
	Role         int               `json:"current_user_role_id"`
	RepoCount    int64             `json:"repo_count"`
	Metadata     map[string]string `json:"metadata"`
}
type ListProjectsOptions struct {
	ListOptions
	Name   string `url:"name,omitempty" json:"name,omitempty"`
	Public bool   `url:"public,omitempty" json:"public,omitempty"`
	Owner  string `url:"owner,omitempty" json:"owner,omitempty"`
}

// ProjectsService handles communication with the user related methods of
// the Harbor API.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L45
type ProjectsService struct {
	client *Client
}

// Return specific project detail information.
//
// This endpoint returns specific project information by project ID.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L149
func (s *ProjectsService) GetProjectByID(pid int64) (*Project, *http.Response, error) {
	var project *Project
	bytes, resp, err := httpGet(s.client, nil, fmt.Sprintf("projects/%d", pid))
	if err != nil {
		return nil, resp, err
	}
	err = json.Unmarshal(bytes, &project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, err
}

// List projects
//
// This endpoint returns all projects created by Harbor,
// and can be filtered by project name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L46
func (s *ProjectsService) ListProject(opt *ListProjectsOptions) ([]*Project, *http.Response, error) {
	var projects []*Project
	queryData, err := Query(*opt)
	if err != nil {
		return nil, nil, err
	}
	bytes, resp, err := httpGet(s.client, queryData, "projects")
	if err != nil {
		return nil, resp, err
	}
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		return nil, resp, err
	}
	return projects, resp, nil
}
