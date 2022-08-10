package handlers

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"homework-1/internal/repository"
	"strconv"
	"strings"
)

func deleteCmdHandler(repository repository.Product, cmdArgs string) string {
	args := strings.Split(cmdArgs, " ")
	if len(args) != 1 {
		return errors.Wrapf(BadArguments, "Invalid arguments count: %d. Require 1", len(args)).Error()
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrapf(BadArguments, "Can't parse id: %s", args[0]).Error()
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	if err = repository.DeleteProduct(ctx, id); err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Deleted: %d", id)
}
