package Solution

import (
	"os"
	"testing"
	"time"
)

func TestReadFile(t *testing.T) {
	var args []string

	filename := time.Now().Format("02_01_2006_03_04_05_000")
	os.Create(filename)

	args = []string{"foo", filename}
	_, err := readFile(args)
	if err != nil {
		t.Log(err.Error())
		t.Fatal("should read created for test file")
	}

	os.WriteFile(filename, []byte("Hello world!"), 0666)
	data, _ := readFile(args)
	if data[0] != "Hello world!" {
		t.Fatal("read gone wrong")
	}

	os.Remove(filename)

	filename = time.Now().Format("02_01_2006_03_04_05_000")
	args = []string{"foo", filename}
	_, err = readFile(args)
	if err != errReadFile {
		t.Fatal("error 'cant read file' should be created")
	}

	args = []string{"foo"}
	_, err = readFile(args)
	if err != errMissingParam {
		t.Fatal("error 'missing parameter' should be created")
	}
}
