package random

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/spf13/cobra"
)

const (
	exampleCity = `
		# Generate a random city
		!random city`
)

func newRandomCityCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "city",
		Short:   "Generate a random city",
		Example: exampleCity,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.OutOrStdout().Write([]byte(randomdata.City()))
		},
	}

	return c
}
