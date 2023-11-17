package libaio

import (
	"errors"
	goaio "go-aio"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	ctx, err := OpenAIOCtx("/tmp/aio-test", Options{
		IODepth: 1024,
		Flag:    os.O_CREATE | os.O_RDWR | os.O_SYNC,
		Perm:    0644,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ctx == nil {
		t.Fatal(errors.New("open aio ctx failed"))
	}

	buf, err := goaio.PosixMemAlign(4096, 4096)
	if err != nil {
		t.Fatal(err)
	}

	buf[0] = 1
	buf[1] = 2
	buf[2] = 3
	buf[3] = 4
	_, err = ctx.WriteAt(buf, 0)
	if err != nil {
		t.Fatal(err)
	}

	err = ctx.Close()
	if err != nil {
		t.Fatal(err)
	}
}
