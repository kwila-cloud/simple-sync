package contract

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestOpenAPISpec(t *testing.T) {
	data, err := ioutil.ReadFile("specs/openapi.yaml")
	assert.NoError(t, err, "should be able to read specs/openapi.yaml")

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
