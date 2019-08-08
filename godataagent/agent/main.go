package main

import (
	"os"

	godataagent ".."
	"github.com/efjoubert/goforit/goplatform"
)

func main() {
	var srvs, _ = goplatform.NewService("godataagent", "", "",
		startGoDataAgent,
		runGoDataAgent,
		stopGoDataAgent)
	if srvs != nil {
		srvs.Execute(os.Args)
	}
}

var dtaAgntsMngr *godataagent.DataAgentsManager

func startGoDataAgent(svs *goplatform.Service, args ...string) {
	if dtaAgntsMngr == nil {
		dtaAgntsMngr = godataagent.NewDataAgentsManager()
	}
}

func runGoDataAgent(svs *goplatform.Service, args ...string) {
	if dtaAgntsMngr != nil {
		dtaAgntsMngr.Startup("", "", "", ":1111")
	}
}

func stopGoDataAgent(svs *goplatform.Service, args ...string) {
	if dtaAgntsMngr != nil {
		dtaAgntsMngr.Shutdown()
		dtaAgntsMngr = nil
	}
}
