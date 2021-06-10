package procrastitracker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestConstructDestination(t *testing.T) {

	u, _ := url.Parse("http://localhost:7070/google.com")

	out, err := ConstructDestination(u)

	if err != nil {
		t.Fatal(err)
	}

	if out != "https://www.google.com" {
		t.Fail()
	}
}

func TestCanMakeRequestViaProxy(t *testing.T) {
	t.Parallel()
	go StartWebProxy()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	}))
	transport := http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			return url.Parse("http://localhost:7070")
		},
	}
	client := http.Client{
		Transport: &transport,
	}

	r, err := client.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if r.Header.Get("x-proxy") != "served by procrastiproxy" {
		t.Errorf("no x-proxy header")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "hello world" {
		t.Fatal(body)
	}
}


func TestBlockedDomainIsBlocked(t *testing.T) {
	u, err := url.Parse("https://youtube.com/my-channel")
	if err != nil {
		t.Fatal(err)
	}
	if !isBlocked(u) {
		t.Fail()
	}
}

func TestAllowedDomainIsNotBlocked(t *testing.T) {
	u, err := url.Parse("https://google.com/search")
	if err != nil {
		t.Fatal(err)
	}
	if isBlocked(u) {
		t.Fail()
	}
}