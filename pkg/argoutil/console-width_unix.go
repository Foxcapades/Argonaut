//go:build aix || darwin || freebsd || linux || netbsd || openbsd || solaris

package argoutil

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/foxcapades/argonaut/v3/internal/log"
)

func GetConsoleWidth() (int, bool) {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		log.DebugLn1("GetConsoleWidth() > tput failed with error: %s", err)
		return 80, false
	}

	width, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		log.DebugLn1("GetConsoleWidth() > string parse failed with error: %s", err)
		return 80, false
	}

	return width, true
}
