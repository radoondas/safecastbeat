// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build integration

package template

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/outputs/elasticsearch/estest"
	"github.com/elastic/beats/libbeat/version"
)

type testTemplate struct {
	t      *testing.T
	client ESClient
	common.MapStr
}

func TestCheckTemplate(t *testing.T) {
	client := estest.GetTestingElasticsearch(t)
	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	loader := &Loader{
		client: client,
	}

	// Check for non existent template
	assert.False(t, loader.CheckTemplate("libbeat-notexists"))
}

func TestLoadTemplate(t *testing.T) {
	// Setup ES
	client := estest.GetTestingElasticsearch(t)
	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	// Load template
	absPath, err := filepath.Abs("../")
	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	fieldsPath := absPath + "/fields.yml"
	index := "testbeat"

	tmpl, err := New(version.GetDefaultVersion(), index, client.GetVersion(), TemplateConfig{}, false)
	assert.NoError(t, err)
	content, err := tmpl.LoadFile(fieldsPath)
	assert.NoError(t, err)

	loader := &Loader{
		client: client,
	}

	// Load template
	err = loader.LoadTemplate(tmpl.GetName(), content)
	assert.Nil(t, err)

	// Make sure template was loaded
	assert.True(t, loader.CheckTemplate(tmpl.GetName()))

	// Delete template again to clean up
	client.Request("DELETE", "/_template/"+tmpl.GetName(), "", nil, nil)

	// Make sure it was removed
	assert.False(t, loader.CheckTemplate(tmpl.GetName()))
}

func TestLoadInvalidTemplate(t *testing.T) {
	// Invalid Template
	template := map[string]interface{}{
		"json": "invalid",
	}

	// Setup ES
	client := estest.GetTestingElasticsearch(t)
	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	templateName := "invalidtemplate"

	loader := &Loader{
		client: client,
	}

	// Try to load invalid template
	err := loader.LoadTemplate(templateName, template)
	assert.Error(t, err)

	// Make sure template was not loaded
	assert.False(t, loader.CheckTemplate(templateName))
}

// Tests loading the templates for each beat
func TestLoadBeatsTemplate(t *testing.T) {
	beats := []string{
		"libbeat",
	}

	for _, beat := range beats {
		// Load template
		absPath, err := filepath.Abs("../../" + beat)
		assert.NotNil(t, absPath)
		assert.Nil(t, err)

		// Setup ES
		client := estest.GetTestingElasticsearch(t)
		if err := client.Connect(); err != nil {
			t.Fatal(err)
		}

		fieldsPath := absPath + "/fields.yml"
		index := beat

		tmpl, err := New(version.GetDefaultVersion(), index, client.GetVersion(), TemplateConfig{}, false)
		assert.NoError(t, err)
		content, err := tmpl.LoadFile(fieldsPath)
		assert.NoError(t, err)

		loader := &Loader{
			client: client,
		}

		// Load template
		err = loader.LoadTemplate(tmpl.GetName(), content)
		assert.Nil(t, err)

		// Make sure template was loaded
		assert.True(t, loader.CheckTemplate(tmpl.GetName()))

		// Delete template again to clean up
		client.Request("DELETE", "/_template/"+tmpl.GetName(), "", nil, nil)

		// Make sure it was removed
		assert.False(t, loader.CheckTemplate(tmpl.GetName()))
	}
}

func TestTemplateSettings(t *testing.T) {
	// Setup ES
	client := estest.GetTestingElasticsearch(t)
	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	// Load template
	absPath, err := filepath.Abs("../")
	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	fieldsPath := absPath + "/fields.yml"

	settings := TemplateSettings{
		Index: common.MapStr{
			"number_of_shards": 1,
		},
		Source: common.MapStr{
			"enabled": false,
		},
	}
	config := TemplateConfig{
		Settings: settings,
	}
	tmpl, err := New(version.GetDefaultVersion(), "testbeat", client.GetVersion(), config, false)
	assert.NoError(t, err)
	content, err := tmpl.LoadFile(fieldsPath)
	assert.NoError(t, err)

	loader := &Loader{
		client: client,
	}

	// Load template
	err = loader.LoadTemplate(tmpl.GetName(), content)
	assert.Nil(t, err)

	// Check that it contains the mapping
	templateJSON := getTemplate(t, client, tmpl.GetName())
	assert.Equal(t, 1, templateJSON.NumberOfShards())
	assert.Equal(t, false, templateJSON.SourceEnabled())

	// Delete template again to clean up
	client.Request("DELETE", "/_template/"+tmpl.GetName(), "", nil, nil)

	// Make sure it was removed
	assert.False(t, loader.CheckTemplate(tmpl.GetName()))
}

