package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mrjones/oauth"
)


type HardCodedSecretGetter map[string]string
func (h HardCodedSecretGetter) secretGetter(key string, header map[string]string) (*oauth.Consumer, error) {
	secret, ok := h[key]
	if !ok {
		return nil, fmt.Errorf("oauth_consumer_key %s is unknown")
	}

	c := oauth.NewConsumer(key, secret, oauth.ServiceProvider{})
	return c, nil
}


func main() {
	var secrets = HardCodedSecretGetter{
		"test": "secret",
	}
	var provider = oauth.NewProvider(secrets.secretGetter)

	http.HandleFunc("/launch", func (w http.ResponseWriter, r *http.Request) {
		authorized, err := provider.IsAuthorized(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if authorized == nil {
			w.WriteHeader(http.StatusUnauthorized)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("launch authorized"))
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:9999", nil))
}
