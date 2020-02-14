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

package cmd

import (
	"github.com/elastic/beats/libbeat/cmd"
	"github.com/elastic/beats/libbeat/cmd/instance"
	"github.com/elastic/beats/winlogbeat/beater"

	// Register fields.
	_ "github.com/elastic/beats/winlogbeat/include"

	// Import processors and supporting modules.
	_ "github.com/elastic/beats/libbeat/processors/script"
	_ "github.com/elastic/beats/libbeat/processors/timestamp"
	_ "github.com/elastic/beats/winlogbeat/processors/script/javascript/module/winlogbeat"
)

// Name of this beat
var Name = "winlogbeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, instance.Settings{Name: Name, HasDashboards: true})