func TestOverwrite(t *testing.T) {
	// Setup ES
	client := estest.GetTestingElasticsearch(t)
	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	beatInfo := beat.Info{
		Beat:        "testbeat",
		IndexPrefix: "testbeatidx",
		Version:     version.GetDefaultVersion(),
	}
	templateName := "testbeatidx-" + version.GetDefaultVersion()

	absPath, err := filepath.Abs("../")
	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	// make sure no template is already there
	client.Request("DELETE", "/_template/"+templateName, "", nil, nil)

	// Load template
	config := TemplateConfig{
		Enabled: true,
		Fields:  absPath + "/fields.yml",
	}
	loader, err := NewLoader(config, client, beatInfo, nil, false)
	assert.NoError(t, err)
	err = loader.Load()
	assert.NoError(t, err)

	// Load template again, this time with custom settings
	config = TemplateConfig{
		Enabled: true,
		Fields:  absPath + "/fields.yml",
		Settings: TemplateSettings{
			Source: map[string]interface{}{
				"enabled": false,
			},
		},
	}
	loader, err = NewLoader(config, client, beatInfo, nil, false)
	assert.NoError(t, err)
	err = loader.Load()
	assert.NoError(t, err)

	// Overwrite was not enabled, so the first version should still be there
	templateJSON := getTemplate(t, client, templateName)
	assert.Equal(t, true, templateJSON.SourceEnabled())

	// Load template again, this time with custom settings AND overwrite: true
	config = TemplateConfig{
		Enabled:   true,
		Overwrite: true,
		Fields:    absPath + "/fields.yml",
		Settings: TemplateSettings{
			Source: map[string]interface{}{
				"enabled": false,
			},
		},
	}
	loader, err = NewLoader(config, client, beatInfo, nil, false)
	assert.NoError(t, err)
	err = loader.Load()
	assert.NoError(t, err)

	// Overwrite was enabled, so the custom setting should be there
	templateJSON = getTemplate(t, client, templateName)
	assert.Equal(t, false, templateJSON.SourceEnabled())

	// Delete template again to clean up
	client.Request("DELETE", "/_template/"+templateName, "", nil, nil)
}

var dataTests = []struct {
	data  common.MapStr
	error bool
}{
	{
		data: common.MapStr{
			"keyword": "test keyword",
			"array":   [...]int{1, 2, 3},
			"object": common.MapStr{
				"hello": "world",
			},
		},
		error: false,
	},
	{
		// Invalid array
		data: common.MapStr{
			"array": common.MapStr{
				"hello": "world",
			},
		},
		error: true,
	},
	{
		// Invalid object
		data: common.MapStr{
			"object": [...]int{1, 2, 3},
		},
		error: true,
	},
	{
		// tests enabled: false values
		data: common.MapStr{
			"array_disabled": [...]int{1, 2, 3},
			"object_disabled": common.MapStr{
				"hello": "world",
			},
		},
		error: false,
	},
}

// Tests if data can be loaded into elasticsearch with right types
func TestTemplateWithData(t *testing.T) {
	fieldsPath, err := filepath.Abs("./testdata/fields.yml")
	assert.NotNil(t, fieldsPath)
	assert.Nil(t, err)

	// Setup ES
	client := estest.GetTestingElasticsearch(t)

	tmpl, err := New(version.GetDefaultVersion(), "testindex", client.GetVersion(), TemplateConfig{}, false)
	assert.NoError(t, err)
	content, err := tmpl.LoadFile(fieldsPath)
	assert.NoError(t, err)

	loader := &Loader{
		client: client,
	}

	// Load template
	err = loader.LoadTemplate(tmpl.GetName(), content)
	assert.Nil(t, err)

	// Make sure template was loaded
	assert.True(t, loader.CheckTemplate(tmpl.GetName()))

	for _, test := range dataTests {
		_, _, err = client.Index(tmpl.GetName(), "_doc", "", nil, test.data)
		if test.error {
			assert.NotNil(t, err)

		} else {
			assert.Nil(t, err)
		}
	}

	// Delete template again to clean up
	client.Request("DELETE", "/_template/"+tmpl.GetName(), "", nil, nil)

	// Make sure it was removed
	assert.False(t, loader.CheckTemplate(tmpl.GetName()))
}

func getTemplate(t *testing.T, client ESClient, templateName string) testTemplate {
	status, body, err := client.Request("GET", "/_template/"+templateName, "", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, status, 200)

	var response common.MapStr
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	return testTemplate{
		t:      t,
		client: client,
		MapStr: common.MapStr(response[templateName].(map[string]interface{})),
	}
}

func (tt *testTemplate) SourceEnabled() bool {
	key := fmt.Sprintf("mappings._source.enabled")

	// _source.enabled is true if it's missing (default)
	b, _ := tt.HasKey(key)
	if !b {
		return true
	}

	val, err := tt.GetValue(key)
	if !assert.NoError(tt.t, err) {
		doc, _ := json.MarshalIndent(tt.MapStr, "", "    ")
		tt.t.Fatal(fmt.Sprintf("failed to read '%v' in %s", key, doc))
	}

	return val.(bool)
}

func (tt *testTemplate) NumberOfShards() int {
	val, err := tt.GetValue("settings.index.number_of_shards")
	require.NoError(tt.t, err)

	i, err := strconv.Atoi(val.(string))
	require.NoError(tt.t, err)
	return i
}
