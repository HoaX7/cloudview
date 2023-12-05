package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

// cloudfront a CDN offered by aws. Use the data from this to check if
// s3 bucket is configured with a cdn
// check if the domain name returned in cloudfront contains s3 bucket names
func (s Services) GetCloudfrontDistributions(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	instance := cloudfront.NewFromConfig(s.Config)
	result, err := instance.ListDistributions(context.TODO(), &cloudfront.ListDistributionsInput{})
	if err != nil {
		logger.Logger.Error("GetCloudFrontDistributions: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.CLOUDFRONT,
			Result: nil,
		}
		return
	}
	if len(result.DistributionList.Items) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.CLOUDFRONT,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.CLOUDFRONT,
		Result: result.DistributionList,
	}
	return
}
