package user

import (
	"errors"
	"syscall"
)

var (
	
	ErrUnsupported = errors.New("user lookup: operating system does not provide passwd-formatted data")
	
	ErrNoPasswdEntries = errors.New("no matching entries in passwd file")
	ErrNoGroupEntries  = errors.New("no matching entries in group file")
)

func lookupUser(filter func(u User) bool) (User, error) {
	
	passwd, err := GetPasswd()
	if err != nil {
		return User{}, err
	}
	defer passwd.Close()

	
	users, err := ParsePasswdFilter(passwd, filter)
	if err != nil {
		return User{}, err
	}

	
	if len(users) == 0 {
		return User{}, ErrNoPasswdEntries
	}

	
	return users[0], nil
}




func CurrentUser() (User, error) {
	return LookupUid(syscall.Getuid())
}




func LookupUser(username string) (User, error) {
	return lookupUser(func(u User) bool {
		return u.Name == username
	})
}




func LookupUid(uid int) (User, error) {
	return lookupUser(func(u User) bool {
		return u.Uid == uid
	})
}

func lookupGroup(filter func(g Group) bool) (Group, error) {
	
	group, err := GetGroup()
	if err != nil {
		return Group{}, err
	}
	defer group.Close()

	
	groups, err := ParseGroupFilter(group, filter)
	if err != nil {
		return Group{}, err
	}

	
	if len(groups) == 0 {
		return Group{}, ErrNoGroupEntries
	}

	
	return groups[0], nil
}




func CurrentGroup() (Group, error) {
	return LookupGid(syscall.Getgid())
}




func LookupGroup(groupname string) (Group, error) {
	return lookupGroup(func(g Group) bool {
		return g.Name == groupname
	})
}




func LookupGid(gid int) (Group, error) {
	return lookupGroup(func(g Group) bool {
		return g.Gid == gid
	})
}
