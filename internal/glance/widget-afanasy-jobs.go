package glance

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var afanasyJobsTemplate = mustParseTemplate("afanasy-jobs.html", "widget-base.html")

type afanasyJobsWidget struct {
	widgetBase    `yaml:",inline"`
	AllowInsecure bool         `yaml:"allow-insecure"`
	Jobs          []AfanasyJob `yaml:"-"`
}

type AfanasyJob struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	State      string  `json:"state"`
	UserName   string  `json:"user_name"`
	Percentage float64 `json:"progress_percentage"`
}

func (widget *afanasyJobsWidget) initialize() error {
	widget.withTitle("Afanasy Jobs").withCacheDuration(30 * time.Second)
	return nil
}

func (widget *afanasyJobsWidget) update(ctx context.Context) {
	url := "http://192.168.90.104:5000/afanasy/jobs"
	var client *http.Client
	if widget.AllowInsecure {
		client = defaultInsecureHTTPClient
	} else {
		client = defaultHTTPClient
	}
	resp, err := client.Get(url)
	if err != nil {
		widget.withError(fmt.Errorf("failed to fetch jobs from n8n: %w", err))
		return
	}
	defer resp.Body.Close()

	// Support JSON with a 'jobs' key
	type jobsResponse struct {
		Jobs []AfanasyJob `json:"jobs"`
	}
	var jobsData jobsResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&jobsData); err != nil {
		widget.withError(fmt.Errorf("failed to decode jobs JSON: %w", err))
		return
	}
	j := jobsData.Jobs

	fmt.Printf("afanasy-jobs widget: loaded %d jobs\n", len(j))
	widget.Jobs = j
}

func (widget *afanasyJobsWidget) Render() template.HTML {
	return widget.renderTemplate(widget, afanasyJobsTemplate)
}
