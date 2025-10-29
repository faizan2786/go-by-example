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

	const schemaSeparator string = "://"
	const pathSeparator string = "/"

	// use a string builder and pre-allocate required bytes to avoid extra mem. allocations caused by string concats
	sb := strings.Builder{}

	urlLen := len(u.Scheme) + len(schemaSeparator) + len(u.Host) + len(pathSeparator) + len(u.Path)
	sb.Grow(urlLen)

	// check scheme is not empty
	if u.Scheme != "" {
		sb.WriteString(u.Scheme)
		sb.WriteString(schemaSeparator)
	}
	if u.Host != "" {
		sb.WriteString(u.Host)
	}
	if u.Path != "" {
		sb.WriteString(pathSeparator + u.Path)
	}
	return sb.String()
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
