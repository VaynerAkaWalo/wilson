package adapter_profile

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang-template/internal/domain/profile"
)

type ddbProfileAdapter struct {
	Id       string `dynamodbav:"id"`
	Element  string `dynamodbav:"element"`
	Owner    string `dynamodbav:"secondary-id"`
	Name     string `dynamodbav:"name"`
	Location string `dynamodbav:"location"`
	Level    int64  `dynamodbav:"level"`
	Gold     int64  `dynamodbav:"gold"`
	Version  int64  `dynamodbav:"version"`
}

func (prof ddbProfileAdapter) getKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(prof.Id)
	if err != nil {
		panic(err)
	}
	element, err := attributevalue.Marshal(ProfileElement)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"id": id, "element": element}
}

func profileToDDB(profile *profile.Profile) *ddbProfileAdapter {
	return &ddbProfileAdapter{
		Id:       string(profile.Id),
		Element:  ProfileElement,
		Owner:    string(profile.Owner),
		Name:     profile.Name,
		Location: string(profile.Location),
		Level:    profile.Level,
		Gold:     profile.Gold,
		Version:  0,
	}
}
