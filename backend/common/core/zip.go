package core

import (
	"archive/zip"
	"bytes"
	"io"
	"time"

	"github.com/pkg/errors"
)

// ZipBuffer は、Zipファイルを構築するための構造体。
type ZipBuffer interface {
	Reader() io.Reader
	AddFile(string, []byte) error
	Close() error
}

type zipBuffer struct {
	w         *zip.Writer
	b         *bytes.Buffer
	createdAt time.Time
	closed    bool
}

// NewZipBuffer は、新しいZipBufferのインスタンスを返す。
func NewZipBuffer() ZipBuffer {
	var b bytes.Buffer
	return &zipBuffer{
		w:         zip.NewWriter(&b),
		b:         &b,
		createdAt: time.Now(),
	}
}

// AddFile は、指定したファイル名と内容を持つファイルをzipに追加する。
func (z *zipBuffer) AddFile(name string, content []byte) error {
	if z.closed {
		return errors.New("cannot add file to closed ZipBuffer")
	}
	f, err := z.w.CreateHeader(&zip.FileHeader{
		Name:     name,
		Modified: z.createdAt,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create zip content (filename: %v)", name)
	}
	_, err = f.Write(content)
	if err != nil {
		return errors.Wrapf(err, "failed to write zip archive (filename: %v)", name)
	}
	return nil
}

// Close は、ZipBufferへの書き込みを完了する。
func (z *zipBuffer) Close() error {
	if z.closed {
		return errors.New("ZipBuffer is already closed")
	}
	if err := z.w.Close(); err != nil {
		return errors.Wrap(err, "failed to close ZipBuffer")
	}
	z.closed = true
	return nil
}

// Reader は、ZipBufferの内容をReaderとして返す。
// 呼び出すたびに異なるReaderを生成する。
func (z *zipBuffer) Reader() io.Reader {
	if !z.closed {
		z.Close()
	}
	copiedBytes := make([]byte, len(z.b.Bytes()))
	copy(copiedBytes, z.b.Bytes())
	return bytes.NewReader(copiedBytes)
}
