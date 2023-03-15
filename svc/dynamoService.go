package svc

import (
	"encoding/json"
	"fmt"
	"xl-app/db"
	"xl-app/types"
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

func (d *DynamoService) SaveEntity(xlDto types.XLDto) types.ServiceResponse {
	fmt.Println("service SaveEntity called...")
	logJson(xlDto)
	if err := d.db.CreateTable(xlDto.DbName); err != nil {
		return types.ServiceResponse{
			IsOk:  false,
			Error: err,
		}
	}

	err := d.db.BulkSave(xlDto)
	if err != nil {
		return types.ServiceResponse{
			IsOk:  false,
			Error: err,
		}
	}
	return types.ServiceResponse{
		IsOk:  true,
		Value: "good mock response...figure this out later",
		Error: nil,
	}
}

func (d *DynamoService) BulkMigration(name string, json string) types.ServiceResponse {
	println(json)
	return types.ServiceResponse{
		IsOk:  true,
		Value: "mock response...no idea yet...",
		Error: nil,
	}
}
