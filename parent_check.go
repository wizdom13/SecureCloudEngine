package main

import (
	"os"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

var allowedParentProcessNames = map[string]struct{}{
	"securecloud":     {},
	"securecloud.exe": {},
}

func isAllowedParentProcessName(processName string) bool {
	normalizedName := strings.ToLower(processName)
	if idx := strings.LastIndexAny(normalizedName, `/\\`); idx >= 0 {
		normalizedName = normalizedName[idx+1:]
	}
	_, allowed := allowedParentProcessNames[normalizedName]
	return allowed
}

func isLaunchedBySecureCloud() bool {
	parentProcess, err := process.NewProcess(int32(os.Getppid()))
	if err != nil {
		return false
	}

	parentName, err := parentProcess.Name()
	if err != nil {
		return false
	}

	return isAllowedParentProcessName(parentName)
}
