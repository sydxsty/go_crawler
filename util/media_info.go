package util

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	err := ExecCmd(execFileDir, filepath.Join(execFileDir, `mediainfo`), []string{filepath.Join(mediaPath, mediaName)}, callback)
	if err != nil {
		return "", errors.Wrap(err, "generate media info failed")
	}

	return result, nil
}

func GetMediaImage(execFileDir, mediaPath, mediaName string) ([]byte, error) {
	// please unzip mtn-linux.rar to ./lib, may need ffmpeg
	execFileDir, _ = filepath.Abs(execFileDir)
	mediaPath, _ = filepath.Abs(mediaPath)
	pathReg := regexp.MustCompile(`output: +(.+\.jpg)`)
	var jpgFile string
	callback := func(str string) error {
		v := pathReg.FindStringSubmatch(str)
		if v != nil {
			jpgFile = filepath.Clean(v[1])
		}
		return nil
	}
	err := ExecCmd(execFileDir, filepath.Join(execFileDir, `mtn`), []string{
		"-f", filepath.Join(execFileDir, `tahomabd.ttf`),
		"-c", "4",
		"-r", "3",
		"-P",
		"-O", execFileDir, filepath.Join(mediaPath, mediaName)},
		callback)
	if err != nil {
		return nil, errors.Wrap(err, "generate media info failed")
	}
	log.Println("media image is saved to", jpgFile)
	data, err := os.ReadFile(jpgFile)
	if err != nil {
		return nil, errors.Wrap(err, "can not read file: "+jpgFile)
	}
	err = os.Remove(jpgFile)
	if err != nil {
		log.Println("file remove failed, ", jpgFile, err)
	}
	return data, nil
}

func ExecCmd(dir, cmdStr string, params []string, callback func(string) error) error {
	cmd := exec.Command(cmdStr, params...)
	log.Println("Exec params", dir, cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "cmd.StdoutPipe err")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "cmd.StderrPipe err")
	}
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
	err = readFrom(stdout, callback)
	if err != nil {
		return err
	}
	return readFrom(stderr, callback)
}

func readFrom(out io.ReadCloser, callback func(string) error) error {
	reader := bufio.NewReader(out)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = callback(line)
		if err != nil {
			return err
		}
	}
	return nil
}
