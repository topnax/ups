package utils

import (
	"testing"
	"ups/sp/server/utils"
)

func TestMax(t *testing.T) {
	got := utils.Max(1, 2)
	if got != 2 {
		t.Errorf("Max(1, 2) = %d; want 2", got)
	}
}

func TestMax2(t *testing.T) {
	got := utils.Max(1, 1)
	if got != 1 {
		t.Errorf("Max(1, 1) = %d; want 1", got)
	}
}

func TestMax3(t *testing.T) {
	got := utils.Max(2, 1)
	if got != 2 {
		t.Errorf("Max(2, 1) = %d; want 2", got)
	}
}

func TestMin(t *testing.T) {
	got := utils.Min(1, 2)
	if got != 1 {
		t.Errorf("Max(1, 2) = %d; want 1", got)
	}
}

func TestMin2(t *testing.T) {
	got := utils.Max(1, 1)
	if got != 1 {
		t.Errorf("Max(1, 1) = %d; want 1", got)
	}
}

func TestMin3(t *testing.T) {
	got := utils.Min(2, 1)
	if got != 1 {
		t.Errorf("Max(2, 1) = %d; want 1", got)
	}
}
