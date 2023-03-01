package svc

type StringMap = map[string]string

type ServiceResponse struct {
	IsOk  bool
	Value string
	Error error
}

type XLDto struct {
	XlData []StringMap
	DbName string
}
