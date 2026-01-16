package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type afanasyJobsWidget struct {
	widgetBase
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
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		widget.withError(fmt.Errorf("failed to fetch jobs from n8n: %w", err))
		return
	}
	defer resp.Body.Close()

	var jobs []AfanasyJob
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&jobs); err != nil {
		widget.withError(fmt.Errorf("failed to decode jobs JSON: %w", err))
		return
	}

	tmpl, err := template.ParseFiles("internal/glance/templates/afanasy-jobs.html")
	if err != nil {
		widget.withError(fmt.Errorf("failed to parse template: %w", err))
		return
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, jobs); err != nil {
		widget.withError(fmt.Errorf("failed to execute template: %w", err))
		return
	}
	widget.templateBuffer.Reset()
	widget.templateBuffer.Write(buf.Bytes())
}

func (widget *afanasyJobsWidget) Render() template.HTML {
	return template.HTML(widget.templateBuffer.String())
}
