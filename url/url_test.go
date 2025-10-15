package url

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestURLString(t *testing.T) {

	url := URL{Scheme: "http", Host: "www.dummyurl.com", Path: "mypage"}
	got := url.String()
	want := "http://www.dummyurl.com/mypage"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestParse(t *testing.T) {
	rawUrl := "https://myurl.com/myblog"
	got, err := Parse(rawUrl)
	if err != nil {
		t.Fatalf("Parse(%q) err = %q, want %v", rawUrl, err, nil)
	}

	want := &URL{Scheme: "https", Host: "myurl.com", Path: "myblog"}
	diff := cmp.Diff(want, got)
	if diff != "" {
		t.Errorf("Parse(%q) output mismatch (-want +got):\n%s", rawUrl, diff)
	}
}
