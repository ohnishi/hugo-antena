package core_test

import (
	"bufio"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ohnishi/antena/backend/common/core"
)

func TestPeekLine(t *testing.T) {
	content := `abc
def
ghi`
	original := strings.NewReader(content)
	peeked, newReader, err := core.PeekLine(original)

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if peeked != "abc" {
		t.Errorf("peeked = %q, want \"abc\"", peeked)
	}

	all, err := ioutil.ReadAll(newReader)

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	if string(all) != content {
		t.Errorf("content = %q, want %q", string(all), content)
	}
}

func TestPeekLine_InfiniteReader(t *testing.T) {
	original := &infiniteReader{i: 0}

	peeked, newReader, err := core.PeekLine(original)

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if peeked != "AB" {
		t.Errorf("peeked = %q, want \"AB\"", peeked)
	}

	sc := bufio.NewScanner(newReader)
	if !sc.Scan() {
		t.Fatal("sc.Scan() = false, want true")
	}
	txt := sc.Text()
	if txt != "AB" {
		t.Errorf("txt = %q, want \"AB\"", txt)
	}
	if !sc.Scan() {
		t.Fatal("sc.Scan() = false, want true")
	}
	txt = sc.Text()
	if txt != "CAB" {
		t.Errorf("txt = %q, want \"CAB\"", txt)
	}
}

type infiniteReader struct {
	i int
}

func (r *infiniteReader) Read(p []byte) (n int, err error) {
	if r.i == 0 {
		p[0] = 'A'
	} else if r.i == 1 {
		p[0] = 'B'
	} else if r.i == 2 {
		p[0] = '\n'
	} else {
		p[0] = 'C'
	}
	r.i = (r.i + 1) % 4
	return 1, nil
}
