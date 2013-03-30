package image

import (
    "strings"
    "strconv"
)

type Filter func(*Manifest) bool
type filterWrapper func(string) Filter

var allFilters = map[string]filterWrapper{
    "owner": OwnerFilter,
    "state": StateFilter,
    "name": NameFilter,
    "version": VersionFilter,
    "public": PublicFilter,
    "os": OsFilter,
    "type": TypeFilter,
}

func GetFilter(name, value string) (Filter, bool) {
    fn, ok := allFilters[name]    
    if !ok {
        return nil, false
    }

    return fn(value), true
}

func MatchManifest(filters []Filter, m *Manifest) bool {
    for _, filter := range filters {
        if !filter(m) {
            return false
        }
    }
    return true
}

func OwnerFilter(owner string) Filter {
    return func(m *Manifest) bool {
        return m.Owner == owner
    }
}

func StateFilter(str string) Filter {
    var filter Filter

    if str == "all" {
        filter = func(m *Manifest) bool {
            return true
        }
    } else {
        state := ManifestState(str)

        filter = func(m *Manifest) bool {
            return m.State == state
        }
    }

    return filter
}

func NameFilter(name string) Filter {
    var filter Filter

    if strings.HasPrefix(name, "~") {
        filter = func(m *Manifest) bool {
            return strings.Contains(m.Name, name[1:])
        }
    } else {
        filter = func(m *Manifest) bool {
            return m.Name == name
        }
    }
    
    return filter
}

func VersionFilter(version string) Filter {
    return func(m *Manifest) bool {
        return m.Version == version
    }
}

func PublicFilter(str string) Filter {
    public, _ := strconv.ParseBool(str)

    return func(m *Manifest) bool {
        return m.Public == public
    }
}

func OsFilter(os string) Filter {
    return func(m *Manifest) bool {
        return m.Os == os
    }
}

func TypeFilter(t string) Filter {
    return func(m *Manifest) bool {
        return m.Type == t
    }
}
