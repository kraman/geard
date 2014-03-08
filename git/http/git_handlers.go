package http

import (
	"fmt"
	"github.com/smarterclayton/geard/containers"
	"github.com/smarterclayton/geard/git"
	gitjobs "github.com/smarterclayton/geard/git/jobs"
	"github.com/smarterclayton/geard/http"
	"github.com/smarterclayton/geard/jobs"
	"github.com/smarterclayton/go-json-rest"
)

func Routes() []http.HttpJobHandler {
	return []http.HttpJobHandler{
		&httpCreateRepositoryRequest{},
		&httpGitArchiveContentRequest{Ref: "*"},
	}
}

type httpCreateRepositoryRequest gitjobs.CreateRepositoryRequest

func (h *httpCreateRepositoryRequest) HttpMethod() string { return "PUT" }
func (h *httpCreateRepositoryRequest) HttpPath() string   { return "/repository" }
func (h *httpCreateRepositoryRequest) Handler(conf *http.HttpConfiguration) http.JobHandler {
	return func(reqid jobs.RequestIdentifier, token *http.TokenData, r *rest.Request) (jobs.Job, error) {
		repositoryId, errg := containers.NewIdentifier(token.ResourceLocator())
		if errg != nil {
			return nil, errg
		}
		// TODO: convert token into a safe clone spec and commit hash
		return &gitjobs.CreateRepositoryRequest{
			git.RepoIdentifier(repositoryId),
			token.ResourceType(),
		}, nil
	}
}

type httpGitArchiveContentRequest gitjobs.GitArchiveContentRequest

func (h *httpGitArchiveContentRequest) HttpMethod() string { return "GET" }
func (h *httpGitArchiveContentRequest) HttpPath() string {
	return "/repository/archive/" + string(h.Ref)
}
func (h *httpGitArchiveContentRequest) Handler(conf *http.HttpConfiguration) http.JobHandler {
	return func(reqid jobs.RequestIdentifier, token *http.TokenData, r *rest.Request) (jobs.Job, error) {
		repoId, errr := containers.NewIdentifier(token.ResourceLocator())
		if errr != nil {
			return nil, jobs.SimpleJobError{jobs.JobResponseInvalidRequest, fmt.Sprintf("Invalid repository identifier: %s", errr.Error())}
		}
		ref, errc := gitjobs.NewGitCommitRef(r.PathParam("*"))
		if errc != nil {
			return nil, jobs.SimpleJobError{jobs.JobResponseInvalidRequest, fmt.Sprintf("Invalid commit ref: %s", errc.Error())}
		}

		return &gitjobs.GitArchiveContentRequest{
			git.RepoIdentifier(repoId),
			ref,
		}, nil
	}
}