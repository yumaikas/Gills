package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"

	"gopkg.in/myesui/uuid.v1"
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

func NoteForFileNames(names []string, tag, notes string) Note {
	var buf bytes.Buffer
	if len(notes) > 0 {
		buf.WriteString(notes)
		buf.WriteString("\n")
	}

	for _, name := range names {
		buf.WriteString("![Uploaded file](/admin/upload/")
		buf.WriteString(name)
		buf.WriteString(")\n\n")
	}
	// If we have a category for this page, use that, otherwise, give it the default category
	if len(tag) > 0 {
		buf.WriteString(tag)
	} else {
		buf.WriteString("@archive @upload-")
		buf.WriteString(UploadUUID)
	}

	return Note{Content: buf.String()}
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
