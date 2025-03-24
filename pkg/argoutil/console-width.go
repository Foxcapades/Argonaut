//go:build !(aix || darwin || freebsd || linux || netbsd || openbsd || solaris)

package argoutil

func GetConsoleWidth() (int, bool) {
	return 80, false
}
