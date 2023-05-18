package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"google.golang.org/api/monitoring/v3"
	"google.golang.org/api/option"
)

func Chart(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("Failed to create service. %v", err)
	}
	projectID := "junior-engineers-gym-2023"
	filter := `metric.type="compute.googleapis.com/instance/cpu/utilization"`

	location, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(location)
	startTime := now.Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := now.Format(time.RFC3339)
	res, err := svc.Projects.TimeSeries.List("projects/" + projectID).Filter(filter).
		IntervalStartTime(startTime).
		IntervalEndTime(endTime).
		AggregationPerSeriesAligner("ALIGN_MAX").
		AggregationAlignmentPeriod("60s").
		AggregationCrossSeriesReducer("REDUCE_MAX").
		Do()
	if err != nil {
		fmt.Printf("Could not execute request: %v", err)
	}

	var timeList []string
	var valueList []opts.LineData
	for _, ts := range res.TimeSeries {
		for _, point := range ts.Points {
			if point.Value.DoubleValue != nil {
				value := *point.Value.DoubleValue * 100
				timeStamp := point.Interval.EndTime
				t, _ := time.Parse(time.RFC3339, timeStamp)
				timeList = append(timeList, t.In(location).Format(time.RFC3339))
				valueList = append(valueList, opts.LineData{Value: value})
			}
		}
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "CPU Utilization",
		}),
	)
	line.SetXAxis(timeList).
		AddSeries("Series1", valueList)

	page := components.NewPage()
	page.AddCharts(line)

	f, err := os.Create("output.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	page.Render(f)
}

func main() {
	http.HandleFunc("/", Chart)
	http.ListenAndServe(":8080", nil)
}
