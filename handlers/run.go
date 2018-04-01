package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mattn/go-shellwords"
)

type (
	Runner     func(command []string) (result interface{}, err error)
	RunHandler struct {
		runner Runner
	}
)

func NewRunHandler(runner Runner) *RunHandler {
	return &RunHandler{runner}
}

func (h *RunHandler) Run(c echo.Context) error {
	params, err := c.FormParams()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot parse form params.")
	}
	command := params.Get("command")
	if command == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "command is required parameter.")
	}
	args, err := shellwords.Parse(command)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot parse command.")
	}
	result, err := h.runner(args)
	if err != nil {
		return err
	}
	return c.JSON(200, result)
}
