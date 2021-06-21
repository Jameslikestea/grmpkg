package main

import (
	"grmpkg/server"
	"os"

	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(10)
	}

	data, _ := os.ReadFile("./files/go.mod")

	path := modfile.ModulePath(data)
	logger.Sugar().Infow("found go module", "path", path)
	file, err := modfile.ParseLax("", data, func(path, version string) (string, error) {
		logger.Sugar().Infow("fixing version", "path", path, "version", version)
		return version, nil
	})
	if err != nil {
		logger.Sugar().Errorw("cannot parse go module", "error", err)
	}

	logger.Sugar().Infow("successfully parsed module", "module", file.Module.Mod.Path, "file", file)

	file.Module.Mod.Path = "somethingelse"
	file.Module.Syntax.Token[1] = "somethingelse"

	output, err := file.Format()
	if err != nil {
		logger.Sugar().Errorw("cannot output newfile", "error", err)
		return
	}

	logger.Sugar().Infow("successfully updated path", "newfile", string(output))

	s := server.New(logger.Sugar().Named("server"))
	s.Start()
}
