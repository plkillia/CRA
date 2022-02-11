package main

import (
	"encoding/json"
	"os"

	"github.com/araddon/dateparse"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func setConfigurationValues(c *cli.Context) error {
	dateStr := c.String(ArgFromDate)

	if dateStr != "" {
		dateVal, err := dateparse.ParseAny(dateStr)
		if err != nil {
			return errors.WithMessage(err, "dateparse.ParseAny")
		}

		GlobalConfig.FromDate = dateVal
		GlobalConfig.FromDateSet = true
	}

	GlobalConfig.OutputFileName = c.String(ArgOutputFile)
	GlobalConfig.RepositoryLocation = c.String(ArgRepositoryLocation)

	return nil
}

func writeOutputToFile(output WalkerOutput) error {
	outBytes, err := json.Marshal(output)
	if err != nil {
		return err
	}

	f, err := os.Create(GlobalConfig.OutputFileName)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(string(outBytes))
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}
