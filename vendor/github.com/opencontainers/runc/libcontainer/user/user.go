package user

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	minId = 0
	maxId = 1<<31 - 1 
)

var (
	ErrRange = fmt.Errorf("uids and gids must be in range %d-%d", minId, maxId)
)

type User struct {
	Name  string
	Pass  string
	Uid   int
	Gid   int
	Gecos string
	Home  string
	Shell string
}

type Group struct {
	Name string
	Pass string
	Gid  int
	List []string
}

func parseLine(line string, v ...interface{}) {
	if line == "" {
		return
	}

	parts := strings.Split(line, ":")
	for i, p := range parts {
		
		
		if len(v) <= i {
			break
		}

		
		
		switch e := v[i].(type) {
		case *string:
			*e = p
		case *int:
			
			*e, _ = strconv.Atoi(p)
		case *[]string:
			
			if p != "" {
				*e = strings.Split(p, ",")
			} else {
				*e = []string{}
			}
		default:
			
			panic(fmt.Sprintf("parseLine only accepts {*string, *int, *[]string} as arguments! %#v is not a pointer!", e))
		}
	}
}

func ParsePasswdFile(path string) ([]User, error) {
	passwd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer passwd.Close()
	return ParsePasswd(passwd)
}

func ParsePasswd(passwd io.Reader) ([]User, error) {
	return ParsePasswdFilter(passwd, nil)
}

func ParsePasswdFileFilter(path string, filter func(User) bool) ([]User, error) {
	passwd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer passwd.Close()
	return ParsePasswdFilter(passwd, filter)
}

func ParsePasswdFilter(r io.Reader, filter func(User) bool) ([]User, error) {
	if r == nil {
		return nil, fmt.Errorf("nil source for passwd-formatted data")
	}

	var (
		s   = bufio.NewScanner(r)
		out = []User{}
	)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		
		
		
		
		
		p := User{}
		parseLine(line, &p.Name, &p.Pass, &p.Uid, &p.Gid, &p.Gecos, &p.Home, &p.Shell)

		if filter == nil || filter(p) {
			out = append(out, p)
		}
	}

	return out, nil
}

func ParseGroupFile(path string) ([]Group, error) {
	group, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer group.Close()
	return ParseGroup(group)
}

func ParseGroup(group io.Reader) ([]Group, error) {
	return ParseGroupFilter(group, nil)
}

func ParseGroupFileFilter(path string, filter func(Group) bool) ([]Group, error) {
	group, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer group.Close()
	return ParseGroupFilter(group, filter)
}

func ParseGroupFilter(r io.Reader, filter func(Group) bool) ([]Group, error) {
	if r == nil {
		return nil, fmt.Errorf("nil source for group-formatted data")
	}

	var (
		s   = bufio.NewScanner(r)
		out = []Group{}
	)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		text := s.Text()
		if text == "" {
			continue
		}

		
		
		
		
		
		p := Group{}
		parseLine(text, &p.Name, &p.Pass, &p.Gid, &p.List)

		if filter == nil || filter(p) {
			out = append(out, p)
		}
	}

	return out, nil
}

type ExecUser struct {
	Uid   int
	Gid   int
	Sgids []int
	Home  string
}





func GetExecUserPath(userSpec string, defaults *ExecUser, passwdPath, groupPath string) (*ExecUser, error) {
	passwd, err := os.Open(passwdPath)
	if err != nil {
		passwd = nil
	} else {
		defer passwd.Close()
	}

	group, err := os.Open(groupPath)
	if err != nil {
		group = nil
	} else {
		defer group.Close()
	}

	return GetExecUser(userSpec, defaults, passwd, group)
}






















