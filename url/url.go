package url

import "fmt"

type URL struct {
	Scheme string
	Host   string
	Path   string
}

func (u *URL) String() string {
	return fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, u.Path)
}

func Parse(rawURL string) (*URL, error) {
	return &URL{Scheme: "https", Host: "github.com", Path: "faizan2786"}, nil
}
