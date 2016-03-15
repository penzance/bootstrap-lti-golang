FROM golang:1.6-alpine
RUN apk add --update git && rm -rf /var/cache/apk/*
ADD . /go/src/github.com/penzance/bootstrap-lti-golang
RUN go get -d -v \
	github.com/gorilla/context \
	github.com/gorilla/mux \
	github.com/gorilla/sessions \
	github.com/mrjones/oauth \
	k8s.io/kubernetes/pkg/util/sets
RUN go install -v github.com/penzance/bootstrap-lti-golang/bootstraplti
ENTRYPOINT ["/go/bin/bootstraplti"]
EXPOSE 9999
