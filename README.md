# harbor-go-client
A Harbor API client enabling Go programs to interact with Harbor in a simple and uniform way

base TimeBye/go-harbor fix concurrent error

## Coverage

This API client package covers most of the existing Harbor API calls and is updated regularly
to add new and/or missing endpoints. Currently the following services are supported:

- [ ] Users
- [x] Projects
- [x] Repositories
- [ ] Jobs
- [ ] Policies
- [ ] Targets
- [ ] SystemInfo
- [ ] LDAP
- [ ] Configurations

## Usage

```go
import "github.com/jiankunking/harbor-go-client"
```

Construct a new Harbor client, then use the various services on the client to
access different parts of the Harbor API. For example, to list all
users:

```go
harborClient := harbor.NewClient(nil, "url","username","password")
```

Some API methods have optional parameters that can be passed. For example,
to list all projects for user "haobor":

```go
harborClient := harbor.NewClient(nil, "url","username","password")
opt := &ListProjectsOptions{Name: "haobor"}
projects, _, err := harborClient.Projects.ListProjects(opt)
``` 