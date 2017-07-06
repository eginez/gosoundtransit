package main

import (
	"testing"
	"time"
)

func TestCalcInitTime(t *testing.T) {
	now := time.Now()
	reference := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	d := CalculateDurationUntilFirstTrigger(17, reference)
	if d.Hours() != 17 {
		t.Error(d)
	}
}

func TestCalcInitTimeClosePast(t *testing.T) {
	now := time.Now()
	reference := time.Date(now.Year(), now.Month(), now.Day(), 5, 30, 0, 0, now.Location())
	d := CalculateDurationUntilFirstTrigger(5, reference)
	if d.Hours() != 0 {
		t.Error(d)
	}
}

func TestCalcInitTimeClose(t *testing.T) {
	now := time.Now()
	reference := time.Date(now.Year(), now.Month(), now.Day(), 4, 30, 0, 0, now.Location())
	d := CalculateDurationUntilFirstTrigger(5, reference)
	if d.Hours() != 0.5 {
		t.Error(d)
	}
}

func TestCalcInitTimePast(t *testing.T) {
	now := time.Now()
	reference := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	d := CalculateDurationUntilFirstTrigger(5, reference)
	if d.Hours() != 19 {
		t.Error(d)
	}
}
