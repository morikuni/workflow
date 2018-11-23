package workflow_test

import (
	"strings"
	"testing"

	"github.com/morikuni/workflow"
	"github.com/stretchr/testify/assert"
)

func TestParseDefinition(t *testing.T) {
	r := strings.NewReader(`
tasks:
  aaa:
    cmd: bbb
  ccc:
    cmd: ddd
`)
	want := workflow.Definition{
		Tasks: map[string]workflow.Task{
			"aaa": {
				CMD: "bbb",
			},
			"ccc": {
				CMD: "ddd",
			},
		},
	}

	def, err := workflow.ParseDefinition(r)
	assert.NoError(t, err)
	assert.Equal(t, want, def)
}

func TestLoadDefinition(t *testing.T) {
	want := workflow.Definition{
		Tasks: map[string]workflow.Task{
			"echo": {
				CMD: `echo "hello"`,
			},
			"pipe": {
				CMD: `echo "world" | cat`,
			},
		},
	}

	def, err := workflow.LoadDefinition("testdata/simple.yaml")
	assert.NoError(t, err)
	assert.Equal(t, want, def)
}
