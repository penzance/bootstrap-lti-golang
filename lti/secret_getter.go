package lti

import (
	"fmt"

	"github.com/mrjones/oauth"
)

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
