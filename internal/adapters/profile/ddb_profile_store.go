package adapter_profile

import (
	"context"
	"errors"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang-template/internal/domain/action"
	"golang-template/internal/domain/profile"
	"golang-template/internal/domain/transaction"
	"log/slog"
	"net/http"
)

const (
	ProfileElement string = "profile"
)

type (
	DynamoDbProfileStore struct {
		client     *dynamodb.Client
		table      *string
		ownerIndex *string
	}
)

func NewDDBProfileStore(cfg aws.Config) *DynamoDbProfileStore {
	return &DynamoDbProfileStore{
		client:     dynamodb.NewFromConfig(cfg),
		table:      aws.String("wilson-entities-stage"),
		ownerIndex: aws.String("secondary-lookup"),
	}
}

func (d DynamoDbProfileStore) GetProfilesByOwner(ctx context.Context, id profile.OwnerId) ([]profile.Profile, error) {
	keyEx := expression.Key("secondary-id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, xhttp.NewError("cannot construct key", http.StatusInternalServerError)
	}

	output, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 d.table,
		IndexName:                 d.ownerIndex,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, xhttp.NewError("unable to find profiles for owner", http.StatusInternalServerError)
	}

	profiles := make([]profile.Profile, 0)
	for _, profLookup := range output.Items {
		var ddbProfile ddbProfileAdapter
		err = attributevalue.UnmarshalMap(profLookup, &ddbProfile)

		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, xhttp.NewError("unable to convert profile", http.StatusInternalServerError)
		}

		profOutput, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: d.table,
			Key:       ddbProfile.getKey(),
		})

		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, xhttp.NewError("unable to convert profile", http.StatusInternalServerError)
		}

		err = attributevalue.UnmarshalMap(profOutput.Item, &ddbProfile)

		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, xhttp.NewError("unable to convert profile", http.StatusInternalServerError)
		}

		profiles = append(profiles, profile.Profile{
			Id:       profile.Id(ddbProfile.Id),
			Name:     ddbProfile.Name,
			Owner:    profile.OwnerId(ddbProfile.Owner),
			Level:    ddbProfile.Level,
			Gold:     ddbProfile.Gold,
			Location: profile.LocationId(ddbProfile.Location),
		})
	}

	return profiles, nil
}

func (d DynamoDbProfileStore) Save(ctx context.Context, p *profile.Profile) error {
	ddbProfile := profileToDDB(p)

	item, err := attributevalue.MarshalMap(ddbProfile)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return xhttp.NewError("unable to convert profile", http.StatusInternalServerError)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: d.table,
		Item:      item,
	})

	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return xhttp.NewError("unable to create profile", http.StatusInternalServerError)
	}

	return nil
}

func (d DynamoDbProfileStore) Get(ctx context.Context, id action.ProfileId) (action.Profile, error) {
	ddbProf := &ddbProfileAdapter{
		Id: string(id),
	}

	output, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: d.table,
		Key:       ddbProf.getKey(),
	})

	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return action.Profile{}, xhttp.NewError("unable to query profile", http.StatusInternalServerError)
	}

	if len(output.Item) == 0 {
		return action.Profile{}, xhttp.NewError("profile not found", http.StatusNotFound)
	}

	err = attributevalue.UnmarshalMap(output.Item, &ddbProf)

	if err != nil {
		return action.Profile{}, xhttp.NewError("unable to convert profile", http.StatusInternalServerError)
	}

	return action.Profile{
		Id:       action.ProfileId(ddbProf.Id),
		Location: action.LocationId(ddbProf.Location),
	}, nil
}

func (d DynamoDbProfileStore) GetAll(ctx context.Context) ([]action.Profile, error) {
	output, err := d.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:      d.table,
		ConsistentRead: aws.Bool(false),
	})

	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, xhttp.NewError("unable to query profiles", http.StatusInternalServerError)
	}

	profiles := make([]action.Profile, 0)
	for _, dbProf := range output.Items {
		var convertedProfile action.Profile
		err = attributevalue.UnmarshalMap(dbProf, &convertedProfile)

		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, xhttp.NewError("unable to convert profile", http.StatusInternalServerError)
		}

		profiles = append(profiles, convertedProfile)
	}

	return profiles, nil
}

func (d DynamoDbProfileStore) GetBalance(ctx context.Context, profile string) (transaction.Balance, error) {
	ddbProf := ddbProfileAdapter{
		Id: profile,
	}

	output, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: d.table,
		Key:       ddbProf.getKey(),
	})

	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return transaction.Balance{}, xhttp.NewError("unable to query profile", http.StatusInternalServerError)
	}

	if len(output.Item) == 0 {
		return transaction.Balance{}, xhttp.NewError("profile not found", http.StatusNotFound)
	}

	err = attributevalue.UnmarshalMap(output.Item, &ddbProf)

	if err != nil {
		return transaction.Balance{}, xhttp.NewError("unable to convert profile", http.StatusInternalServerError)
	}

	return transaction.Balance{
		Profile: profile,
		Gold:    ddbProf.Gold,
		Version: ddbProf.Version,
	}, nil
}

func (d DynamoDbProfileStore) UpdateBalance(ctx context.Context, balance transaction.Balance) (transaction.Balance, error) {
	ddbProf := ddbProfileAdapter{
		Id: balance.Profile,
	}

	update := expression.Set(expression.Name("gold"), expression.Value(balance.Gold))
	update.Set(expression.Name("version"), expression.Value(balance.Version+1))

	condition := expression.Equal(expression.Name("version"), expression.Value(balance.Version))

	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return transaction.Balance{}, xhttp.NewError("unable to build conditions for balance update", http.StatusInternalServerError)
	}

	output, err := d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 d.table,
		Key:                       ddbProf.getKey(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})

	if err != nil {
		if errors.Is(err, &types.ConditionalCheckFailedException{}) {
			return transaction.Balance{}, transaction.VersionMismatchError{}
		}

		slog.ErrorContext(ctx, err.Error())
		return transaction.Balance{}, xhttp.NewError("error while updating balance", http.StatusInternalServerError)
	}

	err = attributevalue.UnmarshalMap(output.Attributes, &ddbProf)

	if err != nil {
		return transaction.Balance{}, xhttp.NewError("unable to convert profile", http.StatusInternalServerError)
	}

	return transaction.Balance{
		Profile: ddbProf.Id,
		Gold:    ddbProf.Gold,
		Version: ddbProf.Version,
	}, nil
}
