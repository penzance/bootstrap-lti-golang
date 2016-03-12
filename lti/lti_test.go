package lti

import (
	"testing"
)

func TestHardCodedSecretGetter(t *testing.T) {
	h := HardCodedSecretGetter{"consumer1": "secret1"}

	_, err := h.SecretGetter("consumer1", nil)
	if err != nil {
		t.Error("unable to get consumer for 'consumer1'", err)
	}
}
