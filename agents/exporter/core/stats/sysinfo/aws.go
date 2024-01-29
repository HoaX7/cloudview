package sysinfo

import "cloudview/agents/exporter/core/iox"

type awsPlatform struct{}

const instanceIdFile = "/var/lib/cloud/data/instance-id"

func isAWSInstance() bool {
	return iox.Exists(instanceIdFile)
}

func (aws *awsPlatform) instanceId() string {
	id, _ := iox.ReadText(instanceIdFile)
	return id
}

func awsProvider() (sys, error) {
	return &awsPlatform{}, nil
}
