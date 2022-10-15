package core

import (
	"bytes"
	"encoding/csv"

	"github.com/pkg/errors"
)

// CSVBuilder は、ヘッダーの値や各行の値からCSVを構築する。
// NewCSVBuilderとNewTSVBuilderという2つのコンストラクタがあり、それぞれCSVとTSVを構築できるが、
// 型はいずれもCSVBuilderとなる。
type CSVBuilder interface {
	AddRow(row ...string)
	Bytes() ([]byte, error)
	String() (string, error)
	WithComma(comma rune) CSVBuilder
	WithCRLF(crlf bool) CSVBuilder
}

type csvBuilder struct {
	rows  [][]string
	comma rune
	crlf  bool
}

// NewCSVBuilder は、与えられたHeaderを持つCSVBuilderを構築して返す。
// 引数が空の時は、HeaderなしのCSVBuilderを返す。
func NewCSVBuilder(header ...string) CSVBuilder {
	var rows [][]string
	if len(header) > 0 { // len = 0の場合はrowsは空にしておく
		rows = append(rows, header)
	}
	return &csvBuilder{
		rows:  rows,
		comma: ',',
		crlf:  true,
	}
}

// NewTSVBuilder は、与えられたHeaderを持つTSVBuilderを構築して返す。
// 引数が空の時は、HeaderなしのTSVBuilderを返す。
func NewTSVBuilder(header ...string) CSVBuilder {
	return NewCSVBuilder(header...).WithComma('\t')
}

// AddRow は、CSVBuilderに行を追加する。
// 引数には、行レコードに含める値を設定する。
func (b *csvBuilder) AddRow(row ...string) {
	b.rows = append(b.rows, row)
}

// Bytes は、csvBuilderからcsvテキストデータをbyte列として返す。
// encodingはUTF-8であり、BOMがつかないため、excelで開くと文字化けすることがあることに留意する。
func (b *csvBuilder) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.Comma = b.comma
	w.UseCRLF = b.crlf
	if err := w.WriteAll(b.rows); err != nil {
		return nil, errors.Wrap(err, "failed to create csv")
	}
	return buf.Bytes(), nil
}

// String は、csvBuilderからcsvテキストデータをstringとして返す。
func (b *csvBuilder) String() (string, error) {
	bytes, err := b.Bytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// WithComma は、値のセパレータを変更した新たなCSVBuilderインスタンスを返す。
func (b *csvBuilder) WithComma(comma rune) CSVBuilder {
	return &csvBuilder{
		rows:  b.rows,
		comma: comma,
		crlf:  b.crlf,
	}
}

// WithCRLF は、CRLFフラグを変更した新たなCSVBuilderインスタンスを返す。
func (b *csvBuilder) WithCRLF(crlf bool) CSVBuilder {
	return &csvBuilder{
		rows:  b.rows,
		comma: b.comma,
		crlf:  crlf,
	}
}
