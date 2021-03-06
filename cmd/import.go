// Copyright © 2019 Thilina Manamgoda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/config"
	"github.com/ThilinaManamgoda/password-manager/pkg/inputs"
	"github.com/ThilinaManamgoda/password-manager/pkg/passwords"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var importCmd = &cobra.Command{
	Use:   "import ",
	Short: "Import passwords",
	Long:  `Import passwords`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mPassword, err := inputs.GetFlagStringVal(cmd, inputs.FlagMasterPassword)
		if err != nil {
			return errors.Wrapf(err, inputs.ErrMsgCannotGetFlag, inputs.FlagMasterPassword)
		}
		if mPassword == "" {
			mPassword, err = inputs.PromptForMPassword()
			if err != nil {
				return errors.Wrap(err, "cannot prompt for Master password")
			}
		}
		csvFile, err := inputs.GetFlagStringVal(cmd, config.FlagCSVFile)
		if err != nil {
			return errors.Wrapf(err, inputs.ErrMsgCannotGetFlag, config.FlagCSVFile)
		}
		if csvFile == "" {
			return errors.New("must provide a medium to import")
		}
		passwordRepo, err := passwords.LoadRepo(mPassword, false)
		if err != nil {
			return errors.Wrap(err, "couldn't initialize password repository")
		}

		err = passwordRepo.Import(passwords.CSVImporterID, map[string]string{passwords.ConfKeyCSVFilePath: csvFile})
		if err != nil {
			return errors.Wrap(err, "couldn't import the CSV file")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringP(config.FlagCSVFile, "c", "", "Import passwords")
}
