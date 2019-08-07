package godataagent

import (
	"github.com/efjoubert/goforit/gonet"
)

//DataAgent - single DataAgent
type DataAgent struct {
}

//DataAgentsManager Data Agent(s) Manager
type DataAgentsManager struct {
	svr *gonet.Server
}

//NewDataAgentsManager - New Data Agent(s) Manager
func NewDataAgentsManager() (dtaAgntsMngr *DataAgentsManager) {
	dtaAgntsMngr = &DataAgentsManager{}
	return
}

//Startup - start Data Agent(s) Manager
func (dtaAgntsMngr *DataAgentsManager) Startup(srvlcontextpath string, srvltpath string, dtagentaias string, port string) {
	dtaAgntsMngr.svr = gonet.NewServer(port, false, "", "", gonet.DefaultServeHTTPCall)
	dtaAgntsMngr.svr.Listen()
}

//Shutdown - shutdown Data Agent(s) Manager
func (dtaAgntsMngr *DataAgentsManager) Shutdown() {

}
