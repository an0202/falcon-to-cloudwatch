/**
 * @Author: jie.an
 * @Description:
 * @File:  cloudwath
 * @Version: 1.0.0
 * @Date: 2019/11/29 17:40
 */
package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/open-falcon/falcon-plus/common/model"
//	"time"
)


//https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/cloudwatch/custom_metrics.go


func main(){
	mvs := []*model.MetricValue{&model.MetricValue{
		Endpoint:  "a",
		Metric:    "b",
		Value:     nil,
		Step:      0,
		Type:      "c",
		Tags:      "iface=eth0",
		Timestamp: 0,
	},&model.MetricValue{
		Endpoint:  "a",
		Metric:    "b",
		Value:     nil,
		Step:      0,
		Type:      "c",
		Tags:      "mount=/,fstype=xfs",
		Timestamp: 0,
	}}
	for j := 0; j < len(mvs); j++ {
		//			mvs[j].Step = sec
		mvs[j].Endpoint = "test-hostname.com"
	}

	chanelprac()
}

func chanelprac(){
	var numList []int
	for i:=0;i<103;i++{
		numList = append(numList, i)
		if i%9 == 0 && i!=0 {
			fmt.Println(numList)
			numList = []int{}
		}
		if i == 102 {
			fmt.Println(numList)
		}
	}
}

func pushdata(){
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)

	_, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("Site/Traffic"),
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String("UniqueVisitors"),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(5885.0),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("SiteName"),
						Value: aws.String("example.com"),
					},
				},
			},
			&cloudwatch.MetricDatum{
				MetricName: aws.String("UniqueVisits"),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(8628.0),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("SiteName"),
						Value: aws.String("example.com"),
					},
				},
			},
			&cloudwatch.MetricDatum{
				MetricName: aws.String("PageViews"),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(18057.0),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("PageURL"),
						Value: aws.String("my-page.html"),
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Error adding metrics:", err.Error())
		return
	}

	// Get information about metrics
	result, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
		Namespace: aws.String("Site/Traffic"),
	})
	if err != nil {
		fmt.Println("Error getting metrics:", err.Error())
		return
	}

	for _, metric := range result.Metrics {
		fmt.Println(*metric.MetricName)

		for _, dim := range metric.Dimensions {
			fmt.Println(*dim.Name + ":", *dim.Value)
			fmt.Println()
		}
	}
}