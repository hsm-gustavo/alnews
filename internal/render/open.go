package render

import "os/exec"

func Open(url string) {
	exec.Command("xdg-open", url).Start()
}