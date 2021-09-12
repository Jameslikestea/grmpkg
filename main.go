package main

import (
	"os"

	"grmpkg.com/grmpkg/server"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(10)
	}

	s := server.New(logger.Sugar().Named("server"))
	s.Start()
	a := make(chan int, 1)
	<-a
}
