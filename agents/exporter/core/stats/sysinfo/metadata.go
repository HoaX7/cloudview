package sysinfo

type sys interface {
	instanceId() string
}

type unknownPlatform struct{}

func currentProvider() (sys, error) {
	if isAWSInstance() {
		return awsProvider()
	}
	return unknownProvider()
}

func (u *unknownPlatform) instanceId() string {
	return ""
}

func unknownProvider() (sys, error) {
	return &unknownPlatform{}, nil
}
