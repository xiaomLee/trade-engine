package router

import (
	"github.com/gin-gonic/gin"
	grpc_end "github.com/xiaomLee/grpc-end"
	grpc_middleware "github.com/xiaomLee/grpc-end/middleware"
	"github.com/xiaomLee/trade-engine/grpchandler"
	"github.com/xiaomLee/trade-engine/httphandler"
	"github.com/xiaomLee/trade-engine/middleware"
)

func NewEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// use middleware here
	engine.Use(middleware.Recover)
	engine.Use(middleware.RequestStart)
	engine.Use(middleware.RequestOut)

	// router here
	engine.Any("/", httphandler.HealthCheck)

	return engine
}

func NewGRpcEngine() *grpc_end.GRpcEngine {
	engine := grpc_end.NewGRpcEngine("MyAppName")
	engine.RegisterFunc("hello", "world", grpchandler.SayHi)

	engine.Use(grpc_middleware.Recover)
	engine.Use(grpc_middleware.Logger)

	return engine
}
