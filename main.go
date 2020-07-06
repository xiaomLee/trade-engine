package main

import (
	"context"
	"flag"
	"math"
	"net/http"
	"time"

	"github.com/judwhite/go-svc/svc"
	"github.com/xiaomLee/trade-engine/router"
	"google.golang.org/grpc"
)

type Service struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

var (
	debug = flag.Bool("debug", true, "debug will start net/http/pprof service at default 9999")
	port  = flag.String("port", "9999", "debug port default 9999")
)

func (s *Service) Init(env svc.Environment) error {

	return nil
}

func (s *Service) Start() error {
	s.httpServer = &http.Server{
		Addr:    ":1234",
		Handler: router.NewEngine(),
	}
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()
	println("http service start success,", "listening on 1234")

	var err error
	s.grpcServer, err = router.NewGRpcEngine().Run(":1235", grpc.MaxRecvMsgSize(math.MaxInt32))
	if err != nil {
		return err
	}
	println("gRpc service start success,", "listening on 1235")

	return nil
}

func (s *Service) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		println(err.Error())
	}
	println("http server stop success")

	s.grpcServer.Stop()
	println("gRpc server stop success")

	// release source here

	return nil
}

func main() {

	if *debug {
		go func() {
			http.ListenAndServe(":"+*port, nil)
			println("start debug pprof on " + *port)
			println(`
// Then use the pprof tool to look at the heap profile:
//
//	go tool pprof http://localhost:9999/debug/pprof/heap
//
// Or to look at a 30-second CPU profile:
//
//	go tool pprof http://localhost:9999/debug/pprof/profile?seconds=30
//
// Or to look at the goroutine blocking profile, after calling
// runtime.SetBlockProfileRate in your program:
//
//	go tool pprof http://localhost:9999/debug/pprof/block
//
// Or to collect a 5-second execution trace:
//
//	wget http://localhost:9999/debug/pprof/trace?seconds=5
//
// Or to look at the holders of contended mutexes, after calling
// runtime.SetMutexProfileFraction in your program:
//
//	go tool pprof http://localhost:9999/debug/pprof/mutex
//
// To view all available profiles, open http://localhost:9999/debug/pprof/
// in your browser.`)
		}()
	}

	if err := svc.Run(&Service{}); err != nil {
		println(err.Error())
	}
}
