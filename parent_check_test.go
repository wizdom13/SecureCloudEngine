package main

import "testing"

func TestIsAllowedParentProcessName(t *testing.T) {
	tests := []struct {
		name        string
		processName string
		want        bool
	}{
		{name: "windows executable name", processName: "SecureCloud.exe", want: true},
		{name: "linux executable name", processName: "SecureCloud", want: true},
		{name: "full windows path", processName: `C:\\Program Files\\SecureCloud\\SecureCloud.exe`, want: true},
		{name: "full unix path", processName: "/opt/securecloud/SecureCloud", want: true},
		{name: "different process", processName: "bash", want: false},
		{name: "empty process name", processName: "", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isAllowedParentProcessName(tc.processName)
			if got != tc.want {
				t.Fatalf("isAllowedParentProcessName(%q) = %v, want %v", tc.processName, got, tc.want)
			}
		})
	}
}
