// This file is part of DiceDB.
// Copyright (C) 2024 DiceDB (dicedb.io).
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package resp

import (
	"testing"
	"time"

	testifyAssert "github.com/stretchr/testify/assert"
)

func TestZCARD(t *testing.T) {
	conn := getLocalConnection()
	defer conn.Close()

	testCases := []struct {
		name   string
		cmds   []string
		expect []interface{}
		delays []time.Duration
	}{
		{
			name: "ZCARD with wrong number of arguments",
			cmds: []string{"ZCARD", "ZCARD myzset field"},
			expect: []interface{}{"ERR wrong number of arguments for 'zcard' command",
				"ERR wrong number of arguments for 'zcard' command"},
			delays: []time.Duration{0, 0},
		},
		{
			name: "ZCARD with wrong type of key",
			cmds: []string{"SET string_key string_value", "ZCARD string_key"},
			expect: []interface{}{"OK",
				"WRONGTYPE Operation against a key holding the wrong kind of value"},
			delays: []time.Duration{0, 0},
		},
		{
			name:   "ZCARD with non-existent key",
			cmds:   []string{"ZADD myzset 1 one", "ZCARD wrong_myzset"},
			expect: []interface{}{int64(1), int64(0)},
			delays: []time.Duration{0, 0},
		},
		{
			name:   "ZCARD with sorted set holding single element",
			cmds:   []string{"ZADD myzset 1 one", "ZCARD myzset"},
			expect: []interface{}{int64(1), int64(1)},
			delays: []time.Duration{0, 0},
		},
		{
			name:   "ZCARD with sorted set holding multiple elements",
			cmds:   []string{"ZADD myzset 1 one 2 two", "ZCARD myzset", "ZADD myzset 3 three", "ZCARD myzset", "ZREM myzset two", "ZCARD myzset"},
			expect: []interface{}{int64(2), int64(2), int64(1), int64(3), int64(1), int64(2)},
			delays: []time.Duration{0, 0, 0, 0, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			FireCommand(conn, "DEL myzset string_key")
			for i, cmd := range tc.cmds {
				if tc.delays[i] > 0 {
					time.Sleep(tc.delays[i])
				}
				result := FireCommand(conn, cmd)
				testifyAssert.Equal(t, tc.expect[i], result, "Value mismatch for cmd %s", cmd)
			}
		})
	}
}
