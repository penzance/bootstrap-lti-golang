package lti

import (
	"testing"
)

func TestHardCodedSecretGetter(t *testing.T) {
	// yes, this is basically just testing a map
	h := HardCodedSecretGetter{"consumer1": "secret1"}

	_, err := h.SecretGetter("consumer1", nil)
	if err != nil {
		t.Error("unable to get consumer for 'consumer1'", err)
	}
	// TODO: we can't access the secret within the consumer's signer, but we
	//		 can try to verify a request with it.
}
