package handlers

import (
	"fmt"
	"github.com/pkg/errors"
	"homework-1/internal/storage"
	"strconv"
	"strings"
)

func deleteCmdHandler(cmdArgs string) string {
	args := strings.Split(cmdArgs, " ")
	if len(args) != 1 {
		return errors.Wrapf(BadArguments, "Invalid arguments count: %d. Require 1", len(args)).Error()
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrapf(BadArguments, "Can't parse id: %s", args[0]).Error()
	}

	if err = storage.Delete(id); err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Deleted: %d", id)
}
