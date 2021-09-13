package cmd

/*
Copyright © 2021 Steven Callister

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitpath <filepath>",
	Short: "Returns the URL to a particular git path",
	RunE:  GitPathCmd,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			log.Error().Err(err)
			return err
		}

		if verbose {
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().BoolP("verbose", "v", false, "Enables verbose logging")
}
