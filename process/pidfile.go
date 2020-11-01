// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const ConfigDirEnvKey = "GOSIVY_CONFIG_DIR"

// PIDFile gives back the path to pid file which the process port is written.
// Pid file is created when the agent is launched.
func PIDFile(pid int) (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, strconv.Itoa(pid)), nil
}

func ConfigDir() (string, error) {
	if configDir := os.Getenv(ConfigDirEnvKey); configDir != "" {
		return configDir, nil
	}

	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "gosivy"), nil
	}

	if xdgConfigDir := os.Getenv("XDG_CONFIG_HOME"); xdgConfigDir != "" {
		return filepath.Join(xdgConfigDir, "gosivy"), nil
	}

	homeDir := guessUnixHomeDir()
	if homeDir == "" {
		return "", fmt.Errorf("unable to get current user home directory: os/user lookup failed; $HOME is empty")
	}
	return filepath.Join(homeDir, ".config", "gosivy"), nil
}

func GetPort(pid int) (string, error) {
	portfile, err := PIDFile(pid)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(portfile)
	if err != nil {
		return "", err
	}
	port := strings.TrimSpace(string(b))
	return port, nil
}

func guessUnixHomeDir() string {
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}
