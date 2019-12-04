/**
 * @Author: jie.an
 * @Description:
 * @File:  cloudwath
 * @Version: 1.0.0
 * @Date: 2019/11/29 17:40
 */
package http

import (
	//"encoding/json"
	//"falcon-to-cloudwatch/g"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/open-falcon/falcon-plus/common/model"
	"strings"
)

//https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/cloudwatch/custom_metrics.go

//func main() {
//	mvs := []*model.MetricValue{&model.MetricValue{
//		Endpoint:  "a",
//		Metric:    "cpu.busy",
//		Value:     3.1415926,
//		Step:      0,
//		Type:      "c",
//		Tags:      "iface=eth0",
//		Timestamp: 0,
//	}, &model.MetricValue{
//		Endpoint:  "a",
//		Metric:    "load.5min",
//		Value:     3.14159555555526,
//		Step:      0,
//		Type:      "c",
//		Tags:      "mount=/,fstype=xfs",
//		Timestamp: 0,
//	}, &model.MetricValue{
//		Endpoint:  "a",
//		Metric:    "load.2min",
//		Value:     33098888888888888,
//		Step:      0,
//		Type:      "c",
//		Tags:      "mount=/,fstype=xfsd,datatype=1",
//		Timestamp: 0,
//	}, &model.MetricValue{
//		Endpoint:  "a",
//		Metric:    "load.20min",
//		Value:     10.1,
//		Step:      0,
//		Type:      "c",
//		Tags:      "",
//		Timestamp: 0,
//	}, &model.MetricValue{
//		Endpoint:  "b",
//		Metric:    "load.bdmin",
//		Value:     5.35487,
//		Step:      0,
//		Type:      "aa",
//		Tags:      "test=test1",
//		Timestamp: 0,
//	}, &model.MetricValue{
//		Endpoint:  "b",
//		Metric:    "load.cd",
//		Value:     1,
//		Step:      0,
//		Type:      "aa",
//		Tags:      "test=test2",
//		Timestamp: 0,
//	}}
//	//for j := 0; j < len(mvs); j++ {
//	//	//			mvs[j].Step = sec
//	//	mvs[j].Endpoint = "test-hostname.com"
//	//}
//	cfg := "cfg.test.json"
//	g.ParseConfig(cfg)
//	fmt.Println(g.Config().MonitoredPorts)
//	fmt.Println(g.Config().MonitoredProcs)
//	//m := make(map[string]map[int]string)
//	//c := make(map[int]string)
//	//c[2] = "ssh -ef"
//	//m["cmd"] = c
//	//data, err := json.Marshal(m)
//	//if err != nil {
//	//	fmt.Println("json marshal failed,err:", err)
//	//	return
//	//}
//	//fmt.Printf("%s\n", string(data))
//	fmt.Print(mvs)
//	//PushToCloudwatch(mvs)
//}

func getFloat64(unk interface{}) (value float64) {
	switch i := unk.(type) {
	case float64:
		return i
	case float32:
		return float64(i)
	case int:
		return float64(i)
	case int8:
		return float64(i)
	case int16:
		return float64(i)
	case int32:
		return float64(i)
	case int64:
		return float64(i)
	// Todo : Unhandle error
	default:
		return 0.000000000000
	}
}

// Push Metric To CloudWatch
func PushToCloudwatch(mvs []*model.MetricValue) {
	if len(mvs) == 0 {
		return
	}
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)
	// mvs process
	var metricdatums []*cloudwatch.MetricDatum
	for i, v := range mvs {
		var dimensions []*cloudwatch.Dimension
		//append instance id  to cloudwatch dimensions
		intancedimension := cloudwatch.Dimension{
			Name:  aws.String("InstanceId"),
			Value: aws.String(v.Endpoint),
		}
		dimensions = append(dimensions, &intancedimension)
		// append tags to cloudwatch metric
		if len(v.Tags) > 0 {
			v.Metric = v.Metric + "/" + v.Tags
		}
		// handle cloudwatch.MetricDatum type
		value := getFloat64(v.Value)
		metricdatum := cloudwatch.MetricDatum{
			MetricName: aws.String(v.Metric),
			Unit:       aws.String("Count"),
			Value:      aws.Float64(value),
			Dimensions: dimensions,
		}
		metricdatums = append(metricdatums, &metricdatum)
		if (i%9 == 0 && i != 0) || i == len(mvs)-1 {
			fmt.Println(metricdatums)
			_, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
				Namespace:  aws.String("Porsche-CloudWatch-Test"),
				MetricData: metricdatums})
			if err != nil {
				fmt.Println("Error adding metrics:", err.Error())
				return
			}
		}
	}
}

// Covert Tags To AWS Dimensions And Push Metric To Cloudwatch
func PushToCloudwatchWithTagsSplit(mvs []*model.MetricValue) {
	if len(mvs) == 0 {
		return
	}
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)
	// mvs process
	fmt.Println(mvs)
	var metricdatums []*cloudwatch.MetricDatum
	for i, v := range mvs {
		var dimensions []*cloudwatch.Dimension
		//append instance id  to cloudwatch dimensions
		intancedimension := cloudwatch.Dimension{
			Name:  aws.String("InstanceId"),
			Value: aws.String(v.Endpoint),
		}
		dimensions = append(dimensions, &intancedimension)
		// append tags to cloudwatch dimensions
		if len(v.Tags) > 0 {
			s := strings.Split(v.Tags, ",")
			for _, tag := range s {
				x := strings.Split(tag, "=")
				tagdimension := cloudwatch.Dimension{
					Name:  aws.String(x[0]),
					Value: aws.String(x[1]),
				}
				dimensions = append(dimensions, &tagdimension)
			}
		}
		// handle cloudwatch.MetricDatum type
		value := getFloat64(v.Value)
		metricdatum := cloudwatch.MetricDatum{
			MetricName: aws.String(v.Metric),
			Unit:       aws.String("Count"),
			Value:      aws.Float64(value),
			Dimensions: dimensions,
		}
		metricdatums = append(metricdatums, &metricdatum)
		if (i%3 == 0 && i != 0) || i == len(mvs)-1 {
			fmt.Println(metricdatums)
			_, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
				Namespace:  aws.String("Porsche-CloudWatch-Test"),
				MetricData: metricdatums})
			if err != nil {
				fmt.Println("Error adding metrics:", err.Error())
				return
			}
		}
	}
}
