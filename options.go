package copy

import "os"

// Options specifies optional actions on copying.
type Options struct {

	// OnSymlink can specify what to do on symlink
	OnSymlink func(src string) SymlinkAction

	// OnDirExists can specify what to do when there is a directory already existing in destination.
	OnDirExists func(src, dest string) DirExistsAction

	// Skip can specify which files should be skipped
	Skip func(src string) (bool, error)

	// AddPermission to every entities,
	// NO MORE THAN 0777
	// @OBSOLETE
	// Use `PermissionControl = AddPermission(perm)` instead
	AddPermission os.FileMode

	// PermissionControl can preserve or even add permission to
	// every entries, for example
	//
	//		opt.PermissionControl = AddPermission(0222)
	//
	// See permission_control.go for more detail.
	PermissionControl PermissionControlFunc

	// Sync file after copy.
	// Useful in case when file must be on the disk
	// (in case crash happens, for example),
	// at the expense of some performance penalty
	Sync bool

	// Preserve the atime and the mtime of the entries.
	// On linux we can preserve only up to 1 millisecond accuracy.
	PreserveTimes bool

	// Preserve the uid and the gid of all entries.
	PreserveOwner bool

	// The byte size of the buffer to use for copying files.
	// If zero, the internal default buffer of 32KB is used.
	// See https://golang.org/pkg/io/#CopyBuffer for more information.
	CopyBufferSize uint

	intent struct {
		src  string
		dest string
	}
}

// SymlinkAction represents what to do on symlink.
type SymlinkAction int

const (
	// Deep creates hard-copy of contents.
	Deep SymlinkAction = iota
	// Shallow creates new symlink to the dest of symlink.
	Shallow
	// Skip does nothing with symlink.
	Skip
)

// DirExistsAction represents what to do on dest dir.
type DirExistsAction int

const (
	// Merge preserves or overwrites existing files under the dir (default behavior).
	Merge DirExistsAction = iota
	// Replace deletes all contents under the dir and copy src files.
	Replace
	// Untouchable does nothing for the dir, and leaves it as it is.
	Untouchable
)

// getDefaultOptions provides default options,
// which would be modified by usage-side.
func getDefaultOptions(src, dest string) Options {
	return Options{
		OnSymlink: func(string) SymlinkAction {
			return Shallow // Do shallow copy
		},
		OnDirExists: nil, // Default behavior is "Merge".
		Skip: func(string) (bool, error) {
			return false, nil // Don't skip
		},
		AddPermission:     0,                  // Add nothing
		PermissionControl: PerservePermission, // Just preserve permission
		Sync:              false,              // Do not sync
		PreserveTimes:     false,              // Do not preserve the modification time
		CopyBufferSize:    0,                  // Do not specify, use default bufsize (32*1024)
		intent: struct {
			src  string
			dest string
		}{src, dest},
	}
}

// assureOptions struct, should be called only once.
// All optional values MUST NOT BE nil/zero after assured.
func assureOptions(src, dest string, opts ...Options) Options {
	defopt := getDefaultOptions(src, dest)
	if len(opts) == 0 {
		return defopt
	}
	if opts[0].OnSymlink == nil {
		opts[0].OnSymlink = defopt.OnSymlink
	}
	if opts[0].Skip == nil {
		opts[0].Skip = defopt.Skip
	}
	if opts[0].AddPermission > 0 {
		opts[0].PermissionControl = AddPermission(opts[0].AddPermission)
	} else if opts[0].PermissionControl == nil {
		opts[0].PermissionControl = PerservePermission
	}
	opts[0].intent.src = defopt.intent.src
	opts[0].intent.dest = defopt.intent.dest
	return opts[0]
}
