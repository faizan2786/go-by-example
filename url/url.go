package url

import (
	"errors"
	"strings"
)

type URL struct {
	Scheme string
	Host   string
	Path   string
}

func (u *URL) String() string {
	str := ""
	if u == nil {
		return str
	}

	// check scheme is not empty
	if sc := u.Scheme; sc != "" {
		str += u.Scheme
		str += "://"
	}
	if h := u.Host; h != "" {
		str += u.Host
	}
	if p := u.Path; p != "" {
		str += "/" + u.Path
	}
	return str
}

func Parse(rawURL string) (*URL, error) {

	scheme, leftOverStr, found := strings.Cut(rawURL, ":")
	if !found || scheme == "" {
		return nil, errors.New("missing scheme")
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
