package services

import (
	"cloudview/app/src/api/middleware/logger"
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func (s *Services) CloudWatchInit() {
	s.CloudWatch = cloudwatch.NewFromConfig(s.Config)
}

//	func (s Services) CloudWatchMetrics() {
//		result, err := s.CloudWatch.GetMetricData()
//	}
type CWChannel struct {
	Metrics    *cloudwatch.GetMetricStatisticsOutput
	BucketName string
}

func (s Services) GetCloudWatchMetricStatistics(respch chan<- CWChannel,
	wg *sync.WaitGroup,
	input *cloudwatch.GetMetricStatisticsInput,
	bucketName string) {
	defer wg.Done()
	result, err := s.CloudWatch.GetMetricStatistics(context.TODO(), input)
	if err != nil {
		logger.Logger.Error("cloudwatch.GetCloudWatchMetricStatistics: ERROR", err)
		respch <- CWChannel{
			BucketName: bucketName,
		}
		return
	}
	respch <- CWChannel{
		BucketName: bucketName,
		Metrics:    result,
	}
	return
}

/*
*
This function can be used to retrieve usage data such as bandwidth/cpu/ram usage.
Input takes an array of metrics to be returned.

	IMPORTANT CONCEPT:
	The Metrics Data returned by this function is grouped by period specified.
	Ex: if `period` is 24 hours. The data will be aggregated to every 24 hour
	timestamp between start - end dates.

	Scenario: To fetch bandwidth data for a month, the `period` specified is
	the total seconds elapsed since the start of the month and MUST BE A MULTIPLE OF 60.
	This will return the total bandwidth for the month in the 0th position of the array.

	To fetch Sum,Average use the `Stat` prop. You can also use this prop to fetch
	Maximum,Minimum values for a period along with timestamp.
*/
func (s Services) getCloudWatchMetricData(input *cloudwatch.GetMetricDataInput) (*cloudwatch.GetMetricDataOutput, error) {
	return s.CloudWatch.GetMetricData(context.TODO(), input)
}

func (s Services) getBandwidth(namespace string, dimensionName string, value string) (*cloudwatch.GetMetricDataOutput, error) {
	var (
		stat                = "Sum"
		MaxDatapoints int32 = 1
		ReturnData          = true
		networkIn           = "NetworkIn"
		networkOut          = "NetworkOut"
	)
	now := time.Now()
	// fetch data from start of the month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Calculating period to fetch total bandwidth consumed for the month.
	var period int32 = calculatePeriod(startOfMonth, now)
	logger.Logger.Log("cloudwatch.GetBandwidth: fetching networkin/networkout metrics for instance:", dimensionName, value, "between:", startOfMonth, now)
	return s.getCloudWatchMetricData(&cloudwatch.GetMetricDataInput{
		MaxDatapoints: &MaxDatapoints,
		StartTime:     &startOfMonth,
		EndTime:       &now,
		MetricDataQueries: []types.MetricDataQuery{{
			ReturnData: &ReturnData,
			MetricStat: &types.MetricStat{
				Period: &period,
				Stat:   &stat,
				Metric: &types.Metric{
					Dimensions: []types.Dimension{{
						Name:  &dimensionName,
						Value: &value,
					}},
					Namespace:  &namespace,
					MetricName: &networkIn,
				},
			},
			Id: aws.String("networkIn"),
		}, {
			Id:         aws.String("networkOut"),
			ReturnData: &ReturnData,
			MetricStat: &types.MetricStat{
				Period: &period,
				Stat:   &stat,
				Metric: &types.Metric{
					Dimensions: []types.Dimension{{
						Name:  &dimensionName,
						Value: &value,
					}},
					Namespace:  &namespace,
					MetricName: &networkOut,
				},
			},
		}},
	})
}

func (s Services) getCpuUsage(namespace string, dimensionName string, value string) (*cloudwatch.GetMetricDataOutput, error) {
	var (
		stat                = "Sum"
		MaxDatapoints int32 = 6
		ReturnData          = true
		metricName          = "CPUUtilization"
		period        int32 = 3600
	)

	now := time.Now()
	startTime := now.Add(-24 * time.Hour)

	logger.Logger.Log("cloudwatch.GetBandwidth: fetching cpu-usage metrics for:", dimensionName, value, "between:", startTime, now)
	return s.getCloudWatchMetricData(&cloudwatch.GetMetricDataInput{
		MaxDatapoints: &MaxDatapoints,
		StartTime:     &startTime,
		EndTime:       &now,
		MetricDataQueries: []types.MetricDataQuery{{
			ReturnData: &ReturnData,
			MetricStat: &types.MetricStat{
				Period: &period,
				Stat:   &stat,
				Metric: &types.Metric{
					Dimensions: []types.Dimension{{
						Name:  &dimensionName,
						Value: &value,
					}},
					Namespace:  &namespace,
					MetricName: &metricName,
				},
			},
			Id: aws.String("cpuUsage"),
		}},
	})
}

func calculatePeriod(start, end time.Time) int32 {
	// Calculate the duration between start and end
	duration := end.Sub(start)

	// Convert duration to seconds
	periodSeconds := int32(duration.Seconds())

	// Round to the nearest multiple of 60
	roundedPeriod := (periodSeconds + 30) / 60 * 60

	return roundedPeriod
}
