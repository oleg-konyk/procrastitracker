package procrastitracker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//var blockedList = []string{
//}

var blockedList []string

func Block(url string) {
	blockedList = append(blockedList, url)
}

func IsBlocked(url *url.URL) bool {
	for _, val := range blockedList {
		if strings.Contains(url.Host, val) {
			return true
		}
	}
	return false
}

//type BList struct {}

func StartWebProxy() {
	log.Printf("web proxy started")
	handler := func(w http.ResponseWriter, r *http.Request) {
		// https://google.com
		destination := r.URL
		log.Printf("proxying %v", destination)

		// do blocking here
		if IsBlocked(destination) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "%s is blocked -- reported to arundel", destination)
			return
		}

		gres, err := http.Get(destination.String())
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("x-proxy", "served by procrastiproxy")
		bb, err := ioutil.ReadAll(gres.Body)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(gres.StatusCode)
		_, err = w.Write(bb)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":7070", nil))
	// create a new blockedList, return that
}
