package glance

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var afanasyMachinesTemplate = mustParseTemplate("afanasy-machines.html", "widget-base.html")

type afanasyMachinesWidget struct {
	widgetBase    `yaml:",inline"`
	AllowInsecure bool             `yaml:"allow-insecure"`
	Machines      []AfanasyMachine `yaml:"-"`
}

type AfanasyMachine struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	CPUUsage   int    `json:"cpu_usage"`
	MemTotalMB int    `json:"mem_total_mb"`
	MemUsedMB  int    `json:"mem_used_mb"`
}

type AfanasyMachinesData struct {
	Machines []AfanasyMachine `json:"renders"`
}

func (widget *afanasyMachinesWidget) initialize() error {
	widget.withTitle("Afanasy Render Machines").withCacheDuration(30 * time.Second)
	return nil
}

func (widget *afanasyMachinesWidget) update(ctx context.Context) {
	url := "http://192.168.90.104:5000/afanasy/renders"
	var client *http.Client
	if widget.AllowInsecure {
		client = defaultInsecureHTTPClient
	} else {
		client = defaultHTTPClient
	}
	resp, err := client.Get(url)
	if err != nil {
		widget.withError(fmt.Errorf("failed to fetch machines from API: %w", err))
		return
	}
	defer resp.Body.Close()

	// Support JSON with a 'renders' key
	type rendersResponse struct {
		Renders []AfanasyMachine `json:"renders"`
	}
	var rendersData rendersResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&rendersData); err != nil {
		widget.withError(fmt.Errorf("failed to decode machine data JSON: %w", err))
		return
	}
	widget.Machines = rendersData.Renders
}

func (widget *afanasyMachinesWidget) Render() template.HTML {
	return widget.renderTemplate(widget, afanasyMachinesTemplate)
}
