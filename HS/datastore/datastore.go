package datastore

import "github.com/aws/aws-sdk-go/aws/credentials"

type Datastore interface {
	//MUST have a pointer to client
}

type MockDB struct {
	//figure out later
}

type DynamoDB struct {
	//figure out later
}

func NewMockDB() (*MockDB, error) {
	return nil, nil
}

func NewDynamoDB(creds *credentials.Credentials) (*DynamoDB, error) {
	return nil, nil
}
