package lti

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/mrjones/oauth"
)

const (
	launchParamsKey string = "LTI_LAUNCH_PARAMS"
	sessionKey string = "LTI_SESSION"
)

func init() {
	gob.Register(&url.Values{})
}

// extend map to function as an oauth SecretGetter
type HardCodedSecretGetter map[string]string

func (h HardCodedSecretGetter) SecretGetter(key string, header map[string]string) (*oauth.Consumer, error) {
	secret, ok := h[key]
	if !ok {
		return nil, fmt.Errorf("oauth_consumer_key %s is unknown")
	}
	c := oauth.NewConsumer(key, secret, oauth.ServiceProvider{})
	return c, nil
}

// TODO: authorization check methods (ie. IsAdministrator)
// TODO: New method to supply default provider/secretgetter

type LTISessionHandler struct {
	Next http.Handler
	Provider oauth.Provider
	Store sessions.Store
}

func (h LTISessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the session object
	session, err := h.Store.Get(r, sessionKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check for a new launch
	if isLaunch(r) {
		authorized, err := h.Provider.IsAuthorized(r)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
		if authorized == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}

		// TODO: better behavior if there's an existing lti session

		// stick it in the session
		r.ParseForm()
		session.Values[launchParamsKey] = r.Form
		if err := session.Save(r,w); err != nil {
			log.Println("unable to save lti launch params to session:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	}

	// grab the launch params from the session, store in the context
	val := session.Values[launchParamsKey]
	launchParams, ok := val.(url.Values)
	if !ok {
		log.Println("wrong type for lti params in session")
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	SetLaunchParams(r, launchParams)

	h.Next.ServeHTTP(w, r)
}

func isLaunch(r *http.Request) bool {
	// save a copy of the body so we can restore it before returning
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		// TODO: should we bubble this up somehow?
		log.Println("unable to read the request body:", err)
		return false
	}
	defer func(body []byte) {
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
	}(body)

	if r.Method == http.MethodPost {
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		if err := r.ParseForm(); err != nil {
			// TODO: should we bubble this up somehow?
			log.Println("unable to parse form from body:", err)
			return false
		}

		return r.Form.Get("lti_message_type") == "basic-lti-launch-request"
	}

	return false
}

func GetLaunchParams(r *http.Request) url.Values {
	if rv := context.Get(r, launchParamsKey); rv != nil {
		return rv.(url.Values)
	}
	return nil
}

func SetLaunchParams(r *http.Request, v url.Values) {
	context.Set(r, launchParamsKey, v)
}
