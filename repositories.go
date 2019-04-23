package harbor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ListRepositoriesOption struct {
	ListOptions
	ProjectId int64  `url:"project_id,omitempty" json:"project_id,omitempty"`
	Q         string `url:"q,omitempty" json:"q,omitempty"`
}

type RepoResp struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ProjectID    int64     `json:"project_id"`
	Description  string    `json:"description"`
	PullCount    int64     `json:"pull_count"`
	StarCount    int64     `json:"star_count"`
	TagsCount    int64     `json:"tags_count"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

// RepoRecord holds the record of an repository in DB, all the infors are from the registry notification event.
type RepoRecord struct {
	RepositoryID int64     `json:"repository_id"`
	Name         string    `json:"name"`
	ProjectID    int64     `json:"project_id"`
	Description  string    `json:"description"`
	PullCount    int64     `json:"pull_count"`
	StarCount    int64     `json:"star_count"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type cfg struct {
	Labels map[string]string `json:"labels"`
}

//ComponentsOverview has the total number and a list of components number of different serverity level.
type ComponentsOverview struct {
	Total   int                        `json:"total"`
	Summary []*ComponentsOverviewEntry `json:"summary"`
}

//ComponentsOverviewEntry ...
type ComponentsOverviewEntry struct {
	Sev   int `json:"severity"`
	Count int `json:"count"`
}

//ImgScanOverview mapped to a record of image scan overview.
type ImgScanOverview struct {
	ID              int64               `json:"-"`
	Digest          string              `json:"image_digest"`
	Status          string              `json:"scan_status"`
	JobID           int64               `json:"job_id"`
	Sev             int                 `json:"severity"`
	CompOverviewStr string              `json:"-"`
	CompOverview    *ComponentsOverview `json:"components,omitempty"`
	DetailsKey      string              `json:"details_key"`
	CreationTime    time.Time           `json:"creation_time,omitempty"`
	UpdateTime      time.Time           `json:"update_time,omitempty"`
}

type TagDetail struct {
	Digest        string    `json:"digest"`
	Name          string    `json:"name"`
	Size          int64     `json:"size"`
	Architecture  string    `json:"architecture"`
	OS            string    `json:"os"`
	DockerVersion string    `json:"docker_version"`
	Author        string    `json:"author"`
	Created       time.Time `json:"created"`
	Config        *cfg      `json:"config"`
}

type Signature struct {
	Tag    string            `json:"tag"`
	Hashes map[string][]byte `json:"hashes"`
}

type TagResp struct {
	TagDetail
	Signature    *Signature       `json:"signature"`
	ScanOverview *ImgScanOverview `json:"scan_overview,omitempty"`
}

// RepositoriesService handles communication with the user related methods of
// the Harbor API.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L891
type RepositoriesService struct {
	client *Client
}

// Get repositories accompany with relevant project and repo name.
//
// This endpoint let user search repositories accompanying
// with relevant project ID and repo name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L892
func (s *RepositoriesService) ListRepository(opt *ListRepositoriesOption) ([]*RepoRecord, *http.Response, error) {
	var v []*RepoRecord
	queryData, err := Query(*opt)
	if err != nil {
		return nil, nil, err
	}
	bytes, resp, err := httpGet(s.client, queryData, "repositories")
	if err != nil {
		return nil, resp, err
	}
	err = json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// Get tags of a relevant repository.
//
// This endpoint aims to retrieve tags from a relevant repository.
// If deployed with Notary, the signature property of
// response represents whether the image is singed or not.
// If the property is null, the image is unsigned.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1054
func (s *RepositoriesService) ListRepositoryTags(repoName string) ([]*TagResp, *http.Response, error) {
	var v []*TagResp
	bytes, resp, err := httpGet(s.client, nil, fmt.Sprintf("repositories/%s/tags", repoName))
	if err != nil {
		return nil, resp, err
	}
	err = json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// Get the tag of the repository.
//
// This endpoint aims to retrieve the tag of the repository.
// If deployed with Notary, the signature property of
// response represents whether the image is singed or not.
// If the property is null, the image is unsigned.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L999
func (s *RepositoriesService) GetRepositoryTag(repoName, tag string) (*TagDetail, *http.Response, error) {
	var v TagDetail
	bytes, resp, err := httpGet(s.client, nil, fmt.Sprintf("repositories/%s/tags/%s", repoName, tag))
	if err != nil {
		return nil, resp, err
	}
	err = json.Unmarshal(bytes, &v)
	return &v, resp, err
}

// Delete a tag in a repository.
//
// This endpoint let user delete tags with repo name and tag.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1025
func (s *RepositoriesService) DeleteRepositoryTag(repoName, tag string) (*http.Response, error) {
	resp, err := httpDelete(s.client, fmt.Sprintf("repositories/%s/tags/%s", repoName, tag))
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Delete a repository.
//
// This endpoint let user delete a repository with name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L944
func (s *RepositoriesService) DeleteRepository(repoName string) (*http.Response, error) {
	resp, err := httpDelete(s.client, fmt.Sprintf("repositories/%s", repoName))
	if err != nil {
		return resp, err
	}
	return resp, nil
}
