package id

import "testing"

func TestIDSnowflake(t *testing.T) {
	idG := newSnowFlakeID(0)

	for i := 0; i < 100; i++ {
		t.Log(idG.NextID())
	}

	t.Log(idG.NextIDStr())
}
