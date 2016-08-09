package main

import (
	"flag"
	"fmt"
	"log"
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
		dimensionName      = flag.String("dimension.name", "", "metric dimension name, must not be empty")
		dimensionValue     = flag.String("dimension.value", "", "metric dimension value, must not be empty")
		monitoringDetailed = flag.Bool("detailed", false, "monitoring resoution: 1m (else: 5m)")
		statistics         = flag.String("stats", "Average", "possible values: Minimum, Maximum, Average, Sum, SampleCount")
		awsRegion          = flag.String("aws.region", "eu-central-1", "AWS region")
		period             = 5 * time.Minute
	)
	flag.Parse()

	if *monitoringDetailed {
		period = time.Minute
	}
	startTime := time.Now().Add(-1 * period)

	if !statsMap[*statistics] || *namespace == "" || *dimensionName == "" || *dimensionValue == "" {
		flag.Usage()
		return
	}

	cloudwatchCli := cloudwatch.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))

	params := cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(time.Now()),
		MetricName: metricName,
		Namespace:  namespace,
		Period:     aws.Int64(int64(period.Seconds())),
		StartTime:  &startTime,
		Statistics: []*string{statistics},
		Dimensions: []*cloudwatch.Dimension{{
			Name:  dimensionName,
			Value: dimensionValue,
		}},
	}
	resp, err := cloudwatchCli.GetMetricStatistics(&params)
	if err != nil {
		log.Fatalf("Could not get metric statistics for %#v: %v", params, err)
	}

	if len(resp.Datapoints) == 0 {
		// no data
		fmt.Println("0")
	} else {
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
	}
}
