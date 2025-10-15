package url

import (
	"testing"
)

func TestParse(t *testing.T) {
	rawUrl := "https://github.com/faizan2786"
	got, err := Parse(rawUrl)
	if err != nil {
		t.Fatalf("Parse(%q) err = %q, want %v", rawUrl, err, nil)
	}

	want := &URL{Scheme: "https", Host: "github.com", Path: "faizan2786"}
	if *got != *want {
		t.Errorf("Parse(%q)\ngot %#v\nwant %#v", rawUrl, got, want)
	}
}

func TestURLString(t *testing.T) {

	url := URL{Scheme: "http", Host: "www.dummyurl.com", Path: "mypage"}
	got := url.String()
	want := "http://www.dummyurl.com/mypage"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
