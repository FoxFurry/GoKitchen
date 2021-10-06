package profiler

import (
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"net/http"
	_ "net/http/pprof"
)

func StartProfiler() {
	logger.LogMessage("Starting profiler server. Access localhost:6060 for more data!")
	go func() {
		logger.LogError(http.ListenAndServe("localhost:6060", nil).Error())
	}()
}