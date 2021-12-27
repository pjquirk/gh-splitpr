package cmd

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh"
)

func ToNwo(r gh.Repository) string {
	if r.Host() != "" {
		return fmt.Sprintf("%s/%s/%s", r.Host(), r.Owner(), r.Name())
	} else {
		return fmt.Sprintf("%s/%s", r.Owner(), r.Name())
	}
}

func ToRepository(nwo string) (gh.Repository, error) {
	parts := strings.Split(nwo, "/")
	if len(parts) == 2 {
		return repo{host: "", owner: parts[0], name: parts[1]}, nil
	} else if len(parts) == 3 {
		return repo{host: parts[0], owner: parts[1], name: parts[2]}, nil
	}
	return nil, fmt.Errorf("Could not extract a host/owner/name from the given identifier: %s", nwo)
}

type repo struct {
	host  string
	name  string
	owner string
}

func (r repo) Host() string {
	return r.host
}

func (r repo) Name() string {
	return r.name
}

func (r repo) Owner() string {
	return r.owner
}
