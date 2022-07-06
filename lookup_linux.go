package user

import (
	"os"
	"strings"
)

type User struct {
	Uid      string
	Gid      string
	Username string
	Name     string
	HomeDir  string
}

type Group struct {
	Gid  string // group ID
	Name string // group name
}

func Current() (*User, error) { return lookup(os.Getenv("USER"), "") }

func Lookup(username string) (*User, error) { return lookup(username, "") }

func LookupId(uid string) (*User, error) { return lookup("", uid) }

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

func lookup(username, uid string) (*User, error) {

	// cache
	if username != "" {
		if user, has := cache[username]; has {
			return user, nil
		}
	} else if uid != "" {
		if user, has := cache[uid]; has {
			return user, nil
		}
	}

	// read /etc/passwd
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		panic(err)
	}

	// new user object
	user := &User{}
	ls := strings.Split(string(data), "\n")
	for i := 0; i < len(ls); i++ {
		fields := strings.Split(ls[i], ":")
		if len(fields) != 7 {
			continue
		}
		user.Username = fields[0]
		user.Uid = fields[2]
		if username != "" {
			if user.Username != username {
				continue
			}
		} else if user.Uid != uid {
			continue
		}
		user.Gid = fields[3]
		user.HomeDir = fields[5]
		user.Name = user.Username

		// add cache
		cache[user.Name] = user
		cache[user.Uid] = user
		break
	}
	return user, nil
}

var cache = make(map[string]*User)
