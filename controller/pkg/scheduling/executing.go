package scheduling

import (
	"context"
	"log"
)

func ExecWithNodeList(ctx context.Context, funcName string,
	fn func(context.Context, string, ...interface{}) error,
	nodesList []string, options ...interface{}) (errCnt int) {
	for _, node := range nodesList {
		err := fn(ctx, node, options...)
		if err != nil {
			errCnt++
			log.Println(funcName, node, err)
		}
	}
	return errCnt
}

func ExecWithParameterList(ctx context.Context, funcName string,
	fn func(context.Context, string, string, ...interface{}) error,
	nodesList []string, parameterList []string, options ...interface{}) (errCnt int) {
	for i, node := range nodesList {
		err := fn(ctx, node, parameterList[i], options...)
		if err != nil {
			errCnt++
			log.Println(funcName, node, err)
		}
	}
	return errCnt
}
