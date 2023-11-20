package goaio

import "io"

type IOCtx interface {
	io.Writer
	io.WriterAt
	io.Reader
	io.ReaderAt
	io.Closer
}
