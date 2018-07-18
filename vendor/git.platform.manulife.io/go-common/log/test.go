package log

import (
	"os"

	"github.com/pkg/errors"
)

func InitTester() error {
	Init("testing")

	if err := os.Setenv("LOGGER_LEVEL", "debug"); err != nil {
		return errors.Wrap(err, "unable to set LOGGER_LEVEL for testing")
	}

	if err := os.Setenv("LOGGER_FMT", "text"); err != nil {
		return errors.Wrap(err, "unable to set LOGGER_FMT for testing")
	}

	return nil
}
