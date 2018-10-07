package system 

import (
	"os"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/mount"
	"github.com/pkg/errors"
)













func EnsureRemoveAll(dir string) error {
	notExistErr := make(map[string]bool)

	
	exitOnErr := make(map[string]int)
	maxRetry := 50

	
	mount.RecursiveUnmount(dir)

	for {
		err := os.RemoveAll(dir)
		if err == nil {
			return nil
		}

		pe, ok := err.(*os.PathError)
		if !ok {
			return err
		}

		if os.IsNotExist(err) {
			if notExistErr[pe.Path] {
				return err
			}
			notExistErr[pe.Path] = true

			
			
			
			
			
			if pe.Path == dir {
				return nil
			}
			continue
		}

		if pe.Err != syscall.EBUSY {
			return err
		}

		if mounted, _ := mount.Mounted(pe.Path); mounted {
			if e := mount.Unmount(pe.Path); e != nil {
				if mounted, _ := mount.Mounted(pe.Path); mounted {
					return errors.Wrapf(e, "error while removing %s", dir)
				}
			}
		}

		if exitOnErr[pe.Path] == maxRetry {
			return err
		}
		exitOnErr[pe.Path]++
		time.Sleep(100 * time.Millisecond)
	}
}
