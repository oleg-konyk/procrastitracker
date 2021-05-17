package procrastitracker

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func StartWebProxy()  {
	handler := func (w http.ResponseWriter, r *http.Request) {
		destination := r.URL.Path

		gres, err := http.Get("https://www." + destination[1:])
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(gres.StatusCode)
		w.Header()

		bb, err := ioutil.ReadAll(gres.Body)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(bb)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":7070", nil))
}

func ConstructDestination(u *url.URL) (string, error) {
	return "https://www.google.com", nil
}