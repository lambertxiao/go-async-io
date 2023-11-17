package libaio

import (
	goaio "go-aio"
	"os"
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

	_, err = ctx.WriteAt(buf, 0)
	s.Nil(err)

	err = ctx.Close()
	s.Nil(err)

	data, err := os.ReadFile(fpath)
	s.Nil(err)

	s.Equal(buf, data)
	err = os.Remove(fpath)
	s.Nil(err)
}
