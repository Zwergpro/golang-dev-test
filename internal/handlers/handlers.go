package handlers

import (
	"fmt"
	"github.com/pkg/errors"
	"homework-1/internal/commander"
	"homework-1/internal/storage"
	"strconv"
	"strings"
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

func listCmdHandler(_ string) string {
	products := storage.List()
	res := make([]string, len(products))

	for _, p := range products {
		res = append(res, p.String())
	}

	return strings.Join(res, "\n")
}

func addCmdHandler(cmdArgs string) string {
	args := strings.Split(cmdArgs, " ")
	if len(args) != 3 {
		return errors.Wrapf(BadArguments, "Invalid arguments count: %d", len(args)).Error()
	}

	name := args[0]

	price, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return errors.Wrapf(BadArguments, "Can't parse price: %s", args[1]).Error()
	}

	quantity, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return errors.Wrapf(BadArguments, "Can't parse quantity: %s", args[2]).Error()
	}

	product, err := storage.NewProduct(name, uint(price), uint(quantity))
	if err != nil {
		return err.Error()
	}

	if err := storage.Add(product); err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Product added: %s", product.String())
}

func AddHandlers(c *commander.Commander) {
	c.RegisterHandler(helpCmd, helpCmdHandler)
	c.RegisterHandler(listCmd, listCmdHandler)
	c.RegisterHandler(addCmd, addCmdHandler)
}
