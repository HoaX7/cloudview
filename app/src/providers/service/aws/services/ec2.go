package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"

	"cloudview/app/src/types"
	"context"
	"sync"

	service_types "cloudview/app/src/types"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Ec2Response struct {
	Reservations []ec2_types.Reservation // recervations give the number of running instances
	Volumes      []ec2_types.Volume
}

func (s Services) GetEC2Instances(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	instance := ec2.NewFromConfig(s.Config)
	result, err := instance.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		logger.Logger.Error("aws.GetEC2Data: Error describing instances:", err)
		respch <- service_types.GetDataOutput{
			Name:   constants.EC2,
			Result: nil,
		}
		return
	}

	/*
		Fetch Image information such as volumne size etc
		for each instance
	*/
	volumeIds := []string{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			volumeIds = append(volumeIds, *instance.BlockDeviceMappings[0].Ebs.VolumeId)
		}
	}
	res, err := instance.DescribeVolumes(context.TODO(), &ec2.DescribeVolumesInput{
		VolumeIds: volumeIds,
	})

	resp := &Ec2Response{
		Reservations: result.Reservations,
	}
	if err != nil {
		logger.Logger.Error("ec2.GetEc2Instances: error describing images for iamgeIds:", volumeIds, err)
	} else {
		resp.Volumes = res.Volumes
	}

	logger.Logger.Log("ec2.GetEc2Instances: success")
	if len(result.Reservations) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.EC2,
			Result: nil,
		}
		return
	}
	respch <- service_types.GetDataOutput{
		Name:   constants.EC2,
		Result: resp,
	}
	return
}

func (s Services) GetEC2Bandwidth(respch chan<- *types.Usage, wg *sync.WaitGroup, instanceId string) {
	defer wg.Done()
	var (
		namespace     = "AWS/EC2"
		dimensionName = "InstanceId"
	)

	result, err := s.getBandwidth(namespace, dimensionName, instanceId)
	if err != nil {
		logger.Logger.Error("ec2.GetBandwidth: ERROR", err)
		respch <- &types.Usage{}
		return
	}
	logger.Logger.Log("ec2.GetBandwidth: success")
	respch <- &types.Usage{
		Bandwidth: result,
	}
	return
}

func (s Services) GetEC2CpuUsage(respch chan<- *types.Usage, wg *sync.WaitGroup, instanceId string) {
	defer wg.Done()
	var (
		namespace     = "AWS/EC2"
		dimensionName = "InstanceId"
	)
	result, err := s.getCpuUsage(namespace, dimensionName, instanceId)
	if err != nil {
		logger.Logger.Error("ec2.GetEC2CpuUsage: ERROR", err)
		respch <- &types.Usage{}
		return
	}
	logger.Logger.Log("ec2.GetEC2CpuUsage: success")
	respch <- &types.Usage{
		Cpu: result,
	}
	return
}
