package core_test

import (
	"reflect"
	"testing"

	"github.com/ohnishi/antena/backend/common/core"
)

func TestChunkRanges(t *testing.T) {
	tests := []struct {
		name      string
		len       int
		chunkSize int
		expected  []core.ChunkRange
	}{
		{
			"length 0",
			0,
			3,
			[]core.ChunkRange{},
		},
		{
			"chunkSize 1",
			1,
			1,
			[]core.ChunkRange{
				{Low: 0, High: 1},
			},
		},
		{
			"chunkSize 1 with a longer length",
			3,
			1,
			[]core.ChunkRange{
				{Low: 0, High: 1},
				{Low: 1, High: 2},
				{Low: 2, High: 3},
			},
		},
		{
			"length shorter than chunk size",
			2,
			3,
			[]core.ChunkRange{
				{Low: 0, High: 2},
			},
		},
		{
			"length equals with chunk size",
			3,
			3,
			[]core.ChunkRange{
				{Low: 0, High: 3},
			},
		},
		{
			"length longer than chunk size",
			4,
			3,
			[]core.ChunkRange{
				{Low: 0, High: 3},
				{Low: 3, High: 4},
			},
		},
		{
			"length is a multiple of chunk size",
			6,
			3,
			[]core.ChunkRange{
				{Low: 0, High: 3},
				{Low: 3, High: 6},
			},
		},
		{
			"length much longer than chunk size",
			11,
			3,
			[]core.ChunkRange{
				{Low: 0, High: 3},
				{Low: 3, High: 6},
				{Low: 6, High: 9},
				{Low: 9, High: 11},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := core.ChunkRanges(tt.len, tt.chunkSize)
			if !reflect.DeepEqual(tt.expected, r) {
				t.Fatalf("expected: %+v, actual: %+v", tt.expected, r)
			}
		})
	}
}

func TestChunkStrings(t *testing.T) {
	type args struct {
		slice     []string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			"empty slice",
			args{
				[]string{},
				3,
			},
			[][]string{},
		},
		{
			"slice size 6",
			args{
				[]string{"a", "b", "c", "d", "e", "f"},
				3,
			},
			[][]string{{"a", "b", "c"}, {"d", "e", "f"}},
		},
		{
			"slice size 7",
			args{
				[]string{"a", "b", "c", "d", "e", "f", "g"},
				3,
			},
			[][]string{{"a", "b", "c"}, {"d", "e", "f"}, {"g"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.ChunkStrings(tt.args.slice, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChunkStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomString(t *testing.T) {
	s0 := core.RandomString(0)
	if s0 != "" {
		t.Errorf("RandomString(0) returned non empty string")
	}
	s16 := core.RandomString(16)
	if len(s16) != 16 {
		t.Errorf("RandomString(16) returned a string without the length of 16")
	}
	chars := map[rune]bool{}
	for _, c := range s16 {
		chars[c] = true
	}
	if len(chars) == 1 {
		t.Errorf("RandomString(16) returned non random string")
	}
}

func TestConcatStrings(t *testing.T) {
	type args struct {
		slices [][]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"no slice",
			args{[][]string{}},
			[]string{},
		},
		{
			"one slice",
			args{[][]string{{"1", "2"}}},
			[]string{"1", "2"},
		},
		{
			"two slices",
			args{[][]string{{"1", "2"}, {"3", "4"}}},
			[]string{"1", "2", "3", "4"},
		},
		{
			"three slices",
			args{[][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}},
			[]string{"1", "2", "3", "4", "5", "6"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.ConcatStrings(tt.args.slices...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("concatStringSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncludeString(t *testing.T) {
	type args struct {
		array []string
		str   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"include",
			args{
				[]string{"aaa", "bbb"},
				"aaa",
			},
			true,
		},
		{
			"not include",
			args{
				[]string{"aaa", "bbb"},
				"ccc",
			},
			false,
		},
		{
			"empty array",
			args{
				[]string{},
				"aaa",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.IncludeString(tt.args.array, tt.args.str); got != tt.want {
				t.Errorf("IncludeString() = %v, want %v", got, tt.want)
			}
		})
	}
}
