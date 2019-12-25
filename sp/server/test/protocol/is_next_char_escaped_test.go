package protocol

import (
	"testing"
	"ups/sp/server/protocol/impl"
)

func TestNotEscaped(t *testing.T) {
	escaped := impl.IsNextByteEscaped([]byte("not escaped"))
	if escaped {
		t.Error("This should not be escaped")
	}
}

func TestNotEscaped2(t *testing.T) {
	escaped := impl.IsNextByteEscaped([]byte("not escaped\\\\"))
	if escaped {
		t.Error("This should not be escaped")
	}
}

func TestNotEscaped3(t *testing.T) {
	escaped := impl.IsNextByteEscaped([]byte("not escaped\\\\\\\\"))
	if escaped {
		t.Error("This should not be escaped")
	}
}

func TestEscaped(t *testing.T) {
	escaped := impl.IsNextByteEscaped([]byte("escaped\\"))
	if !escaped {
		t.Errorf("'%s' should be escaped", "escaped\\")
	}
}

func TestEscaped2(t *testing.T) {
	escaped := impl.IsNextByteEscaped([]byte("escaped\\\\\\"))
	if !escaped {
		t.Errorf("'%s' should be escaped", "escaped\\\\\\")
	}
}

func TestEscaped3(t *testing.T) {
	escaped := impl.IsNextByteEscaped([]byte("escaped\\\\\\\\\\"))
	if !escaped {
		t.Errorf("'%s' should be escaped", "escaped\\\\\\\\\\")
	}
}
