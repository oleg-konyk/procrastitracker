package procrastitracker

import (
	"testing"
	"net/url"
)

func TestConstructDestination(t *testing.T) {
	u, _:= url.Parse("http://localhost:7070/google.com")

	out, err := ConstructDestination(u)

	if err != nil {
		t.Fatal(err)
	}

	if out != "https://www.google.com" {
		t.Fail()
	}
}