package contract

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func findSpecPath() string {
	candidates := []string{
		"specs/openapi.yaml",
		"../specs/openapi.yaml",
		"../../specs/openapi.yaml",
		"./specs/openapi.yaml",
	}
	for _, p := range candidates {
		abs, _ := filepath.Abs(p)
		if _, err := ioutil.ReadFile(p); err == nil {
			_ = abs // silence linter if needed
			return p
		}
	}
	return "specs/openapi.yaml"
}

func TestOpenAPISpec(t *testing.T) {
	path := findSpecPath()
	data, err := ioutil.ReadFile(path)
	assert.NoError(t, err, "should be able to read %s", path)

	var doc map[string]interface{}
	err = yaml.Unmarshal(data, &doc)
	assert.NoError(t, err, "openapi.yaml should be valid YAML")

	// Basic sanity checks
	_, hasOpenAPI := doc["openapi"]
	assert.True(t, hasOpenAPI, "openapi key must be present")

	paths, hasPaths := doc["paths"]
	assert.True(t, hasPaths, "paths must be present in the spec")
	assert.NotNil(t, paths, "paths must not be nil")
}
