package main

import (
	"fmt"
	"testing"
)

// TestCreateEntry 関数のテスト
func TestCreateEntry(t *testing.T) {
	tests := []struct {
		line        string
		expected    LogEntry
		expectError error
	}{
		{
			line:        "2024 May 28 14:12:01 : firewall INFO[55555555]: TCP connection initiated.  src=192.100.1.200 dst=192.100.2.244 proto=tcp srcport=88888 dstport=77777 interface=bnd1 dir=inbound action=accept rule=123 time=2024-05-28T14:12:01",
			expected:    LogEntry{Src: "192.100.1.200", Dst: "192.100.2.244", Interface: "bnd1", Dir: "inbound", Action: "accept", Rule: "123"},
			expectError: nil,
		},
		{
			line: "src=192.168.1.1 dst=192.168.1.2 interface=eth0 dir=in action=allow rule=1",
			expected: LogEntry{
				Src:       "192.168.1.1",
				Dst:       "192.168.1.2",
				Interface: "eth0",
				Dir:       "in",
				Action:    "allow",
				Rule:      "1",
			},
			expectError: nil,
		},
		{
			line:        "src=192.168.1.1 dst=192.168.1.2 interface=eth0 dir=in action=allow",
			expected:    LogEntry{},
			expectError: fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", "192.168.1.1", "192.168.1.2", "eth0", "in", "allow", ""),
		},
		{
			line:        "",
			expected:    LogEntry{},
			expectError: fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", "", "", "", "", "", ""),
		},
		{
			line: "src=10.0.0.1 dst=10.0.0.2 interface=wlan0 dir=out action=deny rule=2",
			expected: LogEntry{
				Src:       "10.0.0.1",
				Dst:       "10.0.0.2",
				Interface: "wlan0",
				Dir:       "out",
				Action:    "deny",
				Rule:      "2",
			},
			expectError: nil,
		},
		{
			line:        "dst=192.168.1.2 interface=eth0 dir=in action=allow rule=1",
			expected:    LogEntry{},
			expectError: fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", "", "192.168.1.2", "eth0", "in", "allow", "1"),
		},
		{
			line:        "src=192.168.1.1 interface=eth0 dir=in action=allow rule=1",
			expected:    LogEntry{},
			expectError: fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", "192.168.1.1", "", "eth0", "in", "allow", "1"),
		},
		{
			line:        "src=192.168.1.1 dst=192.168.1.2 dir=in action=allow rule=1",
			expected:    LogEntry{},
			expectError: fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", "192.168.1.1", "192.168.1.2", "", "in", "allow", "1"),
		},
		{
			line:        "src=192.168.1.1 dst=192.168.1.2 interface=eth0 action=allow rule=1",
			expected:    LogEntry{},
			expectError: fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", "192.168.1.1", "192.168.1.2", "eth0", "", "allow", "1"),
		},
		{
			line:        "src=192.168.1.1 dst=192.168.1.2 interface=eth0 dir=in rule=1",
			expected:    LogEntry{},
			expectError: fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", "192.168.1.1", "192.168.1.2", "eth0", "in", "", "1"),
		},
		{
			line:        "src=192.168.1.1 dst=192.168.1.2 interface=eth0 dir=in action=allow",
			expected:    LogEntry{},
			expectError: fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", "192.168.1.1", "192.168.1.2", "eth0", "in", "allow", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			actual, err := createEntry(tt.line)
			if tt.expectError == nil && err != nil {
				t.Errorf("expected no error, but got: %v", err)
			}
			if tt.expectError != nil && err == nil {
				t.Errorf("expected error: %v, but got none", tt.expectError)
			}
			if tt.expectError != nil && err != nil && err.Error() != tt.expectError.Error() {
				t.Errorf("expected error: %v, but got: %v", tt.expectError, err)
			}
			if actual != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, actual)
			}
		})
	}
}