func GetExecUser(userSpec string, defaults *ExecUser, passwd, group io.Reader) (*ExecUser, error) {
	if defaults == nil {
		defaults = new(ExecUser)
	}

	
	user := &ExecUser{
		Uid:   defaults.Uid,
		Gid:   defaults.Gid,
		Sgids: defaults.Sgids,
		Home:  defaults.Home,
	}

	
	if user.Sgids == nil {
		user.Sgids = []int{}
	}

	
	var userArg, groupArg string
	parseLine(userSpec, &userArg, &groupArg)

	
	
	uidArg, uidErr := strconv.Atoi(userArg)
	gidArg, gidErr := strconv.Atoi(groupArg)

	
	users, err := ParsePasswdFilter(passwd, func(u User) bool {
		if userArg == "" {
			
			return u.Uid == user.Uid
		}

		if uidErr == nil {
			
			return uidArg == u.Uid
		}

		return u.Name == userArg
	})

	
	if err != nil && passwd != nil {
		if userArg == "" {
			userArg = strconv.Itoa(user.Uid)
		}
		return nil, fmt.Errorf("unable to find user %s: %v", userArg, err)
	}

	var matchedUserName string
	if len(users) > 0 {
		
		matchedUserName = users[0].Name
		user.Uid = users[0].Uid
		user.Gid = users[0].Gid
		user.Home = users[0].Home
	} else if userArg != "" {
		
		

		if uidErr != nil {
			
			return nil, fmt.Errorf("unable to find user %s: %v", userArg, ErrNoPasswdEntries)
		}
		user.Uid = uidArg

		
		if user.Uid < minId || user.Uid > maxId {
			return nil, ErrRange
		}

		
	}

	
	
	if groupArg != "" || matchedUserName != "" {
		groups, err := ParseGroupFilter(group, func(g Group) bool {
			
			if groupArg == "" {
				
				for _, u := range g.List {
					if u == matchedUserName {
						return true
					}
				}
				return false
			}

			if gidErr == nil {
				
				return gidArg == g.Gid
			}

			return g.Name == groupArg
		})
		if err != nil && group != nil {
			return nil, fmt.Errorf("unable to find groups for spec %v: %v", matchedUserName, err)
		}

		
		if groupArg != "" {
			if len(groups) > 0 {
				
				user.Gid = groups[0].Gid
			} else if groupArg != "" {
				
				

				if gidErr != nil {
					
					return nil, fmt.Errorf("unable to find group %s: %v", groupArg, ErrNoGroupEntries)
				}
				user.Gid = gidArg

				
				if user.Gid < minId || user.Gid > maxId {
					return nil, ErrRange
				}

				
			}
		} else if len(groups) > 0 {
			
			user.Sgids = make([]int, len(groups))
			for i, group := range groups {
				user.Sgids[i] = group.Gid
			}
		}
	}

	return user, nil
}






func GetAdditionalGroups(additionalGroups []string, group io.Reader) ([]int, error) {
	var groups = []Group{}
	if group != nil {
		var err error
		groups, err = ParseGroupFilter(group, func(g Group) bool {
			for _, ag := range additionalGroups {
				if g.Name == ag || strconv.Itoa(g.Gid) == ag {
					return true
				}
			}
			return false
		})
		if err != nil {
			return nil, fmt.Errorf("Unable to find additional groups %v: %v", additionalGroups, err)
		}
	}

	gidMap := make(map[int]struct{})
	for _, ag := range additionalGroups {
		var found bool
		for _, g := range groups {
			
			
			if g.Name == ag || strconv.Itoa(g.Gid) == ag {
				if _, ok := gidMap[g.Gid]; !ok {
					gidMap[g.Gid] = struct{}{}
					found = true
					break
				}
			}
		}
		
		
		if !found {
			gid, err := strconv.Atoi(ag)
			if err != nil {
				return nil, fmt.Errorf("Unable to find group %s", ag)
			}
			
			if gid < minId || gid > maxId {
				return nil, ErrRange
			}
			gidMap[gid] = struct{}{}
		}
	}
	gids := []int{}
	for gid := range gidMap {
		gids = append(gids, gid)
	}
	return gids, nil
}




func GetAdditionalGroupsPath(additionalGroups []string, groupPath string) ([]int, error) {
	group, err := os.Open(groupPath)
	if err == nil {
		defer group.Close()
	}
	return GetAdditionalGroups(additionalGroups, group)
}
