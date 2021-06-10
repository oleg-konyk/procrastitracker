package procrastitracker_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"procrastitracker"
)

func TestCanMakeRequestViaProxy(t *testing.T) {
	t.Parallel()
	go procrastitracker.StartWebProxy()

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

func TestCannotMakeRequestBlockedSiteViaProxy(t *testing.T) {
	t.Parallel()
	go procrastitracker.StartWebProxy()

	// test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("hit a blocked URL")
	}))
	transport := http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			return url.Parse("http://localhost:7070")
		},
	}
	client := http.Client{
		Transport: &transport,
	}

	res, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("status code should be 403, got %d", res.StatusCode)
	}
}

func TestBlockedDomainIsBlocked(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://youtube.com/my-channel")
	if err != nil {
		t.Fatal(err)
	}
	if !procrastitracker.IsBlocked(u) {
		t.Fail()
	}
}

func TestAllowedDomainIsNotBlocked(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://google.com/search")
	if err != nil {
		t.Fatal(err)
	}
	if procrastitracker.IsBlocked(u) {
		t.Fail()
	}
}

func TestAddingURLToBlocker(t *testing.T) {
	t.Parallel()

	testURL := "http://google.com/"
	u, err := url.Parse(testURL)
	if err != nil {
		t.Fatal(err)
	}
	// url is not blocked
	if procrastitracker.IsBlocked(u) {
		t.Errorf("%s should not be blocked", u)
	}

	procrastitracker.Block("google.com")

	// url is blocked
	if !procrastitracker.IsBlocked(u) {
		t.Errorf("%s should be blocked", u)
	}
}
