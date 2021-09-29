package scheduling

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecWithNodeList(t *testing.T) {
	ctx := context.Background()
	nodesList := []string{"n1", "n2"}
	testingFn := func(ctx context.Context, node string, options ...interface{}) error { return nil }
	errCnt := ExecWithNodeList(ctx, "testingFn", testingFn, nodesList)
	assert.Equal(t, 0, errCnt)

	testingFn = func(ctx context.Context, node string, options ...interface{}) error {
		return errors.New("testing error")
	}
	errCnt = ExecWithNodeList(ctx, "testingFn", testingFn, nodesList)
	assert.Equal(t, 2, errCnt)
}

func TestExecWithParameterList(t *testing.T) {
	ctx := context.Background()
	nodesList := []string{"n1", "n2"}
	parameterList := []string{"n1", "n2"}
	testingFn := func(ctx context.Context, node string, para string, options ...interface{}) error { return nil }
	errCnt := ExecWithParameterList(ctx, "testingFn", testingFn, nodesList, parameterList)
	assert.Equal(t, 0, errCnt)

	testingFn = func(ctx context.Context, node string, para string, options ...interface{}) error {
		return errors.New("testing error")
	}
	errCnt = ExecWithParameterList(ctx, "testingFn", testingFn, nodesList, parameterList)
	assert.Equal(t, 2, errCnt)
}
