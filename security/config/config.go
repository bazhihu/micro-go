package config

import (
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
)

var (
	Logger    *log.Logger
	KitLogger kitlog.Logger
)

func init() {
	Logger = log.New(os.Stderr, "", log.LstdFlags)

	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	KitLogger = kitlog.With(KitLogger, "ts", kitlog.DefaultTimestampUTC)
	KitLogger = kitlog.With(KitLogger, "caller", kitlog.DefaultCaller)
}
