package lab

import (
	"reflect"
	"testing"
)

func TestCompareIMPBurst(t *testing.T) {
	s := Scenario{Arrivals: []Arrival{
		{0, IMP}, {0, IMP}, {0, IMP}, {0, IMP}, {0, IMP}, {0, IMP},
		{0, Host}, {0, Host}, {2, Host},
	}}
	got, err := Compare(s)
	if err != nil {
		t.Fatal(err)
	}
	wantClasses := []Summary{{Host, 3, 3, 1}, {IMP, 6, 2, 5}}
	if !reflect.DeepEqual(got.Classes, wantClasses) {
		t.Fatalf("summaries = %+v, want %+v", got.Classes, wantClasses)
	}
	wantWaits := [][2]int{{0, 0}, {0, 1}, {1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 0}, {3, 1}, {2, 0}}
	for i, want := range wantWaits {
		result := got.Results[i]
		if result.Message != i+1 || result.Arrival != s.Arrivals[i] ||
			result.SharedWait != want[0] || result.SplitWait != want[1] {
			t.Errorf("message %d = %+v, want waits %v", i+1, result, want)
		}
	}
}

func TestSingleClassBurstAndIdleGap(t *testing.T) {
	got, err := Compare(Scenario{Arrivals: []Arrival{
		{0, Host}, {0, Host}, {0, Host}, {10000, Host}, {10000, Host},
	}})
	if err != nil {
		t.Fatal(err)
	}
	wantWaits := [][2]int{{0, 0}, {0, 1}, {1, 2}, {0, 0}, {0, 1}}
	for i, want := range wantWaits {
		r := got.Results[i]
		if r.SharedWait != want[0] || r.SplitWait != want[1] {
			t.Errorf("message %d = %+v, want waits %v", i+1, r, want)
		}
	}
	if got.Classes[1] != (Summary{Class: IMP}) {
		t.Errorf("empty class has unexpected summary: %+v", got.Classes[1])
	}
}

func TestCompareRejectsInvalidInput(t *testing.T) {
	for name, arrivals := range map[string][]Arrival{
		"empty": nil, "too many": make([]Arrival, 501), "class": {{0, "data"}},
		"negative": {{-1, Host}}, "too late": {{10001, IMP}}, "out of order": {{1, IMP}, {0, Host}},
	} {
		t.Run(name, func(t *testing.T) {
			if got, err := Compare(Scenario{Arrivals: arrivals}); err == nil || len(got.Results) != 0 {
				t.Fatalf("Compare = %+v, %v, want error and no results", got, err)
			}
		})
	}
}
