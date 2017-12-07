package datastore

import "github.com/aws/aws-sdk-go/aws/credentials"

type Datastore interface {
	AddUser(UserInfo)
}

type MockDB struct {
	User []UserInfo
}

type UserInfo struct {
	ID    string
	UName string
}

type DynamoDB struct {
	//figure out later
}

func NewMockDB() (*MockDB, error) {
	m := MockDB{
		User: []UserInfo{},
	}
	return &m, nil
}

func (m *MockDB) AddUser(u UserInfo) {
	m.User = append(m.User, u)
}

func NewDynamoDB(creds *credentials.Credentials) (*DynamoDB, error) {
	return nil, nil
}
