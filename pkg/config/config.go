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

// Package config holds the functionality for configuration.
package config

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/storage"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	// YAMLFileType is the file type YAML.
	YAMLFileType = "yaml"
	// EnvPrefix is the environment variable prefix.
	EnvPrefix = "PM"
	// ErrMsgUnableToReadConf is an error message.
	ErrMsgUnableToReadConf = "unable to load configuration file %s"
	// FilePathEnv is env that represents the main configuration path.
	FilePathEnv = "CONF_PATH"
	// FlagCSVFile is the CSV file flag.
	FlagCSVFile = "csv-file"
	// FlagHTMLFile is the HTML file flag.
	FlagHTMLFile = "html-file"
	// DefaultFilePermission represents the password db file default permission.
	DefaultFilePermission = "0640"
	// DefaultFileStorageEnabled default value for whether the File storage is enabled.
	DefaultFileStorageEnabled = true
	// DefaultGoogleDriveStorageEnabled default value for whether the Google drive storage is enabled.
	DefaultGoogleDriveStorageEnabled = !DefaultFileStorageEnabled
	// DefaultPasswordDBFile default value for password database file name.
	DefaultPasswordDBFile = "password-db"
	// DefaultSelectListSize default value for select list size.
	DefaultSelectListSize = 5
	// DefaultDirectoryName default value for password-manager directory name.
	DefaultDirectoryName = "password-manager"
	// DefaultTokenFileName default value for GDrive token file name.
	DefaultTokenFileName = "g_drive_tokenFile"
)

// Config struct represent the configuration for the tool.
type Config struct {
	EncryptorID    string  `mapstructure:"encryptorID"`
	SelectListSize int     `mapstructure:"selectListSize"`
	DirectoryPath  string  `mapstructure:"directoryPath"`
	Storage        Storage `mapstructure:"storage"`
}

// TransformedConfig is constructed from Config.
// This configuration is formed to support different kind of components that is been used.
type TransformedConfig struct {
	StorageID      string
	EncryptorID    string
	SelectListSize int
	DirectoryPath  string
	Storage        map[string]string
}

// Storage represent storage configurations.
type Storage struct {
	File        File        `mapstructure:"file"`
	GoogleDrive GoogleDrive `mapstructure:"googleDrive"`
}

// File represent file storage configurations.
type File struct {
	Enable         bool   `mapstructure:"enable"`
	PasswordDBFile string `mapstructure:"passwordDBFile"`
	Permission     string `mapstructure:"permission"`
}

// GoogleDrive represent Google Drive storage configurations.
type GoogleDrive struct {
	Enable         bool   `mapstructure:"enable"`
	Directory      string `mapstructure:"directory"`
	PasswordDBFile string `mapstructure:"passwordDBFile"`
}

// Init function configures the viper.
func Init() {
	viper.SetConfigType(YAMLFileType)
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// Configuration method loads the configuration.
func Configuration() (*TransformedConfig, error) {
	err := defaultConf()
	if err != nil {
		return nil, errors.Wrap(err, "cannot set default config")
	}

	err = loadConfigFile()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal the configuration")
	}
	parsedConfig := &TransformedConfig{
		EncryptorID:    config.EncryptorID,
		SelectListSize: config.SelectListSize,
		DirectoryPath:  config.DirectoryPath,
	}

	storageConf := make(map[string]string)
	if isGoogleDriveStorage(config) {
		storageConf[storage.ConfKeyDirectory] = config.Storage.GoogleDrive.Directory
		storageConf[storage.ConfKeyPasswordDBFile] = config.Storage.GoogleDrive.PasswordDBFile
		storageConf[storage.ConfKeyTokenFilePath] = filepath.Join(config.DirectoryPath, "/"+DefaultTokenFileName)
		parsedConfig.StorageID = storage.GoogleDriveStorageID
	} else if isFileStorage(config) {
		storageConf[storage.ConfKeyFilePath] = filepath.Join(config.DirectoryPath, config.Storage.File.PasswordDBFile)
		storageConf[storage.ConfKeyFilePermission] = config.Storage.File.Permission
		parsedConfig.StorageID = storage.FileStorageID
	}
	parsedConfig.Storage = storageConf
	return parsedConfig, nil
}

func isFileStorage(c *Config) bool {
	return c.Storage.File.Enable
}

func isGoogleDriveStorage(c *Config) bool {
	return c.Storage.GoogleDrive.Enable
}

func loadConfigFile() error {
	confFile, exists := os.LookupEnv(EnvPrefix + "_" + FilePathEnv)
	if exists {
		viper.SetConfigFile(confFile)
		if err := viper.ReadInConfig(); err != nil {
			return errors.Wrapf(err, ErrMsgUnableToReadConf, confFile)
		}
	}
	return nil
}

// Must set default values to parse configs from environment variables.
func defaultConf() error {
	home, err := homedir.Dir()
	if err != nil {
		return errors.Wrap(err, "cannot retrieve Home directory path")
	}
	directoryPath := filepath.Join(home, "/"+DefaultDirectoryName)
	viper.SetDefault("storage.file.enable", DefaultFileStorageEnabled)
	viper.SetDefault("storage.file.passwordDBFile", DefaultPasswordDBFile)
	viper.SetDefault("storage.file.permission", DefaultFilePermission)
	viper.SetDefault("storage.googleDrive.enable", DefaultGoogleDriveStorageEnabled)
	viper.SetDefault("storage.googleDrive.passwordDBFile", DefaultPasswordDBFile)
	viper.SetDefault("storage.googleDrive.directory", DefaultDirectoryName)
	viper.SetDefault("encryptorID", utils.AESEncryptID)
	viper.SetDefault("selectListSize", DefaultSelectListSize)
	viper.SetDefault("directoryPath", directoryPath)
	return nil
}
