package main

import (
	"errors"
	"fmt"
	"gopkg.in/myesui/uuid.v1"
	"io"
	"os"
	"path"
)

var pathDir string

func PathForName(name string) string {
	return path.Join(pathDir, name)
}

func SaveUploadedFile(formFile io.Reader, extension string) (string, error) {
	randName := uuid.BulkV4(1)[0].String() + "." + extension
	ostream, err := os.OpenFile(path.Join(pathDir, randName), os.O_CREATE|os.O_WRONLY, 0750)
	defer ostream.Close()
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(ostream, formFile); err != nil {
		return "", err
	}
	return randName, nil
}

func NoteForFileName(name string) Note {
	content := fmt.Sprint(
		"![Uploaded file](/admin/upload/", name, ")\n\n",
		"@archive @upload-", UploadUUID)

	return Note{Content: content}
}

var ErrNotDir = errors.New("$PWD/uploads is not a directory!")

func prepUploadFolder() error {
	uuid.SwitchFormat(uuid.FormatHex)
	var wd, err = os.Getwd()
	if err != nil {
		return err
	}
	pathDir = path.Join(wd, "uploads")
	if info, err := os.Stat(path.Join(wd, "uploads")); os.IsNotExist(err) {
		err := os.Mkdir(pathDir, 0750)
		if err != nil {
			return err
		}
	} else {
		if !info.IsDir() {
			return ErrNotDir
		}
	}
	return nil
}
