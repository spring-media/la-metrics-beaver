package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

var statsMap = map[string]bool{"Minimum": true, "Maximum": true, "Average": true, "Sum": true, "SampleCount": true}

func main() {
	var (
		namespace          = flag.String("namespace", "", "metric namespace, must not be empty")
		metricName         = flag.String("metric.name", "", "metric name, must not be empty")
		monitoringDetailed = flag.Bool("detailed", false, "monitoring resoution: 1m (else: 5m)")
		statistics         = flag.String("stats", "Average", "possible values: Minimum, Maximum, Average, Sum, SampleCount")
		awsRegion          = flag.String("aws.region", "eu-central-1", "AWS region")
		dimensions         = flag.String("dimensions", "", "dimension name and dimension value")
		period             = 5 * time.Minute
	)
	flag.Parse()

	if *monitoringDetailed {
		period = time.Minute
	}
	startTime := time.Now().Add(-1 * period)

	if !statsMap[*statistics] || *namespace == "" || *dimensions == "" {
		flag.Usage()
		return
	}

	err, d := parseDimensions(*dimensions)

	if err != nil {
		fmt.Println(err)
	}

	cloudwatchCli := cloudwatch.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))

	params := cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(time.Now()),
		MetricName: metricName,
		Namespace:  namespace,
		Period:     aws.Int64(int64(period.Seconds())),
		StartTime:  &startTime,
		Statistics: []*string{statistics},
		Dimensions: d,
	}
	resp, err := cloudwatchCli.GetMetricStatistics(&params)
	if err != nil {
		log.Fatalf("Could not get metric statistics for %#v: %v", params, err)
	}

	if len(resp.Datapoints) != 0 {
		switch *statistics {
		case "Sum":
			fmt.Println(*resp.Datapoints[0].Sum)
		case "Minimum":
			fmt.Println(*resp.Datapoints[0].Minimum)
		case "Maximum":
			fmt.Println(*resp.Datapoints[0].Maximum)
		case "Average":
			fmt.Println(*resp.Datapoints[0].Average)
		case "SampleCount":
			fmt.Println(*resp.Datapoints[0].SampleCount)
		default:
			fmt.Println(*resp.Datapoints[0].Sum)
		}

	} else {
		fmt.Println("nodata")
	}
}

func parseDimensions(dimensions string) (error, []*cloudwatch.Dimension) {
	splitted := strings.Split(dimensions, ":")

	if len(splitted) != 2 {
		return fmt.Errorf("Parsing dimensions failed, there should be two values, got %v. Check you format. Should be dValue:dName", len(splitted)), nil
	}

	return nil, []*cloudwatch.Dimension{{
		Name:  &splitted[0],
		Value: &splitted[1],
	}}
}
