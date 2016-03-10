package oauth

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//
// OAuth1 2-legged provider
// Contributed by https://github.com/jacobpgallagher
//

// Provide an buffer reader which implements the Close() interface
type oauthBufferReader struct {
	*bytes.Buffer
}

// So that it implements the io.ReadCloser interface
func (m oauthBufferReader) Close() error { return nil }

type ConsumerGetter func(key string, header map[string]string) (*Consumer, error)

// Provider provides methods for a 2-legged Oauth1 provider
type Provider struct {
	ConsumerGetter ConsumerGetter

	// For mocking
	clock clock
}

// NewProvider takes a function to get the consumer secret from a datastore.
// Returns a Provider
func NewProvider(secretGetter ConsumerGetter) *Provider {
	provider := &Provider{
		secretGetter,
		&defaultClock{},
	}
	return provider
}

// Combine a URL and Request to make the URL absolute
func makeURLAbs(url *url.URL, request *http.Request) {
	if !url.IsAbs() {
		url.Host = request.Host
		if request.TLS != nil || request.Header.Get("X-Forwarded-Proto") == "https" {
			url.Scheme = "https"
		} else {
			url.Scheme = "http"
		}
	}
}

func getOauthParamsFromAuthHeader(authHeader string) (map[string]string, error) {
	var err error

	if strings.EqualFold(OAUTH_HEADER, authHeader[0:5]) {
		return nil, fmt.Errorf("no OAuth Authorization header")
	}

	authHeader = authHeader[5:]
	params := strings.Split(authHeader, ",")
	pars := make(map[string]string)
	for _, param := range params {
		vals := strings.SplitN(param, "=", 2)
		k := strings.Trim(vals[0], " ")
		v := strings.Trim(strings.Trim(vals[1], "\""), " ")
		if strings.HasPrefix(k, "oauth") {
			pars[k], err = url.QueryUnescape(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return pars, nil
}

// IsAuthorized takes an *http.Request and returns a pointer to a string containing the consumer key,
// or nil if not authorized
func (provider *Provider) IsAuthorized(request *http.Request) (*string, error) {
	var consumerKey string
	var err error
	var oauthSignature string
	var oauthTimeNumber int
	var ok bool

	makeURLAbs(request.URL, request)

	// if the oauth params are in the Authorization header, grab them
	params := map[string]string{}
	authHeader := request.Header.Get(HTTP_AUTH_HEADER)
	if authHeader != "" {
		params, err := getOauthParamsFromAuthHeader(authHeader)
		if err != nil {
			return nil, err
		}

		oauthSignature, ok = params[SIGNATURE_PARAM]
		if !ok {
			return nil, fmt.Errorf("no oauth signature")
		}
		delete(params, SIGNATURE_PARAM)

		oauthTimeNumber, err = strconv.Atoi(params[TIMESTAMP_PARAM])
		if err != nil {
			return nil, err
		}

		consumerKey, ok = params[CONSUMER_KEY_PARAM]
		if !ok {
			return nil, fmt.Errorf("no consumer key")
		}
	}

	userParams, err := parseBody(request)
	if err != nil {
		return nil, err
	}

	// if we got the oauth params from the body, remove the signature
	if oauthSignature == "" {
		signatureIndex := -1
		for i, pair := range userParams {
			if pair.key == SIGNATURE_PARAM {
				signatureIndex = i
			}
			if pair.key == TIMESTAMP_PARAM {
				oauthTimeNumber, err = strconv.Atoi(pair.value)
				if err != nil {
					return nil, err
				}
			}
			if pair.key == CONSUMER_KEY_PARAM {
				consumerKey = pair.value
			}
		}
		if signatureIndex == -1 {
			return nil, fmt.Errorf("no oauth signature")
		}
		oauthSignature = userParams[signatureIndex].value
		userParams = append(userParams[:signatureIndex], userParams[signatureIndex+1:]...)
		if oauthTimeNumber == 0 {
			return nil, fmt.Errorf("no timestamp")
		}
		if consumerKey == "" {
			return nil, fmt.Errorf("no consumer key")
		}
	}

	// Check the timestamp
	if math.Abs(float64(int64(oauthTimeNumber)-provider.clock.Seconds())) > 5*60 {
		return nil, fmt.Errorf("too much clock skew")
	}


	consumer, err := provider.ConsumerGetter(consumerKey, params)
	if err != nil {
		return nil, err
	}

	if consumer.serviceProvider.BodyHash {
		bodyHash, err := calculateBodyHash(request, consumer.signer)
		if err != nil {
			return nil, err
		}

		sentHash, ok := params[BODY_HASH_PARAM]

		if bodyHash == "" && ok {
			return nil, fmt.Errorf("body_hash must not be set")
		} else if sentHash != bodyHash {
			return nil, fmt.Errorf("body_hash mismatch")
		}
	}

	allParams := NewOrderedParams()
	for key, value := range params {
		allParams.Add(key, value)
	}

	for i := range userParams {
		allParams.Add(userParams[i].key, userParams[i].value)
	}

	baseString := consumer.requestString(request.Method, canonicalizeUrl(request.URL), allParams)
	err = consumer.signer.Verify(baseString, oauthSignature)
	if err != nil {
		return nil, err
	}

	return &consumerKey, nil
}
