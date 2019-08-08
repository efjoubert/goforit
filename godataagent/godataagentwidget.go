package godataagent

import (
	"io"
	"strings"

	"github.com/efjoubert/goforit/goblet/embed"

	"github.com/efjoubert/goforit/goblet"
)

//DataAgentWidget - Data Agent Widget
type DataAgentWidget struct {
	*goblet.BaseWidget
}

//RegisterAgentSession - session
func (dtaAgntwdgt *DataAgentWidget) RegisterAgentSession() {

}

func initDataAgentWidget() {
	goblet.RegisterEmbededReaders(
		"/jquery.js", embed.JQueryJS,
		"/blockui.js", embed.BlockuiJS,
		"/webactions.js", func() io.Reader { return embed.WebactionsJS(true) },
		"/bootstrap/css/bootstrap.css", embed.BootstrapCSS,
		"/bootstrap/js/bootstrap.js")
	goblet.RegisterServletContextPath("/", "./data").RegisterServlet("/", nil)
	goblet.RegisterWidgetPath("/data", "agent.html", newDataAgentWidget)
}

const dataagentwidgethtml string = `<!doctype html>
<html>
<head>
<script type="text/javascript" src="/jquery.js"></script>
<script type="text/javascript" src="/blockui.js"></script>
<script type="text/javascript" src="/webactions.js"></script>
</head>
<body><span>DATA AGENT</span>
<body></html>`

func dataAgentWidgetHTML() io.Reader {

	return strings.NewReader(dataagentwidgethtml)
}

func newDataAgentWidget(wdgtbrkr *goblet.WidgetBroker) goblet.Widget {
	var dtaagntwdgt = &DataAgentWidget{BaseWidget: goblet.NewBaseWidget(wdgtbrkr)}
	dtaagntwdgt.DefaultWidgetMarkupHandle(dataAgentWidgetHTML)
	return (goblet.Widget)(dtaagntwdgt)
}

func init() {
	initDataAgentWidget()
}
