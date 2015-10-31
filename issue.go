package gcode

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type IssueService struct {
	svc *GcodeService
}

type IssueGetService struct {
	id  string
	svc *GcodeService
}

type IssuePutService struct {
	issue *Issue
	svc   *GcodeService
}

func (svc *IssueService) Get(ID string) *IssueGetService {
	return &IssueGetService{
		id:  ID,
		svc: svc.svc,
	}
}

func (svc *IssueGetService) URL() string {
	return fmt.Sprintf("https://code.google.com/feeds/issues/p/chromium/issues/full/%v", svc.id)
}

func (svc *IssueGetService) Do() (*Issue, error) {
	client := svc.svc.client

	resp, err := client.Get(svc.URL())
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	issue := new(Issue)
	err = xml.Unmarshal(data, issue)
	return issue, err
}

func (svc *IssueService) Put(issue *Issue) *IssuePutService {
	return &IssuePutService{
		issue: issue,
		svc:   svc.svc,
	}
}

func (svc *IssuePutService) Do() (*Issue, error) {
	return nil, nil
}
