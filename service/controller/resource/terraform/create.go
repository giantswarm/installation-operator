package terraform

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/installation-operator/service/controller/controllercontext"
	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	createTablesState, err := toTablesState(createChange)
	if err != nil {
		return microerror.Mask(err)
	}
	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, tableInput := range createTablesState {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating dynamodb table %#q", tableInput.Name))

		{
			i := &dynamodb.CreateTableInput{
				TableName: aws.String(fmt.Sprintf("%s-lock", cr.Name)),
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					{
						AttributeName: aws.String("LockID"),
						AttributeType: aws.String("S"),
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("LockID"),
						KeyType:       aws.String("HASH"),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			}

			_, err = cc.Client.AWS.DynamoDB.CreateTable(i)
			if IsTableAlreadyExists(err) {
				// Fall through.
			} else if err != nil {
				return microerror.Mask(err)
			}
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created dynamodb table %#q", tableInput.Name))
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentBuckets, err := toTablesState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredBuckets, err := toTablesState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var createState []TableState
	for _, bucket := range desiredBuckets {
		if !containsTableState(bucket.Name, currentBuckets) {
			createState = append(createState, bucket)
		}
	}

	return createState, nil
}
