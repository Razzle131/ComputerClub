package Solution

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func TestOpenFile(t *testing.T) {
	var args []string

	filename := time.Now().Format("02_01_2006_03_04_05_000")
	os.Create(filename)

	args = []string{"foo", filename}
	file, err := openFile(args)
	if err != nil {
		os.Remove(filename)
		t.Log(err.Error())
		t.Fatal("should read created for test file")
	}

	os.WriteFile(filename, []byte("Hello file test!"), 0666)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if scanner.Text() != "Hello file test!" {
		os.Remove(filename)
		t.Fatal("file expected content doesnt match with readed from file data")
	}

	os.Remove(filename)

	filename = time.Now().Format("02_01_2006_03_04_05_000")
	args = []string{"foo", filename}
	_, err = openFile(args)
	if err != errOpenFile {
		t.Fatal("error 'cant open file' should be created")
	}

	args = []string{"foo"}
	_, err = openFile(args)
	if err != errMissingParam {
		t.Fatal("error 'missing parameter' should be created")
	}
}

func TestShiftDate(t *testing.T) {
	events := []string{"18:00", "01:00", "00:30"}
	curEventTime, _ := time.Parse(timeLayout, events[2])
	curEventTimeCopy := curEventTime
	prevEventTime, _ := time.Parse(timeLayout, events[1])

	shiftDate(prevEventTime, &curEventTime)
	if curEventTime != curEventTimeCopy.AddDate(0, 0, 1) {
		t.Fatal("no time shift occured for 1 day interval")
	}

	curEventTime, _ = time.Parse(timeLayout, events[2])
	curEventTimeCopy = curEventTime
	prevEventTime, _ = time.Parse(timeLayout, events[1])
	firstEventTime, _ := time.Parse(timeLayout, events[0])

	shiftDate(firstEventTime, &prevEventTime)
	shiftDate(prevEventTime, &curEventTime)
	if curEventTime != curEventTimeCopy.AddDate(0, 0, 2) {
		t.Fatal("no time shift occured for 2 days interval")
	}

	events = []string{"18:00", "01:00", "00:30"}
	curEventTime, _ = time.Parse(timeLayout, events[0])
	curEventTimeCopy = curEventTime
	prevEventTime, _ = time.Parse(timeLayout, events[1])

	shiftDate(prevEventTime, &curEventTime)
	if curEventTime != curEventTimeCopy {
		t.Fatal("time shift occured while it should not be")
	}
}
