package main

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/mrjones/oauth"
	"github.com/penzance/bootstrap-lti-golang/lti"
)

var (
	secrets = lti.HardCodedSecretGetter{
		"test": "secret",
	}
	provider     = oauth.NewProvider(secrets.SecretGetter)
	pageTemplate = template.Must(template.New("ltiBootstrap").Parse(pageTemplateString))
	store sessions.Store
)

func main() {
	authenticationKey, err := base64.StdEncoding.DecodeString("aSQ9SeTlr2GDaGPUWG/3i26NswFOaNgu4UEOwBL1Z8M=")
	if err != nil {
		log.Fatal("unable to decode session authentication key")
	}
	encryptionKey, err := base64.StdEncoding.DecodeString("Z/mtqwiuzhP/7UsXisDH1+2RkGSS2NnfW5iBeDMwFyg=")
	if err != nil {
		log.Fatal("unable to decode session encryption key")
	}
	store = sessions.NewFilesystemStore("", authenticationKey, encryptionKey)
	router := mux.NewRouter()
	router.HandleFunc("/launch", launchHandler).Methods("POST")
	router.HandleFunc("/launch_params", launchParamsHandler).Methods("GET")

	//http.Handle("/", router)
	ltiSessionHandler := lti.LTISessionHandler{router, *provider, store}
	log.Fatal(http.ListenAndServe("0.0.0.0:9999", ltiSessionHandler))
}

func launchHandler(w http.ResponseWriter, r *http.Request) {
	// if we get here, the middleware is happy with the launch
	url := "/launch_params" // TODO: resolve it via router
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func launchParamsHandler(w http.ResponseWriter, r *http.Request) {
	ltiLaunchParams := lti.GetLaunchParams(r)
	if ltiLaunchParams == nil {
		log.Println("unable to retrieve LTI params from context")
		http.Error(w, "unable to find LTI launch params",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	pageTemplate.Execute(w, ltiLaunchParams)
}

// TODO: figure out how to bundle template files with go binaries such that
//       i can split this out without requiring you to run it from src dir
const pageTemplateString = `
<html>
  <head>
    <title>Bootstrap LTI</title>
  </head>

  <body>
    <table>
      <caption>LTI Launch Parameters</caption>
      <thead>
        <tr>
          <th>Key</th>
          <th>Values</th>
        </tr>
      </thead>
      <tbody>
        {{ range $key, $values := . }}
        <tr>
          <td>{{ $key }}</td>
          <td>{{ $values }}</td>
        </tr>
        {{ end }}
      </tbody>
    </table>
  </body>
</html>
	`
