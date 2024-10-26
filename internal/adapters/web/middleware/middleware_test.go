package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MiddlewareTest struct {
	suite.Suite
	server *httptest.Server
	app    *iris.Application
}

func (suite *MiddlewareTest) SetupTest() {
	app := iris.New()
	app.UseError(ErrorHandler)
	err := app.Build()
	if err != nil {
		log.Fatalf("build iris server error: %v", err)
		return
	}
	//app.Router
	suite.app = app
	suite.server = httptest.NewServer(app)
}

func (suite *MiddlewareTest) TearDownTest() {
	suite.server.Close()
}

func (suite *MiddlewareTest) Test_Middleware_ErrorHandler_Not_Found_Error() {
	resp, err := http.Get(suite.server.URL + "/nonexistent")
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func TestErrorHandler(t *testing.T) {
	suite.Run(t, new(MiddlewareTest))
}
