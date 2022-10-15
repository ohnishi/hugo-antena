package core_test

import (
	"testing"
	"time"

	"github.com/ohnishi/antena/backend/common/core"
)

func TestNewLocalDate(t *testing.T) {
	t.Run("normal usage", func(t *testing.T) {
		date := core.NewLocalDate(2020, time.March, 16)

		rfc3399 := date.Format(time.RFC3339)
		expectedWithoutLocation := "2020-03-16T00:00:00"
		if rfc3399[:19] != expectedWithoutLocation {
			t.Errorf("date.Format(RFC3339) = %q, it should start with %q", rfc3399, expectedWithoutLocation)
		}
	})

	t.Run("can use int as month", func(t *testing.T) {
		date := core.NewLocalDate(2020, 1, 2)

		rfc3399 := date.Format(time.RFC3339)
		expectedWithoutLocation := "2020-01-02T00:00:00"
		if rfc3399[:19] != expectedWithoutLocation {
			t.Errorf("date.Format(RFC3339) = %q, it should start with %q", rfc3399, expectedWithoutLocation)
		}
	})

}

func TestParseLocal(t *testing.T) {
	tt := []struct {
		format  string
		value   string
		want    time.Time
		isError bool
	}{
		{"20060102", "20180401", time.Date(2018, 4, 1, 0, 0, 0, 0, time.Local), false},
		{"abcdefhi", "20180401", time.Time{}, true},
	}
	for _, tc := range tt {
		tv, err := core.ParseLocal(tc.format, tc.value)
		if tc.isError && err == nil {
			t.Fatal("err is expected not to be nil")
		} else if !tc.isError && err != nil {
			t.Fatalf("%+v", err)
		} else {
			t.Logf("time: %v", tv)
		}
	}
}

func TestTruncateDayLocal(t *testing.T) {
	datetime := time.Date(2006, time.January, 2, 3, 4, 5, 6, time.Local)
	actual := core.TruncateDayLocal(datetime)
	expected := core.NewLocalDate(2006, time.January, 2)
	if actual != expected {
		t.Errorf("TruncateDayLocal(2006-01-02T-03:04:05) = %q, want %q", actual, expected)
	}
}
