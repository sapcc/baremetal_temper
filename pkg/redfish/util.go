/**
 * Copyright 2021 SAP SE
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redfish

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/stmcginnis/gofish/redfish"
)

func parseMac(s string, sep rune) (string, error) {
	if len(s) < 12 {
		return "", fmt.Errorf("invalid MAC address: %s", s)
	}
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, "-", "")
	var buf bytes.Buffer
	for i, char := range s {
		buf.WriteRune(char)
		if i%2 == 1 && i != len(s)-1 {
			buf.WriteRune(sep)
		}

	}

	return buf.String(), nil
}

func calcTotalMemory(mem []*redfish.Memory) (totalMem int) {
	totalMem = 0
	for _, m := range mem {
		totalMem = totalMem + m.CapacityMiB
	}
	return
}
