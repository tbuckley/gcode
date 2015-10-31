package gcode

import (
	"encoding/xml"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type QueryService struct {
	project string
	query   []string
	params  map[string]string

	offset int
	limit  int

	svc *GcodeService
}

func (svc *QueryService) clone() *QueryService {
	query := make([]string, len(svc.query))
	for i, value := range svc.query {
		query[i] = value
	}

	params := make(map[string]string)
	for key, value := range svc.params {
		params[key] = value
	}

	return &QueryService{
		project: svc.project,
		query:   query,
		params:  params,
		offset:  svc.offset,
		limit:   svc.limit,
	}
}

func (svc *QueryService) Can(can string) *QueryService {
	clone := svc.clone()
	clone.params["can"] = can
	return clone
}

func (svc *QueryService) Open() *QueryService {
	return svc.Can("open")
}

func (svc *QueryService) All() *QueryService {
	return svc.Can("all")
}

func (svc *QueryService) Label(label string) *QueryService {
	clone := svc.clone()
	clone.params["label"] = label
	return clone
}

func (svc *QueryService) Query(query string) *QueryService {
	clone := svc.clone()
	clone.query = append(clone.query, query)
	return clone
}

func (svc *QueryService) addDateQuery(attribute string, date time.Time) *QueryService {
	dateString := date.Format("2006/01/02")
	query := attribute + ":" + dateString
	return svc.Query(query)
}

func (svc *QueryService) OpenedBefore(date time.Time) *QueryService {
	return svc.addDateQuery("opened-before", date)
}

func (svc *QueryService) OpenedAfter(date time.Time) *QueryService {
	return svc.addDateQuery("opened-after", date)
}

func (svc *QueryService) OpenedInRange(start time.Time, end time.Time) *QueryService {
	return svc.All().OpenedAfter(start).OpenedBefore(end)
}

func (svc *QueryService) ClosedBefore(date time.Time) *QueryService {
	return svc.addDateQuery("closed-before", date)
}

func (svc *QueryService) ClosedAfter(date time.Time) *QueryService {
	return svc.addDateQuery("closed-after", date)
}

func (svc *QueryService) ClosedInRange(start time.Time, end time.Time) *QueryService {
	return svc.All().ClosedAfter(start).ClosedBefore(end)
}

func (svc *QueryService) Offset(offset int) *QueryService {
	clone := svc.clone()
	clone.offset = offset
	return clone
}

func (svc *QueryService) Limit(limit int) *QueryService {
	clone := svc.clone()
	clone.limit = limit
	return clone
}

func (svc *QueryService) RemainingPages(totalPages int) []*QueryService {
	pages := make([]*QueryService, totalPages-1)
	for i := 1; i < totalPages; i++ {
		pages[i-1] = svc.Offset(svc.limit * i)
	}
	return pages
}

func (svc *QueryService) URL() string {
	values := url.Values{}
	for key, value := range svc.params {
		values.Set(key, value)
	}
	values.Set("max-results", strconv.Itoa(svc.limit))
	values.Set("start-index", strconv.Itoa(svc.offset+1))

	if len(svc.query) > 0 {
		values.Set("q", strings.Join(svc.query, " "))
	}

	u := url.URL{
		Scheme:   "https",
		Host:     "code.google.com",
		Path:     "/feeds/issues/p/" + svc.project + "/issues/full",
		RawQuery: values.Encode(),
	}
	return u.String()
}

func (svc *QueryService) Do() (*IssuesFeed, error) {
	client := svc.svc.client

	resp, err := client.Get(svc.URL())
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	feed := new(IssuesFeed)
	err = xml.Unmarshal(data, feed)
	return feed, err
}
