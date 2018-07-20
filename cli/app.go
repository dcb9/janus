package cli

import (
	"fmt"
	"os"

	"github.com/dcb9/janus/pkg/server"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("janus", "Qtum adapter to Ethereum JSON RPC")

	qtumRPC = app.Flag("qtum-rpc", "URL of qtum RPC service").Envar("QTUM_RPC").Default("").String()
	bind    = app.Flag("bind", "network interface to bind to (e.g. 0.0.0.0) ").Default("localhost").String()
	port    = app.Flag("port", "port to serve proxy").Default("23889").Int()
	devMode = app.Flag("dev", "[Insecure] Developer mode").Default("false").Bool()
)

func action(pc *kingpin.ParseContext) error {
	addr := fmt.Sprintf("%s:%d", *bind, *port)
	logger := log.NewLogfmtLogger(os.Stdout)

	if !*devMode {
		logger = level.NewFilter(logger, level.AllowWarn())
	}

	s, err := server.New(
		*qtumRPC,
		addr,
		server.SetLogger(logger),
		server.SetDebug(*devMode),
	)

	if err != nil {
		return errors.Wrap(err, "new proxy")
	}

	return s.Start()
}

func Run() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func init() {
	app.Action(action)
}
