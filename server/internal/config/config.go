package config

import (
	"FaucetServer/internal/encryption"
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Config = make(map[string]string)
var mandatory = []string{"listen", "lootdir", "password"}

func Init() error {
	confHandle, err := os.Open("global.conf")
	if err != nil {
		return err
	}
	defer confHandle.Close()

	scanner := bufio.NewScanner(confHandle)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			sections := strings.Split(line, "=")
			key := sections[0]
			value := strings.Join(sections[1:], "=")
			Config[key] = value
		}
	}

	for _, option := range mandatory {
		if _, exists := Config[option]; !exists {
			return fmt.Errorf("option \"%s\" not specified", option)
		}
	}

	encryption.Key = Config["password"]

	os.Mkdir(Config["lootdir"], 0666)

	return nil
}
