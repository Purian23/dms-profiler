package launch

import (
	"os"
	"os/exec"
	"syscall"
)

// Args is argv for: browser --profile-directory=<dir> <url>.
func Args(browserCmd, profileDir, rawURL string) []string {
	return []string{browserCmd, "--profile-directory=" + profileDir, rawURL}
}

func Exec(argv0 string, argv []string, env []string) error {
	return syscall.Exec(argv0, argv, env)
}

func LookPath(cmd string) (string, error) {
	return exec.LookPath(cmd)
}

func Environ() []string {
	return os.Environ()
}
