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
			name:          "Case 1: Basic input",
			inputLine:     "src=192.168.1.1 dest=192.168.1.2 interface=eth0",
			expectedEntry: LogEntry{Src: "192.168.1.1", Dest: "192.168.1.2", Interface: "eth0"},
			expectedError: nil,
		},
		{
			name:          "Case 2: Missing src field",
			inputLine:     "dest=10.0.0.2 interface=eth1",
			expectedEntry: LogEntry{},
			expectedError: fmt.Errorf("ログファイルの解析でエラーが発生しました。ログファイルの形式が不正です: dest=10.0.0.2 interface=eth1"),
		},
		{
			name:          "Case 3: Missing dest field",
			inputLine:     "src=172.16.0.1 interface=eth2",
			expectedEntry: LogEntry{},
			expectedError: fmt.Errorf("ログファイルの解析でエラーが発生しました。ログファイルの形式が不正です: src=172.16.0.1 interface=eth2"),
		},
		{
			name:          "Case 4: Missing interface field",
			inputLine:     "src=8.8.8.8 dest=8.8.4.4",
			expectedEntry: LogEntry{},
			expectedError: fmt.Errorf("ログファイルの解析でエラーが発生しました。ログファイルの形式が不正です: src=8.8.8.8 dest=8.8.4.4"),
		},
		{
			name:          "Case 5: Empty input line",
			inputLine:     "",
			expectedEntry: LogEntry{},
			expectedError: fmt.Errorf("ログファイルの解析でエラーが発生しました。ログファイルの形式が不正です: "),
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
