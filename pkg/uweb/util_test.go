package uweb

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UWebUtilHandlerSuite struct {
	suite.Suite

	ctx context.Context
	//req *http.Request
	res *httptest.ResponseRecorder
}

func (s *UWebUtilHandlerSuite) SetupTest() {
	s.ctx = context.Background()

	s.res = httptest.NewRecorder()
}

func (s *UWebUtilHandlerSuite) TestContentType() {

	// act
	ToJson(s.res, nil)

	// assert
	header := s.res.Header().Get("Content-Type")
	s.Assert().Equal("application/json", header)
}

func (s *UWebUtilHandlerSuite) TestIfHasErrorsShoudReturn400() {

	// act
	ToJson(s.res, nil, errors.New("erro"))

	// assert
	header := s.res.Header().Get("Content-Type")
	s.Assert().Equal("application/json", header)
	s.Assert().Equal(http.StatusBadRequest, s.res.Code)
}

func (s *UWebUtilHandlerSuite) TestIfNotSerializeShouldReturn500() {
	// act
	f := func() bool { return false }
	ToJson(s.res, f)

	// assert
	header := s.res.Header().Get("Content-Type")
	s.Assert().Equal("application/json", header)
	s.Assert().Equal(http.StatusInternalServerError, s.res.Code)
}

func TestUWebUtilHandlerSuite(t *testing.T) {
	suite.Run(t, new(UWebUtilHandlerSuite))
}
