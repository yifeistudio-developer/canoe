package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/suite"
	"github.com/yifeistudio-developer/canoe/internal/application/core/domain"
)

type MiddlewareTest struct {
	suite.Suite
	server *httptest.Server
	app    *iris.Application
}

func (suite *MiddlewareTest) SetupTest() {
	app := iris.Default()
	app.Get("/error/bad-request", func(ctx iris.Context) {
		panic(domain.Fail(400, "Bad Request Msg"))
	})
	app.Get("/error/business", func(ctx iris.Context) {
		panic(domain.Fail(1024, "Business Msg"))
	})
	app.Get("/error/unexpected", func(ctx iris.Context) {
		panic("unexpected error")
	})
	app.UseGlobal(ContextErrorHandler)
	app.UseError(ErrorHandler)
	err := app.Build()
	if err != nil {
		log.Fatalf("build iris server error: %v", err)
		return
	}
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
	var result = domain.Result{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.Nil(err)
	suite.Equal(result, domain.Fail(404, "Not Found"))
}

func (suite *MiddlewareTest) Test_Middleware_ErrorHandler_Bad_Request_Error() {
	resp, err := http.Get(suite.server.URL + "/error/bad-request")
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	var result = domain.Result{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.Nil(err)
	suite.Equal(domain.Fail(400, "Bad Request Msg"), result)
}

func (suite *MiddlewareTest) Test_Middleware_ErrorHandler_Business_Error() {
	resp, err := http.Get(suite.server.URL + "/error/business")
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var result = domain.Result{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.Nil(err)
	suite.Equal(domain.Fail(1024, "Business Msg"), result)
}

func (suite *MiddlewareTest) Test_Middleware_ErrorHandler_Unexpected_Error() {
	resp, err := http.Get(suite.server.URL + "/error/unexpected")
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
	var result = domain.Result{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.Nil(err)
	suite.Equal(domain.Fail(500, "Internal Server Error"), result)
}

func TestErrorHandler(t *testing.T) {
	suite.Run(t, new(MiddlewareTest))
}
