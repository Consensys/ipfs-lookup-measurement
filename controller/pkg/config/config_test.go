package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testNodesListFileName = "./testing.out"
)

func TestGetNodesList(t *testing.T) {
	createTestingFile()
	actual := GetNodesList(testNodesListFileName)
	assert.Equal(t, []string{"test1", "test2"}, actual)
}

func createTestingFile() {
	nodesList := []byte(`test1

	test2
	`)
	_ = os.WriteFile(testNodesListFileName, nodesList, 0644)
}
