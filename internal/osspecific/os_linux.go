//go:build linux

package osspecific

import (
	"os"

	"golang.org/x/sys/unix"
)

//-----------------------------------------------------------------------------
// Set the terminal echo on or off.
//
// If the teminal is redirected "IoctlGetTermios" will return "unix.ENOTTY".
// In that case, the caller should ignore the error and continue. This is
// useful for example when running in a docker container
//-----------------------------------------------------------------------------

func SetEcho(enable bool) error {

	fd := int(os.Stdin.Fd())

	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)

	if err != nil {

		return err
	}

	if enable {

		termios.Lflag |= unix.ECHO

	} else {

		termios.Lflag &^= unix.ECHO
	}

	return unix.IoctlSetTermios(fd, unix.TCSETS, termios)
}

func FlushStdin() error {

	return unix.IoctlSetInt(int(os.Stdin.Fd()), unix.TCFLSH, unix.TCIFLUSH)
}
