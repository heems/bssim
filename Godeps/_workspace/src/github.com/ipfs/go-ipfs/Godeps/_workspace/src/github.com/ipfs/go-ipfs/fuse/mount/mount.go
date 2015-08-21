// package mount provides a simple abstraction around a mount point
package mount

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	goprocess "github.com/heems/go-ipfs/Godeps/_workspace/src/github.com/jbenet/goprocess"

	u "github.com/heems/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-ipfs/util"
)

var log = u.Logger("mount")

var MountTimeout = time.Second * 5

// Mount represents a filesystem mount
type Mount interface {
	// MountPoint is the path at which this mount is mounted
	MountPoint() string

	// Unmounts the mount
	Unmount() error

	// Process returns the mount's Process to be able to link it
	// to other processes. Unmount upon closing.
	Process() goprocess.Process
}

// ForceUnmount attempts to forcibly unmount a given mount.
// It does so by calling diskutil or fusermount directly.
func ForceUnmount(m Mount) error {
	point := m.MountPoint()
	log.Infof("Force-Unmounting %s...", point)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("diskutil", "umount", "force", point)
	case "linux":
		cmd = exec.Command("fusermount", "-u", point)
	default:
		return fmt.Errorf("unmount: unimplemented")
	}

	errc := make(chan error, 1)
	go func() {
		defer close(errc)

		// try vanilla unmount first.
		if err := exec.Command("umount", point).Run(); err == nil {
			return
		}

		// retry to unmount with the fallback cmd
		errc <- cmd.Run()
	}()

	select {
	case <-time.After(2 * time.Second):
		return fmt.Errorf("umount timeout")
	case err := <-errc:
		return err
	}
}

// ForceUnmountManyTimes attempts to forcibly unmount a given mount,
// many times. It does so by calling diskutil or fusermount directly.
// Attempts a given number of times.
func ForceUnmountManyTimes(m Mount, attempts int) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = ForceUnmount(m)
		if err == nil {
			return err
		}

		<-time.After(time.Millisecond * 500)
	}
	return fmt.Errorf("Unmount %s failed after 10 seconds of trying.", m.MountPoint())
}
