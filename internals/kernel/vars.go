package kernel

import (
	"os"

	"github.com/withmandala/go-log"
)

// Logger our global logger
var Logger = log.New(os.Stderr).WithColor().WithDebug()
