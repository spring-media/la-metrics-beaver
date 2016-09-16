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
		namespace      = flag.String("namespace", "", "metric namespace, must not be empty")
		metricName     = flag.String("metric.name", "", "metric name, must not be empty")
		dimensionName  = flag.String("dimension.name", "", "metric dimension name")
		dimensionValue = flag.String("dimension.value", "", "metric dimension value")
		period         = flag.Duration("period", time.Minute, "the time span in minutes like 1m, 2m, 1440m.")
		statistics     = flag.String("stats", "Average", "possible values: Minimum, Maximum, Average, Sum, SampleCount")
		awsRegion      = flag.String("aws.region", "eu-central-1", "AWS region")
	)
	flag.Parse()

	// Normally we are using a large time span and a smaller period.
	// Here we are seeking for only one datapoint.
	startTime := time.Now().Add(-1 * (*period))

	if !statsMap[*statistics] || *namespace == "" {
		flag.Usage()
		return
	}

	cloudwatchCli := cloudwatch.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))

	dimensions := []*cloudwatch.Dimension{}

	dimensionsNames := strings.Split(*dimensionName, ",")
	dimensionsValues := strings.Split(*dimensionValue, ",")

	empty := (len(dimensionsNames) < 1 || len(dimensionsValues) < 1)
	sameCount := len(dimensionsNames) == len(dimensionsValues)
	if !empty && sameCount {
		for i := 0; i < len(dimensionsNames); i++ {
			dimN := dimensionsNames[i]
			dimV := dimensionsValues[i]

			dim := cloudwatch.Dimension{
				Name:  &dimN,
				Value: &dimV,
			}

			dimensions = append(dimensions, &dim)
		}
	}

	if *namespace == "AWS/CloudFront" {
		// if not global region, cloudwatch always return 0 for cloudfront
		dimensions = append(dimensions, &cloudwatch.Dimension{
			Name:  aws.String("Region"),
			Value: aws.String("Global")})
	}

	if *namespace == "AWS/Billing" {
		dimensions = append(dimensions, &cloudwatch.Dimension{
			Name:  aws.String("Currency"),
			Value: aws.String("USD")})
	}

	params := cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(time.Now()),
		MetricName: metricName,
		Namespace:  namespace,
		Period:     aws.Int64(int64(period.Seconds())),
		StartTime:  &startTime,
		Statistics: []*string{statistics},
		Dimensions: dimensions,
	}

	resp, err := cloudwatchCli.GetMetricStatistics(&params)

	if err != nil {
		log.Fatalf("Could not get metric statistics for %#v: %v", params, err)
	}

	if len(resp.Datapoints) == 0 {
		// no data
		fmt.Println("0")

	} else if len(resp.Datapoints) == 1 {
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
		log.Fatalf("There should be no more than one Datapoint: %v", resp)
	}
}
