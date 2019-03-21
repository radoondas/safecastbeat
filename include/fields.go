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

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package include

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("safecastbeat", "fields.yml", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
