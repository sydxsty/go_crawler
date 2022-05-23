package util

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetMediaInfo mediaName can contain sub path of main medias path
func GetMediaInfo(execFileDir, mediaPath, mediaName string) (string, error) {
	// apt-get install mediainfo to install mediainfo cli
	execFileDir, _ = filepath.Abs(execFileDir)
	mediaPath, _ = filepath.Abs(mediaPath)

	var result string
	callback := func(str string) error {
		result += strings.ReplaceAll(str, mediaPath, "{Hidden absolute path}")
		return nil
	}
	err := ExecCmd(`./`, filepath.Join(execFileDir, `mediainfo`), []string{filepath.Join(mediaPath, mediaName)}, callback)
	if err != nil {
		return "", errors.Wrap(err, "generate media info failed")
	}

	return result, nil
}

func ExecCmd(dir, cmdStr string, params []string, callback func(string) error) error {
	cmd := exec.Command(cmdStr, params...)
	log.Println("Exec params", dir, cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "cmd.StdoutPipe err")
	}
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		return err
	}
	// close pipe async
	defer func() {
		go func() {
			err := cmd.Wait()
			if err != nil {
				log.Println("pipe closed error, ", err)
			}
		}()
	}()
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = callback(line)
		if err != nil {
			return err
		}
	}
}
