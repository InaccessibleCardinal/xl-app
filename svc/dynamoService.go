package svc

import (
	"encoding/json"
	"fmt"
	"xl-app/db"
)

type DynamoService struct {
	db *db.Db
}

func logJson(st any) {
	jn, _ := json.Marshal(st)
	println(string(jn))
}

func NewDynamoService(db *db.Db) *DynamoService {
	return &DynamoService{db: db}
}

func (d *DynamoService) SaveEntity(xlDto XLDto) ServiceResponse {
	fmt.Println("service SaveEntity called...")
	logJson(xlDto)
	d.db.CreateTable(xlDto.DbName)
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
