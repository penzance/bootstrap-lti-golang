FROM golang:1.6-alpine
RUN apk add --update git && rm -rf /var/cache/apk/*
ADD . /go/src/github.com/penzance/bootstrap-lti-golang
RUN go get -d -v \
	github.com/gorilla/context \
	github.com/gorilla/mux \
	github.com/gorilla/sessions \
	github.com/mrjones/oauth
RUN go install -v github.com/penzance/bootstrap-lti-golang/bootstrap-lti
ENTRYPOINT ["/go/bin/bootstrap-lti"]
EXPOSE 9999
