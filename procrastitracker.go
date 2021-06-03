package procrastitracker

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func StartWebProxy() {
	log.Printf("web proxy started")
	handler := func(w http.ResponseWriter, r *http.Request) {
		destination := r.URL
		log.Printf("proxying %v", destination)
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
}

func ConstructDestination(u *url.URL) (string, error) {
	return "https://www.google.com", nil
}
