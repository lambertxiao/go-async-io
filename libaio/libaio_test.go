package libaio

import (
	goaio "go-aio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestLibAIOTestSuite(t *testing.T) {
	suite.Run(t, new(LibAIOTestSuite))
}

type LibAIOTestSuite struct {
	suite.Suite
}

const (
	SIZE_4K = 4096
)

func (s *LibAIOTestSuite) openAIOCtx(fpath string) goaio.IOCtx {
	ctx, err := OpenAIOCtx(fpath, Options{
		IODepth: 1024,
		Flag:    os.O_CREATE | os.O_RDWR | os.O_SYNC,
		Perm:    0644,
		Timeout: 10,
	})
	s.Nil(err)
	s.NotNil(ctx)
	return ctx
}

func (s *LibAIOTestSuite) TestWriteAt() {
	fpath := "/tmp/aio-test"
	ctx := s.openAIOCtx(fpath)
	buf, err := goaio.PosixMemAlign(SIZE_4K, SIZE_4K)
	s.Nil(err)

	copy(buf, []byte("hello, aio"))

	n, err := ctx.WriteAt(buf, 0)
	s.Nil(err)
	s.Equal(SIZE_4K, n)

	err = ctx.Close()
	s.Nil(err)

	data, err := os.ReadFile(fpath)
	s.Nil(err)

	s.Equal(buf, data)
	err = os.Remove(fpath)
	s.Nil(err)
}

func (s *LibAIOTestSuite) TestWrite() {
	fpath := "/tmp/aio-test"
	ctx := s.openAIOCtx(fpath)
	buf, err := goaio.PosixMemAlign(SIZE_4K, SIZE_4K)
	s.Nil(err)

	copy(buf, []byte("hello, aio"))

	n, err := ctx.Write(buf)
	s.Nil(err)
	s.Equal(SIZE_4K, n)

	n, err = ctx.Write(buf)
	s.Nil(err)
	s.Equal(SIZE_4K, n)

	err = ctx.Close()
	s.Nil(err)

	data, err := os.ReadFile(fpath)
	s.Nil(err)

	s.Equal(append(buf, buf...), data)

	err = os.Remove(fpath)
	s.Nil(err)
}

func (s *LibAIOTestSuite) TestReadAt() {
	fpath := "/tmp/aio-test"
	fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0600)
	s.Nil(err)

	data := []byte("hello, aio")
	n, err := fd.Write(data)
	s.Nil(err)
	s.Equal(len(data), n)

	fd.Close()

	ctx := s.openAIOCtx(fpath)
	buf, err := goaio.PosixMemAlign(SIZE_4K, SIZE_4K)
	s.Nil(err)

	n, err = ctx.ReadAt(buf, 0)
	s.Nil(err)
	s.Equal(len(data), n)

	s.Equal(data, buf[:n])

	err = ctx.Close()
	s.Nil(err)

	err = os.Remove(fpath)
	s.Nil(err)
}

func (s *LibAIOTestSuite) TestRead() {
	fpath := "/tmp/aio-test"
	fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0600)
	s.Nil(err)

	data := []byte(strings.Repeat("a", 8192))
	n, err := fd.Write(data)
	s.Nil(err)
	s.Equal(len(data), n)

	fd.Close()

	ctx := s.openAIOCtx(fpath)
	buf, err := goaio.PosixMemAlign(SIZE_4K, SIZE_4K)
	s.Nil(err)

	n, err = ctx.Read(buf)
	s.Nil(err)
	s.Equal(data[:4096], buf[:n])

	n, err = ctx.Read(buf)
	s.Nil(err)
	s.Equal(data[4096:], buf[:n])

	err = ctx.Close()
	s.Nil(err)

	err = os.Remove(fpath)
	s.Nil(err)
}
