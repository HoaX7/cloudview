package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cw_types "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Metrics struct {
	Name       string
	Statistics *cloudwatch.GetMetricStatisticsOutput
}
type S3Result struct {
	Data    *s3.ListBucketsOutput
	Metrics []S3Metrics
	ACLList []ACLList
}
type ACLList struct {
	Bucket string
	ACL    *s3.GetBucketAclOutput
}

func (s Services) GetS3Buckets(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()

	instance := s3.NewFromConfig(s.Config)
	result, err := instance.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		logger.Logger.Error("s3.GetS3Buckets: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.S3,
			Result: nil,
		}
		return
	}

	s.CloudWatchInit()
	ch := make(chan CWChannel, len(result.Buckets))
	swg := &sync.WaitGroup{}

	bucketNameLabel := "BucketName"
	namespace := "AWS/S3"

	// Replace 'your-metric-name' with the desired metric name (e.g., 'BucketSizeBytes')
	metricName := "BucketSizeBytes"

	// Replace 'your-start-time' and 'your-end-time' with the desired time range
	// cloudwatch updates data every 48 hours.
	startTime := time.Now().Add(-48 * time.Hour) // Two days ago
	endTime := time.Now()

	// Replace 'your-period' with the desired period in seconds (e.g., 3600 for 1 hour)
	var period int32 = 3600

	storageType := "StandardStorage"
	storageName := "StorageType"

	aclList := []ACLList{}

	for _, bucket := range result.Buckets {
		logger.Logger.Log("s3.GetBuckets: fetching metrics & ACL for bucket:", *bucket.Name)

		/*
			Fetch bucket ACL to see read/write permissions
		*/
		swg.Add(1)
		go func(bucketName string) {
			defer swg.Done()
			res, err := instance.GetBucketAcl(context.TODO(), &s3.GetBucketAclInput{
				Bucket: &bucketName,
			})
			if err != nil {
				return
			}
			aclList = append(aclList, ACLList{
				Bucket: bucketName,
				ACL:    res,
			})
			return
		}(*bucket.Name)

		swg.Add(1)
		/**
		This method allows us to fetch S3 bucket size
		*/
		go s.GetCloudWatchMetricStatistics(ch, swg, &cloudwatch.GetMetricStatisticsInput{
			Namespace:  &namespace,
			MetricName: &metricName,
			Statistics: []cw_types.Statistic{"Sum"},
			Dimensions: []cw_types.Dimension{{
				Name:  &bucketNameLabel,
				Value: bucket.Name,
			}, {
				Name:  &storageName,
				Value: &storageType,
			}},
			Unit:      "Bytes",
			StartTime: &startTime,
			EndTime:   &endTime,
			Period:    &period,
		}, *bucket.Name)
	}

	go func() {
		swg.Wait()
		close(ch)
	}()

	s3MetricData := []S3Metrics{}
	for metricData := range ch {
		// Appending to associate each bucket with its own metric data
		s3MetricData = append(s3MetricData, S3Metrics{
			Name:       metricData.BucketName,
			Statistics: metricData.Metrics,
		})
	}

	logger.Logger.Log("s3.GetS3Buckets: success")
	if len(result.Buckets) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.S3,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name: constants.S3,
		Result: S3Result{
			Data:    result,
			Metrics: s3MetricData,
			ACLList: aclList,
		},
	}
	return
}
