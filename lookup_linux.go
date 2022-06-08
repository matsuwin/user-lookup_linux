package user

import (
	"io/ioutil"
	"os"
	"strings"
)

// User represents a user account.
type User struct {
	// Uid is the user ID.
	// On POSIX systems, this is a decimal number representing the uid.
	// On Windows, this is a security identifier (SID) in a string format.
	// On Plan 9, this is the contents of /dev/user.
	Uid string
	// Gid is the primary group ID.
	// On POSIX systems, this is a decimal number representing the gid.
	// On Windows, this is a SID in a string format.
	// On Plan 9, this is the contents of /dev/user.
	Gid string
	// Username is the login name.
	Username string
	// Name is the user's real or display name.
	// It might be blank.
	// On POSIX systems, this is the first (or only) entry in the GECOS field
	// list.
	// On Windows, this is the user's display name.
	// On Plan 9, this is the contents of /dev/user.
	Name string
	// HomeDir is the path to the user's home directory (if they have one).
	HomeDir string
}

// Group represents a grouping of users.
//
// On POSIX systems Gid contains a decimal number representing the group ID.
type Group struct {
	Gid  string // group ID
	Name string // group name
}

// Lookup looks up a user by username. If the user cannot be found, the
// returned error is of type UnknownUserError.
func Lookup(username string) (*User, error) {
	return lookup(username, "")
}

// LookupId looks up a user by userid. If the user cannot be found, the
// returned error is of type UnknownUserIdError.
func LookupId(uid string) (*User, error) {
	return lookup("", uid)
}

// LookupGroup looks up a group by name. If the group cannot be found, the
// returned error is of type UnknownGroupError.
func LookupGroup(name string) (*Group, error) {
	u, err := lookup(name, "")
	if err != nil {
		return nil, err
	}
	g := &Group{}
	g.Gid = u.Gid
	g.Name = u.Username
	return g, nil
}

// LookupGroupId looks up a group by groupid. If the group cannot be found, the
// returned error is of type UnknownGroupIdError.
func LookupGroupId(gid string) (*Group, error) {
	u, err := lookup("", gid)
	if err != nil {
		return nil, err
	}
	g := &Group{}
	g.Gid = u.Gid
	g.Name = u.Username
	return g, nil
}

// Current returns the current user.
//
// The first call will cache the current user information.
// Subsequent calls will return the cached value and will not reflect
// changes to the current user.
func Current() (*User, error) {
	return lookup(os.Getenv("USER"), "")
}

func lookup(username, uid string) (*User, error) {
	if username != "" {
		if user, has := cache[username]; has {
			return user, nil
		}
	} else if uid != "" {
		if user, has := cache[uid]; has {
			return user, nil
		}
	}
	u := &User{}
	ls := strings.Split(cat("/etc/passwd"), "\n")
	for i := 0; i < len(ls); i++ {
		vs := strings.Split(ls[i], ":")
		if len(vs) != 7 {
			continue
		}
		u.Username = vs[0]
		u.Uid = vs[2]
		if username != "" {
			if u.Username != username {
				continue
			}
		} else if u.Uid != uid {
			continue
		}
		u.Gid = vs[3]
		u.HomeDir = vs[5]
		u.Name = u.Username

		cache[u.Name] = u
		cache[u.Uid] = u

		break
	}
	return u, nil
}

func cat(fp string) string {
	data, _ := ioutil.ReadFile(fp)
	return string(data)
}

var cache = make(map[string]*User)
