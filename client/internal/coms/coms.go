package coms

import (
	"FaucetClient/internal/encryption"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net"
	"sync"
)

var ending = []byte("\x00END\x00")
var lock sync.Mutex

func Send(data map[string]any, conn net.Conn) error {
	lock.Lock()
	defer lock.Unlock()

	dataM, _ := json.Marshal(data)

	encrypted, err := encryption.Encrypt(dataM)
	if err != nil {
		return err
	}

	_, err = conn.Write(append([]byte(base64.StdEncoding.EncodeToString(encrypted)), ending...))
	if err != nil {
		return err
	}

	return nil
}

func Read(conn net.Conn) (map[string]any, error) {
	lock.Lock()
	defer lock.Unlock()

	data := &bytes.Buffer{}
	for {
		smallBuffer := make([]byte, 1048576+len(ending))
		n, err := conn.Read(smallBuffer)

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		if bytes.Contains(smallBuffer, ending) {
			data.Write(smallBuffer[0:bytes.Index(smallBuffer, ending)])
			break
		}
		data.Write(smallBuffer[0:n])
	}

	dataDecoded, err := base64.StdEncoding.DecodeString(data.String())
	if err != nil {
		return nil, err
	}

	dataDecrypted, err := encryption.Decrypt(dataDecoded)
	if err != nil {
		return nil, err
	}

	payload := make(map[string]any)
	err = json.Unmarshal(dataDecrypted, &payload)
	if err != nil {
		return nil, err
	}

	if payload["result"] == "error" {
		return nil, errors.New("something went wrong client side")
	}

	return payload, nil
}
