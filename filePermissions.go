package fsutil

import (
	"github.com/develar/errors"
	"github.com/phayes/permbits"
	"os"
)

// https://github.com/electron-userland/electron-builder/issues/3452#issuecomment-438619535
// quite a lot sources don't have proper permissions to be distributed
func FixPermissions(filePath string, fileMode os.FileMode, isForceSetIfExecutable bool) (permbits.PermissionBits, permbits.PermissionBits, error) {
	originalPermissions := permbits.PermissionBits(fileMode)
	permissions := originalPermissions

	if originalPermissions.UserExecute() {
		permissions.SetGroupExecute(true)
		permissions.SetOtherExecute(true)
	}

	permissions.SetUserRead(true)
	permissions.SetGroupRead(true)
	permissions.SetOtherRead(true)

	permissions.SetSetuid(false)
	permissions.SetSetgid(false)

	if originalPermissions == permissions && (!originalPermissions.UserExecute() || !isForceSetIfExecutable) {
		return originalPermissions, permissions, nil
	}

	err := permbits.Chmod(filePath, permissions)
	return originalPermissions, permissions, errors.WithStack(err)
}
