package vmdetect

import "testing"

func TestCommonCheck(t *testing.T) {
	inVM, msg := CommonChecks()
	if inVM && msg != "nothing" {
		t.Errorf("inside vm but got %s, expect else", msg)
	}
}
