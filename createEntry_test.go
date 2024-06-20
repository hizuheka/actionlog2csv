package main

import (
	"fmt"
	"testing"
)

// TestCreateEntry 関数のテスト
func TestCreateEntry(t *testing.T) {
	testCases := []struct {
		name          string
		inputLine     string
		expectedEntry LogEntry
		expectedError error
	}{
		{
			name:          "Case 0",
			inputLine:     "2024 May 28 14:12:01 : firewall INFO[55555555]: TCP connection initiated.  src=192.100.1.200 dst=192.100.2.244 proto=tcp srcport=88888 dstport=77777 interface=bnd1 dir=inbound action=accept rule=123 time=2024-05-28T14:12:01",
			expectedEntry: LogEntry{Src: "192.100.1.200", Dst: "192.100.2.244", Interface: "bnd1", Dir: "inbound", Action: "accept", Rule: "123"},
			expectedError: nil,
		},
		{
			name:          "Case 1: Basic input",
			inputLine:     "src=192.168.1.1 dst=192.168.1.2 interface=eth0 dir=outbound action=reject rule=999",
			expectedEntry: LogEntry{Src: "192.168.1.1", Dst: "192.168.1.2", Interface: "eth0", Dir: "outbound", Action: "reject", Rule: "999"},
			expectedError: nil,
		},
		{
			name:          "Case 2: Missing src field",
			inputLine:     "dst=10.0.0.2 interface=eth1",
			expectedEntry: LogEntry{},
			expectedError: fmt.Errorf("ログファイルの形式が不正です: src=, dst=10.0.0.2, interface=eth1, dir=, action=, rule="),
		},
		{
			name:          "Case 3: Missing dst field",
			inputLine:     "src=172.16.0.1 interface=eth2",
			expectedEntry: LogEntry{},
			expectedError: fmt.Errorf("ログファイルの形式が不正です: src=172.16.0.1, dst=, interface=eth2, dir=, action=, rule="),
		},
		{
			name:          "Case 4: Missing interface field",
			inputLine:     "src=8.8.8.8 dst=8.8.4.4",
			expectedEntry: LogEntry{},
			expectedError: fmt.Errorf("ログファイルの形式が不正です: src=8.8.8.8, dst=8.8.4.4, interface=, dir=, action=, rule="),
		},
		{
			name:          "Case 5: Empty input line",
			inputLine:     "",
			expectedEntry: LogEntry{},
			expectedError: fmt.Errorf("ログファイルの形式が不正です: src=, dst=, interface=, dir=, action=, rule="),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := createEntry(tc.inputLine)

			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("Expected no error, but got error: %v", err)
				}
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected error message '%v', but got '%v'", tc.expectedError, err)
				}
			} else {
				if tc.expectedError != nil {
					t.Errorf("Expected error '%v', but got no error", tc.expectedError)
				}
				if entry != tc.expectedEntry {
					t.Errorf("Expected entry '%v', but got '%v'", tc.expectedEntry, entry)
				}
			}
		})
	}
}
