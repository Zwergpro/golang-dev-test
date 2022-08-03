package handlers

import (
	"github.com/pkg/errors"
	"homework-1/internal/commander"
	"homework-1/internal/repository"
	"time"
)

const (
	helpCmd   = "help"
	addCmd    = "add"
	updateCmd = "update"
	deleteCmd = "delete"
	listCmd   = "list"

	maxTimeout = time.Millisecond * 30
)

var BadArguments = errors.New("bad arguments")

func helpCmdHandler(_ repository.Product, _ string) string {
	return `/help - list of commands
/list [page] [size]  - list of products
/add <name> <price> <quantity> - add new product
/update <id> <name> <price> <quantity> - update product by id
/delete <id> - delete product
`
}

func AddHandlers(c *commander.Commander) {
	c.RegisterHandler(helpCmd, helpCmdHandler)
	c.RegisterHandler(listCmd, listCmdHandler)
	c.RegisterHandler(addCmd, addCmdHandler)
	c.RegisterHandler(deleteCmd, deleteCmdHandler)
	c.RegisterHandler(updateCmd, updateCmdHandler)
}
