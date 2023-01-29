package svc

import "encoding/json"

type ServiceResponse struct {
	IsOk  bool
	Value string
	Error error
}

type DynamoService struct{}

func logJson(st any) {
	jn, _ := json.Marshal(st)
	println(string(jn))
}

func NewDynamoService() *DynamoService {
	return &DynamoService{}
}

func (d *DynamoService) SaveEntity(stringMap map[string]string) ServiceResponse {
	logJson(stringMap)
	return ServiceResponse{
		IsOk:  true,
		Value: "mock response...figure this out later",
		Error: nil,
	}
}

func (d *DynamoService) BulkMigration(name string, json string) ServiceResponse {
	println(json)
	return ServiceResponse{
		IsOk:  true,
		Value: "mock response...no idea yet...",
		Error: nil,
	}
}
