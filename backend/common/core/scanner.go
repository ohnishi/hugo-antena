package core

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// UnmarshalFunc は1行データを読みだすための関数型(eg. json.Unmarshal)
type UnmarshalFunc func([]byte, interface{}) error

// Scanner は読み込みを担う型。
type Scanner struct {
	reader        io.ReadCloser
	scanner       *bufio.Scanner
	decompressor  io.ReadCloser
	unmarshalFunc UnmarshalFunc
}

// Next は次の行へ移動する。EOFならfalseを返す。
func (s Scanner) Next() bool {
	return s.scanner.Scan()
}

// Scan は現在行を引数vに対して読み込む。
func (s Scanner) Scan(v interface{}) (int, error) {
	b := s.scanner.Bytes()
	if len(b) == 0 {
		return 0, nil
	}

	err := s.unmarshalFunc(b, v)
	return len(b), err
}

// Close は関連リソースを一括で閉じる
func (s Scanner) Close() error {
	var errs error
	if err := s.decompressor.Close(); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := s.reader.Close(); err != nil {
		errs = multierror.Append(errs, err)
	}
	return errs
}

// NewScanner は指定パスからデータを読みだすためのScannerを生成して返す。
// unmarshalFuncはjson固定, decompressorはgzip固定。
func NewScanner(in io.ReadCloser) (*Scanner, error) {
	dcmp, err := gzip.NewReader(in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gzip reader")
	}
	return &Scanner{
		reader:        in,
		scanner:       bufio.NewScanner(dcmp),
		decompressor:  dcmp,
		unmarshalFunc: json.Unmarshal,
	}, nil
}

// NewFileScanner は指定パスからデータを読みだすためのScannerを生成して返す。
func NewFileScanner(filepath string) (*Scanner, error) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", filepath)
	}
	s, err := NewScanner(f)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create scanner: %s", filepath)
	}
	return s, nil
}
