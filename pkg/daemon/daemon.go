package daemon

import (
	"github.com/godpm/godpm/pkg/http"
	"github.com/godpm/godpm/pkg/log"
	"github.com/godpm/godpm/pkg/pprof"
	"github.com/godpm/godpm/pkg/process"

	"github.com/sevlyar/go-daemon"
)

// Start start daemon and process from config dir
func Start(logFile, pidFile string) {
	context := daemon.Context{PidFileName: pidFile, LogFileName: logFile}

	child, err := context.Reborn()
	if err != nil {
		log.Fatal().Fatal(err.Error())
	}

	if child != nil {
		return
	}

	defer func() {
		if err := context.Release(); err != nil {
			log.Error().Println("daemon context release failed ", err)
		}
	}()

	run()
}

func run() {
	go func() {
		log.Info().Println("Try to start process manager")
		process.InitAndStart()
	}()

	go func() {
		log.Info().Println("Try to start pprof server")
		pprof.Run()
	}()

	log.Info().Println("Try to start HTTP server")
	http.RunServer()
}
