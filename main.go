package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func main() {
	namespace := flag.String("namespace", "", "the namespace of the metric")
	metricName := flag.String("metric-name", "", "the name of the desired metric")
	dimensionName := flag.String("dimension-name", "", "the name of the metric's dimension")
	dimensionValue := flag.String("dimension-value", "", "the value of the metric's dimension")
	monitoringType := flag.String("monitoring-type", "", "monitoring type, choose between basic and detailed")
	statistics := flag.String("statistics", "", "Minimum, Maximum, Average, Sum, SampleCount")
	awsRegion := flag.String("aws-region", "eu-central-1", "AWS region")
	awsAccessKey := flag.String("aws-access-key", "", "AWS access key")
	awsSecretKey := flag.String("aws-secret-key", "", "AWS secret key")

	flag.Parse()

	fmt.Println("namespace:", *namespace)
	fmt.Println("metricName:", *metricName)
	fmt.Println("dimensionName:", *dimensionName)
	fmt.Println("dimensionValue:", *dimensionValue)
	fmt.Println("monitoringType:", *monitoringType)
	fmt.Println("statistics:", *statistics)
	fmt.Println("awsRegion:", *awsRegion)
	fmt.Println("awsAccessKey", *awsAccessKey)
	fmt.Println("awsSecretKey", *awsSecretKey)

	var period int64
	var startTime time.Time

	if *monitoringType == "detailed" {
		startTime = oneMinuteAgo()
		period = 60
	} else {
		startTime = fiveMinutesAgo()
		period = 360
	}

	svc := cloudwatch.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))

	params := &cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(time.Now()),
		MetricName: aws.String(*metricName),
		Namespace:  aws.String(*namespace),
		Period:     aws.Int64(period),
		StartTime:  aws.Time(startTime),
		Statistics: []*string{
			aws.String(*statistics),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String(*dimensionName),
				Value: aws.String(*dimensionValue),
			},
		},
		//Unit: aws.String("Seconds"),
	}

	resp, err := svc.GetMetricStatistics(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}

func oneMinuteAgo() time.Time {
	return time.Now().Add(-1 * time.Minute)
}

func fiveMinutesAgo() time.Time {
	return time.Now().Add(-5 * time.Minute)
}
