package svc

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

var (
	maxItems int32 = 1000
)

type LambdaIface interface {
	ListFunctions(context.Context,
		*lambda.ListFunctionsInput,
		...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
	UpdateFunctionConfiguration(
		context.Context,
		*lambda.UpdateFunctionConfigurationInput,
		...func(*lambda.Options)) (*lambda.UpdateFunctionConfigurationOutput, error)
}

func GetClient(cfg aws.Config) *lambda.Client {
	return lambda.NewFromConfig(cfg)
}

type LambdaServiceIface interface {
	FilterFunctionsByRuntime(runtime types.Runtime) ([]string, error)
	ListFunctions() ([]types.FunctionConfiguration, error)
	UpdateRuntime(functionName string) (*types.ImageConfigResponse, error)
}

type LambdaService struct {
	client LambdaIface
	ctx    context.Context
}

func New(ctx context.Context, client *lambda.Client) LambdaServiceIface {
	return &LambdaService{client: client, ctx: ctx}
}

func (l LambdaService) ListFunctions() ([]types.FunctionConfiguration, error) {
	res, err := l.client.ListFunctions(l.ctx, &lambda.ListFunctionsInput{MaxItems: &maxItems})
	if err != nil {
		return nil, err
	}
	return res.Functions, nil
}

func (l LambdaService) FilterFunctionsByRuntime(runtime types.Runtime) ([]string, error) {
	functions, err := l.ListFunctions()
	if err != nil {
		return nil, err
	}
	var results []string
	for _, fn := range functions {
		if fn.Runtime == runtime {
			results = append(results, *fn.FunctionName)
		}
	}
	return results, nil
}

func (l LambdaService) UpdateRuntime(functionName string) (*types.ImageConfigResponse, error) {
	log.Printf("running for function %s\n", functionName)
	res, err := l.client.UpdateFunctionConfiguration(l.ctx, &lambda.UpdateFunctionConfigurationInput{
		FunctionName: &functionName,
		Runtime:      types.RuntimePython311,
	})
	if err != nil {
		return nil, err
	}
	return res.ImageConfigResponse, nil
}

/*
sample usage

func GetServices() lambdasvc.LambdaServiceIface {
	ctx := context.Background()
	cfg, err := awsconfig.GetAwsConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load aws config: %s", err.Error())
	}
	client := lambdasvc.GetClient(cfg)
	return lambdasvc.New(ctx, client)
}

func doTheNeedful() {
	Load()
	lambdaService := GetServices()

	functionNames, err := lambdaService.FilterFunctionsByRuntime(types.RuntimePython310)
	if err != nil {
		log.Fatalf("error listing functions %s\n", err.Error())
	}

	wg := sync.WaitGroup{}
	for _, fn := range functionNames {
		wg.Add(1)
		go func(fn string) {
			defer wg.Done()
			_, err := lambdaService.UpdateRuntime(fn)
			if err != nil {
				log.Fatal(err)
			}
		}(fn)
	}
	wg.Wait()

	log.Println("updates successful")
}

*/
