package utils

import (
	"testing"
	"ups/sp/server/protocol/impl"
)

func TestIndexOfFirst(t *testing.T) {
	got := impl.IndexOfNth("foobar", "b", 1)
	want := 3
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestIndexOfFirst2(t *testing.T) {
	got := impl.IndexOfNth("foobar", "r", 1)
	want := 5
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestIndexOfSecond(t *testing.T) {
	got := impl.IndexOfNth("foobarbaz", "b", 2)
	want := 6
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestIndexOfSecond2(t *testing.T) {
	got := impl.IndexOfNth("foobarbabbbb", "b", 2)
	want := 6
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestIndexOfSecond3(t *testing.T) {
	got := impl.IndexOfNth("foobarbabbbb", "b", 3)
	want := 8
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestIndexOfNotFound(t *testing.T) {
	got := impl.IndexOfNth("foobarbaz", "x", 3)
	want := -1
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestIndexOfNotFound2(t *testing.T) {
	got := impl.IndexOfNth("foobarbaz", "x", 3)
	want := -1
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestIndexOfNotFound3(t *testing.T) {
	got := impl.IndexOfNth("foobarbaz", "x", 1)
	want := -1
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestIndexOfNotFound4(t *testing.T) {
	got := impl.IndexOfNth("foobarbaz", "f", 3)
	want := -1
	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}
