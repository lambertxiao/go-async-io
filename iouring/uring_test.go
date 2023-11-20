package iouring

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type IOUringTestSuite struct {
	suite.Suite
}

func TestLibAIOTestSuite(t *testing.T) {
	suite.Run(t, new(IOUringTestSuite))
}

func (s *IOUringTestSuite) TestNewAIOCtx() {
	ctx, err := New(10)
	s.Nil(err)
	s.NotNil(ctx)
}
