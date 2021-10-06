package profiler

import (
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"net/http"
	_ "net/http/pprof"
)

func StartProfiler() {
	logger.LogMessage("Starting kitchen profiler. Access http://localhost:6060/debug/pprof/ for more useful data!")
	go func() {
		logger.LogError(http.ListenAndServe("localhost:6060", nil).Error())
	}()
}