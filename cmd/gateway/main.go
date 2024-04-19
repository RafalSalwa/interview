package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/config"
	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/internal/cqrs"
	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/internal/handler"
	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/internal/router"
	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/internal/server"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/fatih/color"
	"github.com/pkg/profile"
)

const (
	profFlag        = "--prof"
	pprofServerFlag = "--pprof-server"

	cpuProf      = "cpu"
	memProf      = "mem"
	blockingProf = "blocking"
	traceProf    = "trace"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}
	l := logger.NewConsole()

	checkFlags(l)

	service, err := cqrs.NewService(cfg.Grpc)
	if err != nil {
		return err
	}

	r := router.NewHTTPRouter(l)

	cryptonHandler := handler.NewCryptonHandler(l)
	cryptonHandler.RegisterRoutes(r, cfg.Auth)

	authHandler := handler.NewAuthHandler(service, l)
	authHandler.RegisterRoutes(r, cfg.Auth)

	userHandler := handler.NewUserHandler(service, l)
	userHandler.RegisterRoutes(r, cfg.Auth.JWTToken)

	srv := server.NewServer(cfg, r, l)

	srv.ServeHTTP()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint
	srv.Shutdown()
	return nil
}

func checkFlags(l *logger.Logger) {
	args := os.Args[1:]

	if len(args) > 0 {
		var doneDebugFlags bool
		for !doneDebugFlags && len(args) > 0 {
			switch args[0] {
			case profFlag:
				switch args[1] {
				case cpuProf:
					l.Println("cpu profiling enabled.\n")
					defer profile.Start(profile.CPUProfile, profile.NoShutdownHook).Stop()
				case memProf:
					l.Println("mem profiling enabled.\n")
					defer profile.Start(profile.MemProfile, profile.NoShutdownHook).Stop()
				case blockingProf:
					l.Println("block profiling enabled\n")
					defer profile.Start(profile.BlockProfile, profile.NoShutdownHook).Stop()
				case traceProf:
					l.Println("trace profiling enabled\n")
					defer profile.Start(profile.TraceProfile, profile.NoShutdownHook).Stop()
				default:
					panic("Unexpected prof flag: " + args[1])
				}
				args = args[2:]

			case pprofServerFlag:
				// serve the pprof endpoints setup in the init function run when "net/http/pprof" is imported
				go func() {
					cyanStar := color.CyanString("*")
					l.Print(cyanStar, "Starting pprof server on port 6060.")
					l.Print("Go to", "http://localhost:6060/debug/pprof in a browser to see supported endpoints.")

					srv := &http.Server{
						Addr:         "0.0.0.0:6060",
						Handler:      nil,
						ReadTimeout:  3 * time.Second,
						WriteTimeout: 5 * time.Second,
						IdleTimeout:  5 * time.Second,
					}
					if err := srv.ListenAndServe(); err != nil {
						l.Error().Err(err).Msg("Failed to serve Prometheus metrics:")
					}
				}()
				args = args[1:]
			default:
				doneDebugFlags = true
			}
		}
	}
}
