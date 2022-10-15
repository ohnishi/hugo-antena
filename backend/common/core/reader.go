package core

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"
)

// PeekLine は io.Reader から1行覗き込む。
// 戻り値の io.Reader は Peek で消費する前の io.Reader のように振舞う。
func PeekLine(r io.Reader) (string, io.Reader, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	sc := bufio.NewScanner(tee)
	if !sc.Scan() {
		return "", nil, errors.New("seems to have no content")
	}
	peeked := sc.Text()
	recoveredReader := io.MultiReader(&buf, r)
	return peeked, recoveredReader, nil
}
