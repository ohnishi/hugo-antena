package core_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohnishi/antena/backend/common/core"
)

func ExampleNewCSVBuilder() {
	b := core.NewCSVBuilder("header1", "header2")
	b.AddRow("value1-1", "value1-2")
	b.AddRow("value2-1", "value2-2")
	textBytes, err := b.Bytes()
	fmt.Println(err)
	lines := strings.Split(string(textBytes), "\r\n")
	for _, v := range lines {
		fmt.Println(v)
	}
	// このライブラリの出力CSVは最後のレコードの末尾にもCRLFを持つので、1(header行) + 2(レコード数) + 1(最後の空行) となる
	fmt.Printf("length = %v", len(lines))
	// Output:
	// <nil>
	// header1,header2
	// value1-1,value1-2
	// value2-1,value2-2
	//
	// length = 4
}

// Headerに何も指定しない => Header行なしのCSV/TSVを構築する
func ExampleNewTSVBuilder() {
	b := core.NewTSVBuilder()
	b.AddRow("value1-1", "value1-2")
	b.AddRow("value2-1", "value2-2")
	textBytes, err := b.Bytes()
	fmt.Println(err)
	lines := strings.Split(string(textBytes), "\r\n")
	for _, v := range lines {
		fmt.Println(v)
	}
	// 改行で分割すると、要素数は3 = 2(レコード数) + 1(最後の空行) となる
	fmt.Printf("length = %v", len(lines))
	// Output:
	// <nil>
	// value1-1	value1-2
	// value2-1	value2-2
	//
	// length = 3
}

// csvの値の中にある引用符は2重化され、全体も引用符で囲まれる (See, https://www.ietf.org/rfc/rfc4180.txt)
func TestCSVBuilder_QuotedValues(t *testing.T) {
	b := core.NewCSVBuilder("\"ファイル名\"", "見出し")
	b.AddRow("test.txt", "テスト\"(ダブルクオーテーションつき)")
	text, err := b.String()
	if err != nil {
		t.Fatalf("unexpected error while creating csv: %v", err)
	}
	lines := strings.Split(text, "\r\n")
	expectHeader := "\"\"\"ファイル名\"\"\",見出し"
	if lines[0] != expectHeader {
		t.Fatalf("expect header %s, actual %s", expectHeader, lines[0])
	}
	expectRow := "test.txt,\"テスト\"\"(ダブルクオーテーションつき)\""
	if lines[1] != expectRow {
		t.Fatalf("expect first row %s, actual %s", expectRow, lines[1])
	}
}

func TestString(t *testing.T) {
	b := core.NewCSVBuilder()
	b.AddRow("value1-1", "value1-2")
	str, err := b.String()
	if err != nil {
		t.Fatalf("failed to create csv string %v", err)
	}
	bytes, err := b.Bytes()
	if err != nil {
		t.Fatalf("failed to create csv bytes %v", err)
	}
	if str != string(bytes) {
		t.Fatalf("expect %s, got %s", str, bytes)
	}
}

func TestWithCRLF(t *testing.T) {
	b := core.NewCSVBuilder("header1", "header2")
	b.AddRow("value1-1", "value1-2")
	b.AddRow("value2-1", "value2-2")
	b = b.WithCRLF(false) // CRLF = falseとなる新たなビルダーを生成する
	text, _ := b.String()
	splitWithLF := strings.Split(text, "\n")
	if len(splitWithLF) != 4 {
		t.Fatalf("should separated with LF but has length %v", len(splitWithLF))
	}
	trySplitWithCRLF := strings.Split(text, "\r\n")
	if len(trySplitWithCRLF) != 1 {
		t.Fatalf("should not separated with CRLF but has length %v", len(trySplitWithCRLF))
	}
}

// header行とレコード行の値の数が不一致でもCSV出力できる
func TestCSVBuilder_InconsistentRowLength(t *testing.T) {
	b := core.NewCSVBuilder("header1", "header2")
	b.AddRow("value1-1")
	b.AddRow("value2-1", "value2-2", "value2-3")

	text, err := b.String()
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(text, "\r\n")
	if expect := "value2-1,value2-2,value2-3"; lines[2] != expect {
		t.Fatalf("expect line2 = %s, actual %s", expect, lines[2])
	}
}
