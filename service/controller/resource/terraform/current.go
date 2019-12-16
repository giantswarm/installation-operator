package terraform

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/giantswarm/microerror"
	"golang.org/x/sync/errgroup"

	"github.com/giantswarm/installation-operator/service/controller/controllercontext"
	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	bucketStateNames := []string{
		key.TargetLogBucketName(cr),
		key.BucketName(&cr, cc.Status.AWSAccountID),
	}

	var currentBucketState []TableState
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "finding the S3 buckets")

		g := &errgroup.Group{}
		m := sync.Mutex{}

		for _, inputBucketName := range bucketStateNames {
			bucketName := inputBucketName

			g.Go(func() error {
				inputBucket := TableState{
					Name: bucketName,
				}

				// TODO this check should not be done here. Here we only fetch the
				// current state. We have to make a request anyway so fetching what we
				// want and handling the not found errors as usual should be the way to
				// go.
				//
				//
				//     https://github.com/giantswarm/giantswarm/issues/5246
				//
				isCreated, err := r.isTableCreated(ctx, bucketName)
				if err != nil {
					return microerror.Mask(err)
				}
				if !isCreated {
					return nil
				}

				r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finding the S3 bucket %#q", bucketName))

				m.Lock()
				currentBucketState = append(currentBucketState, inputBucket)
				m.Unlock()

				r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found the S3 bucket %#q", bucketName))

				return nil
			})
		}

		err := g.Wait()
		if err != nil {
			return nil, microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "found the S3 buckets")
	}

	return currentBucketState, nil
}

func (r *Resource) isTableCreated(ctx context.Context, name string) (bool, error) {
	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return false, microerror.Mask(err)
	}

	headInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	}
	_, err = cc.Client.AWS.DynamoDB.DescribeTable(headInput)
	if IsTableNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

