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

type afanasyMachinesWidget struct {
	widgetBase
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
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		widget.withError(fmt.Errorf("failed to fetch machines from API: %w", err))
		return
	}
	defer resp.Body.Close()

	var data []AfanasyMachinesData
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		widget.withError(fmt.Errorf("failed to decode machine data JSON: %w", err))
		return
	}
	var machines []AfanasyMachine
	if len(data) > 0 {
		machines = data[0].Machines
	}

	tmpl, err := template.ParseFiles("internal/glance/templates/afanasy-machines.html")
	if err != nil {
		widget.withError(fmt.Errorf("failed to parse template: %w", err))
		return
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, machines); err != nil {
		widget.withError(fmt.Errorf("failed to execute template: %w", err))
		return
	}
	widget.templateBuffer.Reset()
	widget.templateBuffer.Write(buf.Bytes())
}

func (widget *afanasyMachinesWidget) Render() template.HTML {
	return template.HTML(widget.templateBuffer.String())
}
