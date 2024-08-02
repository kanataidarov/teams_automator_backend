package types

import (
	"github.com/kanataidarov/teams_automator/config"
	"log/slog"
)

const (
	Group    = "group"
	Meeting  = "meeting"
	OneOnOne = "oneOnOne"
)

type CtxVals struct {
	Config *config.Config
	Logger *slog.Logger
}
