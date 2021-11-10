package commands

import (
	"sync"

	"github.com/alash3al/xyr/internals/kernel"
	"github.com/urfave/cli/v2"
)

// CommandFunc a command handler factory
type CommandFunc func(*kernel.Env) *cli.Command

var (
	commands     = []CommandFunc{}
	commandsLock = &sync.RWMutex{}
)

// RegisterCommand registers the specified command via its factory
func RegisterCommand(c CommandFunc) {
	commandsLock.Lock()
	defer commandsLock.Unlock()

	commands = append(commands, c)
}

// GetRegisteredCommands return a list with the registered command factories
func GetRegisteredCommands() []CommandFunc {
	commandsLock.RLock()
	defer commandsLock.RUnlock()

	return commands
}
