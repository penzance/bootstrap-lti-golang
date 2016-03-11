package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mrjones/oauth"
	"github.com/penzance/bootstrap-lti-golang/lti"
)


func main() {
	var secrets = lti.HardCodedSecretGetter{
		"test": "secret",
	}
	var provider = oauth.NewProvider(secrets.SecretGetter)
	// TODO: figure out how to bundle template files with go binaries
	var pageTemplate = template.Must(template.New("ltiBootstrap").Parse(pageTemplateString))

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
		r.ParseForm()
		pageTemplate.Execute(w, r.Form)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:9999", nil))
}

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
