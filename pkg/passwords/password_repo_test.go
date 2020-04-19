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

package passwords

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/config"
	"gotest.tools/assert"
	"os"
	"path"
	"strings"
	"testing"
)

var repo *Repository

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	setupEnvs(wd)
	// Init config package. In a command flow, this done with the "root" command.
	config.Init()

	err = InitRepo("mPassword")
	if err != nil {
		panic(err)
	}
	repo, err = LoadRepo("mPassword")
	if err != nil {
		panic(err)
	}
	err = repo.Import(CSVImporterID, map[string]string{ConfKeyCSVFilePath: path.Join(wd, "../../test/mock-data/data.csv")})
	if err != nil {
		panic(err)
	}
}

func setupEnvs(wd string) {
	unSetEnvs()
	setEnv("PM_DIRECTORYPATH", path.Join(wd, "password-manager-tmp"))
}

func setEnv(env, val string) {
	err := os.Setenv(env, val)
	if err != nil {
		panic(err)
	}
}

func unSetEnvs() {
	for _, val := range os.Environ() {
		key := strings.SplitN(val, "=", 2)[0]
		if strings.HasPrefix(key, "PM_") {
			err := os.Unsetenv(key)
			if err != nil {
				panic(err)
			}
		}
	}
}

func TestGet(t *testing.T) {
	entry, err := repo.GetPasswordEntry("bmcandie15@devhub.com")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Binny", entry.Username)
	assert.Equal(t, "Qa88ookYyY", entry.Password)

	_, err = repo.GetPasswordEntry("invalid@id.com")
	assert.Error(t, err, ErrInvalidID("invalid@id.com").Error())
}

func TestSearchID(t *testing.T) {
	entries, err := repo.SearchEntriesByID("bluckcock")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "bluckcockro@answers.com", entries[0].ID)

	_, err = repo.SearchEntriesByID("invalid@id.com")
	assert.Error(t, err, ErrCannotFindMatchForID("invalid@id.com").Error())
}

func TestSearchLabel(t *testing.T) {
	entries, err := repo.SearchLabel("five")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 130, len(entries))
	assert.Equal(t, true, isPasswordEntriesContainsID("agwyerb@wisc.edu", entries))
}

func isPasswordEntriesContainsID(ID string, entries []Entry) bool {
	for _, entry := range entries {
		if entry.ID == ID {
			return true
		}
	}
	return false
}

func TestAdd(t *testing.T) {
	err := repo.Add("test@test.com", "test", "password", "test-desc", []string{"l1", "l2"})
	if err != nil {
		t.Error(err)
	}
	entry, ok := repo.db.Entries["test@test.com"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "test", entry.Username)
	assert.Equal(t, "password", entry.Password)
	assert.Equal(t, "test-desc", entry.Description)

	ids, err := repo.SearchLabel("l1")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test@test.com", ids[0].ID)
	ids, err = repo.SearchLabel("l2")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test@test.com", ids[0].ID)
}

func TestRemove(t *testing.T) {
	err := repo.Remove("bluckcockro@answers.com")
	if err != nil {
		t.Error(err)
	}
	_, ok := repo.db.Entries["bluckcockro@answers.com"]
	assert.Equal(t, false, ok)
}

func TestChangePasswordEntry(t *testing.T) {
	err := repo.ChangePasswordEntry("lgaggd@purevolume.com", Entry{Username: "change1", Password: "change2"})
	if err != nil {
		t.Error(err)
	}
	entry, ok := repo.db.Entries["lgaggd@purevolume.com"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "change1", entry.Username)
	assert.Equal(t, "change2", entry.Password)
}
