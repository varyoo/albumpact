package source

import (
	"testing"
)

func TestVASplit(t *testing.T) {
	if s := VASplit("Lockhart , Fishbach", ","); s != "Fishbach & Lockhart" {
		t.Error(s)
	}
	if s := VASplit("Cléa Vincent", "lol"); s != "Cléa Vincent" {
		t.Error(s)
	}
	if s := VASplit("", ""); s != "" {
		t.Error(s)
	}
}
