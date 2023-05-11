package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
	"google.golang.org/api/monitoring/v3"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	keyJson := os.Getenv("KEY_JSON")
	if keyJson == "" {
		fmt.Println("There's no KEY_JSON.")
		return
	} else {
		fmt.Println("KEY_JSON: OK")
	}
	svc, err := monitoring.NewService(ctx, option.WithCredentialsJSON([]byte(keyJson)))
	if err != nil {
		fmt.Println("Failed to create Service.")
	}
	projectID := "junior-engineers-gym-2023"
	filter := `metric.type="compute.googleapis.com/instance/cpu/utilization"`
	startTime := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := time.Now().UTC().Format(time.RFC3339)
	res, err := svc.Projects.TimeSeries.List("projects/" + projectID).Filter(filter).IntervalStartTime(startTime).IntervalEndTime(endTime).Do()
	if err != nil {
		fmt.Printf("Could not execute request: %v", err)
	}

	var values []chart.Value
	for _, ts := range res.TimeSeries {
		for _, point := range ts.Points {
			value := *point.Value.DoubleValue
			timeStamp := point.Interval.EndTime
			t, _ := time.Parse(time.RFC3339, timeStamp)
			values = append(values, chart.Value{Value: value, Label: t.Format(time.RFC3339)})
		}
	}

	// Create a new chart.
	graph := chart.BarChart{
		Title: "CPU Utilization",
		Bars:  values,
	}

	// Save the chart to a file.
	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
