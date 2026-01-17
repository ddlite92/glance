package glance

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

// Response struct for jobs API
type RenderJobResponse struct {
	Jobs []RenderJob `json:"jobs"`
}

var renderJobWidgetTemplate = mustParseTemplate("render-job.html", "widget-base.html")

// Job represents a single job entry from the API
// Adjust fields as needed to match the actual API response
// Example fields: ID, Name, Status, Progress, etc.
type RenderJob struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
}

type renderJobWidget struct {
	widgetBase `yaml:",inline"`
	Jobs       []RenderJob `yaml:"-"`
	APIUrl     string      `yaml:"api-url"`
}

func (widget *renderJobWidget) initialize() error {
	widget.withTitle("Render Jobs").withCacheDuration(1 * time.Minute)
	if widget.APIUrl == "" {
		widget.APIUrl = "http://192.168.90.104:5000/afanasy/jobs"
	}
	return nil
}

func (widget *renderJobWidget) update(ctx context.Context) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(widget.APIUrl)
	if err != nil {
		widget.withError(fmt.Errorf("failed to fetch jobs: %w", err))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		widget.withError(fmt.Errorf("failed to read response: %w", err))
		return
	}
	var response RenderJobResponse
	if err := json.Unmarshal(body, &response); err != nil {
		widget.withError(fmt.Errorf("failed to parse jobs: %w", err))
		return
	}
	widget.Jobs = response.Jobs
	widget.withError(nil)
}

func (widget *renderJobWidget) Render() template.HTML {
	return widget.renderTemplate(widget, renderJobWidgetTemplate)
}
