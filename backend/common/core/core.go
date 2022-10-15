// Package core はGoのビルトイン型に対するユーティリティ関数群を提供する。
package core

import (
	"math/rand"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// ChunkRange はスライスの範囲を表す(LowはinclusiveでHighはexclusive)
type ChunkRange struct {
	Low  int
	High int
}

// ChunkRanges はスライスをchunkSize毎に処理するために区切るインデックスの範囲のリストを生成する
func ChunkRanges(len, chunkSize int) []ChunkRange {
	numChunks := (len + chunkSize - 1) / chunkSize
	ranges := make([]ChunkRange, 0, numChunks)
	for i := 0; i < numChunks; i++ {
		l := chunkSize * i
		h := chunkSize * (i + 1)
		if h > len {
			h = len
		}
		ranges = append(ranges, ChunkRange{l, h})
	}
	return ranges
}

// ChunkStrings []stringをchunkSizeごとに分ける
func ChunkStrings(slice []string, chunkSize int) [][]string {
	ranges := ChunkRanges(len(slice), chunkSize)
	chunks := make([][]string, 0, len(ranges))
	for _, r := range ranges {
		chunks = append(chunks, slice[r.Low:r.High])
	}
	return chunks
}

var randomRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandomString は|length|文字のランダムな文字列を生成する。
func RandomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, length)
	for i := range b {
		b[i] = randomRunes[r.Intn(len(randomRunes))]
	}
	return string(b)
}

// ConcatStrings はstringのsliceを連結して返す。
func ConcatStrings(slices ...[]string) []string {
	merged := []string{}
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	return merged
}

// IncludeString はarrayがsを含む場合trueを返す。
func IncludeString(array []string, str string) bool {
	for _, s := range array {
		if str == s {
			return true
		}
	}
	return false
}

// CheckDuplicateInStringSlice は、stringのsliceを走査し重複があればerrorを返す
func CheckDuplicateInStringSlice(array []string) error {
	var errs error
	set := StringSet{}
	for _, s := range array {
		if set.Include(s) {
			errs = multierror.Append(errs, errors.Errorf("Duplicate found: %s", s))
		}
		set.Add(s)
	}
	return errs
}
