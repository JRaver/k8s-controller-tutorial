package testutil

import (
	"testing"
)

func TestInt32Ptr(t *testing.T) {
	i := int32(1)
	ptr := int32Ptr(i)
	if ptr == nil || *ptr != i {
		t.Errorf("int32Ptr(%d) = %v, want %v", i, ptr, i)
	}
}
