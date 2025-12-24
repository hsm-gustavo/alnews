package platform

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func IsUnsupportedPlatform() error {
	switch runtime.GOOS {
		case "windows":
			return fmt.Errorf("windows is not supported")
		case "darwin":
			return fmt.Errorf("mac is not supported")
		default:
			if IsWSL() {
				return fmt.Errorf("windows is not supported")
			}
			return nil
		}
}

func IsWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}