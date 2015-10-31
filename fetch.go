package gcode

import (
	"net/http"
)

type GcodeService struct {
	client *http.Client
}

func New() *GcodeService {
	return NewFromClient(http.DefaultClient)
}

func NewFromClient(client *http.Client) *GcodeService {
	return &GcodeService{
		client: client,
	}
}

func (svc *GcodeService) Query(project string) *QueryService {
	return &QueryService{
		project: project,
		query:   nil,
		params:  map[string]string{"can": "open"},

		offset: 0,
		limit:  25,

		svc: svc,
	}
}

func (svc *GcodeService) Issue() *IssueService {
	return &IssueService{
		svc: svc,
	}
}
