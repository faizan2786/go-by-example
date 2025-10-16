package url

import (
	"errors"
	"fmt"
	"strings"
)

type URL struct {
	Scheme string
	Host   string
	Path   string
}

func (u *URL) String() string {
	return fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, u.Path)
}

func Parse(rawURL string) (*URL, error) {

	scheme, leftOverStr, found := strings.Cut(rawURL, ":")
	if !found {
		return nil, errors.New("missing ':' in the provided url string")
	}

	// if there is no '//' in left over substring (i.e. an opaque URL), return just scheme
	if !strings.HasPrefix(leftOverStr, "//") {
		return &URL{Scheme: scheme}, nil
	}

	leftOverStr = leftOverStr[2:]                  // skip "//"
	host, path, _ := strings.Cut(leftOverStr, "/") // url may or may not have a sub-path

	// construct the URL obj and return the pointer
	return &URL{Scheme: scheme, Host: host, Path: path}, nil
}
