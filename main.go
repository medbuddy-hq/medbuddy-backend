package main

import (
	"context"
	"fmt"
	mdb "medbuddy-backend/pkg/repository/mongo"
	"medbuddy-backend/service/jobs"
	"medbuddy-backend/utility"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"medbuddy-backend/internal/config"
	"medbuddy-backend/pkg/router"

	"github.com/go-playground/validator/v10"
	// rdb "brief/pkg/repository/storage/redis"
)

func init() {
	config.Setup()
	mdb.ConnectToDB()

	// Start background cron jobs
	cJobs := jobs.NewCronJob()
	cJobs.StartJobs()
	// redis.SetupRedis() uncomment when you need redis
}

func main() {
	//Load config
	logger := utility.NewLogger()
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	//gin.SetMode(gin.ReleaseMode)
	e := router.Setup(validatorRef, logger)

	// The HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", getConfig.ServerPort),
		Handler: e,
	}

	// Server run context
	serverCtx, serverCancel := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		mdb.DisconnectDB(shutdownCtx)
		jobs.StopJobs()

		// Store counter variable in redis
		// redis.StoreCounter()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Fatal(err)
		}
		shutdownCancel()
		serverCancel()
	}()

	// Run the server
	logger.Infof("Server is now listening on port: %s\n", getConfig.ServerPort)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
