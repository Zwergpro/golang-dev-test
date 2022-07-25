package handlers

import (
	"github.com/pkg/errors"
	"homework-1/internal/commander"
)

const (
	helpCmd   = "help"
	addCmd    = "add"
	updateCmd = "update"
	deleteCmd = "delete"
	listCmd   = "list"
)

var BadArguments = errors.New("bad arguments")

func helpCmdHandler(_ string) string {
	return `/help - list of commands
/list - list of products
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
