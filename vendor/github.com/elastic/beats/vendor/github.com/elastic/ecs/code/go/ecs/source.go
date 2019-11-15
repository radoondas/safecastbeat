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

// Code generated by scripts/gocodegen.go - DO NOT EDIT.

package ecs

// Source fields describe details about the source of a packet/event.
// Source fields are usually populated in conjunction with destination fields.
type Source struct {
	// Some event source addresses are defined ambiguously. The event will
	// sometimes list an IP, a domain or a unix socket.  You should always
	// store the raw address in the `.address` field.
	// Then it should be duplicated to `.ip` or `.domain`, depending on which
	// one it is.
	Address string `ecs:"address"`

	// IP address of the source.
	// Can be one or multiple IPv4 or IPv6 addresses.
	IP string `ecs:"ip"`

	// Port of the source.
	Port int64 `ecs:"port"`

	// MAC address of the source.
	MAC string `ecs:"mac"`

	// Source domain.
	Domain string `ecs:"domain"`

	// Bytes sent from the source to the destination.
	Bytes int64 `ecs:"bytes"`

	// Packets sent from the source to the destination.
	Packets int64 `ecs:"packets"`

	// Translated ip of source based NAT sessions (e.g. internal client to
	// internet)
	// Typically connections traversing load balancers, firewalls, or routers.
	NatIP string `ecs:"nat.ip"`

	// Translated port of source based NAT sessions. (e.g. internal client to
	// internet)
	// Typically used with load balancers, firewalls, or routers.
	NatPort int64 `ecs:"nat.port"`
}
