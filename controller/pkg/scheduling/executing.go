package scheduling

import (
	"context"
	"log"
)

func ExecWithNodeList(ctx context.Context, funcName string,
	fn func(context.Context, string) error,
	nodesList []string) (errCnt int) {
	for _, node := range nodesList {
		err := fn(ctx, node)
		if err != nil {
			errCnt++
			log.Println(funcName, node, err)
		}
	}
	return errCnt
}

func ExecWithParameterList(ctx context.Context, funcName string,
	fn func(context.Context, string, string) error,
	nodesList []string, parameterList []string) (errCnt int) {
	for i, node := range nodesList {
		err := fn(ctx, node, parameterList[i])
		if err != nil {
			errCnt++
			log.Println(funcName, node, err)
		}
	}
	return errCnt
}
