// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package appsearch

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "appsearch", asset.ModuleFieldsPri, AssetAppsearch); err != nil {
		panic(err)
	}
}

// AssetAppsearch returns asset data.
// This is the base64 encoded gzipped contents of module/appsearch.
func AssetAppsearch() string {
	return "eJzs101u2zoQB/C9TzHIPj6AFw8IXhdtgQRFC7RLhRL/lpmIHyFHTnT7glKsyrIV2S4sbcqlbM78OB5T5C09o1qRcC5A+GyzIGLFBVZ0c+cc/agf3iyIPAqIgBWlYLEgkgiZV46VNSv6b0FE9GcCaSvLAguitUIhw6r+/JaM0NhPFgdXDivKvS3d+5MjwfdDdcMFFhzap8fCxdH378bRVM3orKefo2/pep62eu/5kGgkexxff94fpB5K3yVoaOurpAwix8GXPvKcYIrjvo5PQ/GHeF3iBsIlyihephWjv8B9aWFNPvCFE7BxfIZwFLNRGSAprYg3qMurDNWA5bg2Tp1O+xfQzGqtmCfS/r/LVicntpeQtXibBHsv3hrmUHk/pBprkmkb98Ga27Z5lSje/9oXdkfrn7ZD2kW0ac9bxo7/UqI8wF68uf6y/hn+eNCxDVYYUVSsspBgC8NhmdnS8OBWO1DCE8r3UOoUnuyaHIxUJqcnm4ZYp1iyx77jsVnOYRfs4NJmpYbhRCKwtxX8XPRDySgeJlcG89P7jlG4MhJviZByPnOHcBoXMpE2Szy03c7M3qeM8rVQxXziJvso0mPtETZJ+zeoubNtJAOcE5bR/EKzlbsFjFJDtoEWSemk4Pm8+4pR9Dp2k5wA2yRqrK9CcYSzpRTkwV5BHn8le7yUCAcXlItfyt+beKRj0uzc1/L1y+SRwXBRxaOMKxCPMgMl6N531VKWXsQM17iN3X37Qh7BWRNArDQGqkcnXs3ENl/qa58J77bwIkcPrgxpVRQqILNGfnwYj1eGqzPjfeE84o73inS5m5nEmdf46T+JsEmt8PJfA1yJOd4AvwMAAP//P8LP+A=="
}
