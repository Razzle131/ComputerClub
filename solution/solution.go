package Solution

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Razzle131/ComputerClub/solution/client"
	"github.com/Razzle131/ComputerClub/solution/club"
	"github.com/Razzle131/ComputerClub/solution/event"
)

const timeLayout string = "15:04"

const (
	clientComeEvent       string = "1"
	clientClaimTableEvent string = "2"
	clientWaitsEvent      string = "3"
	clientGoneEvent       string = "4"
)

func Solve() {
	inputFile, err := openFile(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	processFile(inputFile)
}

var errMissingParam error = errors.New("missing parameter, provide file name")
var errOpenFile error = errors.New("cant open given file")

func openFile(args []string) (*os.File, error) {
	if len(args) < 2 {
		return nil, errMissingParam
	}

	file, err := os.Open(args[1])
	if err != nil {
		return nil, errOpenFile
	}

	return file, nil
}

// checks if previous event was in new day or its time is bigger than new event, so modifies current date by delta days
func shiftDate(prevEventDate time.Time, curEventDate *time.Time) {
	if prevEventDate != (time.Time{}) {
		if delta := prevEventDate.Day() - curEventDate.Day(); delta > 0 {
			*curEventDate = curEventDate.AddDate(0, 0, delta)
		}
		if prevEventDate.Compare(*curEventDate) > 0 {
			*curEventDate = curEventDate.AddDate(0, 0, 1)
		}
	}
}

func getInitData(scanner *bufio.Scanner) (tableNum int, startTime time.Time, endTime time.Time, priceForHour int, err error) {
	scanner.Scan()
	tableCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, time.Time{}, time.Time{}, 0, errors.New("error with parsing number of tables (first line in file), recheck it")
	}
	if tableCount <= 0 {
		return 0, time.Time{}, time.Time{}, 0, errors.New("number of tables must be positive (first line in file)")
	}

	scanner.Scan()
	workTime := strings.Split(scanner.Text(), " ")
	if len(workTime) < 2 {
		return 0, time.Time{}, time.Time{}, 0, errors.New("error with parsing time working hours, is it divided by space?")
	}

	if workTime[0] == "24:00" {
		workTime[0] = "00:00"
	}
	if workTime[1] == "24:00" {
		workTime[1] = "00:00"
	}

	start, err := time.Parse(timeLayout, workTime[0])
	if err != nil {
		return 0, time.Time{}, time.Time{}, 0, errors.New("error with parsing club start time (second line in file, first value), recheck it")
	}
	end, err := time.Parse(timeLayout, workTime[1])
	if err != nil {
		return 0, time.Time{}, time.Time{}, 0, errors.New("error with parsing club end time (second line in file, second value), recheck it")
	}

	// if club closes before its start time it means that it is new day for closure definitely
	if end.Compare(start) <= 0 {
		end = end.AddDate(0, 0, 1)
	}

	scanner.Scan()
	price, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, time.Time{}, time.Time{}, 0, errors.New("error with parsing price (third line in file), recheck it")
	}
	if price < 0 {
		return 0, time.Time{}, time.Time{}, 0, errors.New("price must be non-negative (third line in file)")
	}

	return tableCount, start, end, price, nil
}

func processFile(inputFile *os.File) {
	scanner := bufio.NewScanner(inputFile)

	tableCount, start, end, price, err := getInitData(scanner)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var computerClub club.Club = club.New(tableCount)

	var fileStr []string
	var curEventTime time.Time
	var prvEventTime time.Time
	var fileEvent event.Event

	fmt.Println(start.Format(timeLayout))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		fileStr = strings.Split(scanner.Text(), " ")
		if len(fileStr) < 3 {
			for j := 0; j < len(fileStr); j++ {
				fmt.Printf("%v ", fileStr[j])
			}
			fmt.Println("\nnot enought data on last event line, moving to next event...")
			continue
		}

		curEventTime, err = time.Parse(timeLayout, fileStr[0])
		if err != nil {
			for j := 0; j < len(fileStr); j++ {
				fmt.Printf("%v ", fileStr[j])
			}
			fmt.Printf("\nError with parsing time on last event line, moving to next event...\n")
			continue
		}

		fileEvent = event.Event{EventTime: fileStr[0], EventId: fileStr[1], EventClient: fileStr[2], EventTable: "", EventErrMessage: ""}

		// if this is not the first event, compares previous event time to current event time and adds shift in days if it is needed
		shiftDate(prvEventTime, &curEventTime)

		// if club should be closed before current event, we close it now
		if curEventTime.Compare(end) > 0 {
			computerClub.Close(end, price, timeLayout)
		}

		switch fileEvent.EventId {
		case clientComeEvent:
			if start.Compare(curEventTime) > 0 || curEventTime.Compare(end) > 0 {
				fileEvent.EventErrMessage = "NotOpenYet"
				fileEvent.PrintEvent()
				break
			}

			if computerClub.IsInClub(fileEvent.EventClient) {
				fileEvent.EventErrMessage = "YouShallNotPass"
				fileEvent.PrintEvent()
				break
			}

			computerClub.AddClient(fileEvent.EventClient)
		case clientClaimTableEvent:
			if !computerClub.IsInClub(fileEvent.EventClient) {
				fileEvent.EventErrMessage = "ClientUnknown"
				fileEvent.PrintEvent()
				break
			}

			eventTable, err := strconv.Atoi(fileStr[3])
			fileEvent.EventTable = fileStr[3]
			if err != nil {
				fmt.Println("Error with parsing last event table id, cant check is it busy or not, moving to next event...")
				continue
			}
			if eventTable > tableCount {
				fmt.Println("last event`s table id is bigger than table count, cant check data, moving to next event...")
				continue
			}

			if computerClub.IsTableBusy(eventTable - 1) {
				fileEvent.EventErrMessage = "PlaceIsBusy"
				fileEvent.PrintEvent()
				break
			}

			computerClub.StartNewSession(client.Client{ClientName: fileEvent.EventClient, StartedTime: curEventTime}, price, eventTable-1)
		case clientWaitsEvent:
			if computerClub.QueueIsTooBig(tableCount) {
				computerClub.DeleteClient(fileEvent.EventClient)
				fileEvent.EventId = "11"
				fileEvent.PrintEvent()
				break
			}

			if computerClub.FindTableByClient("") >= 0 {
				fileEvent.EventErrMessage = "ICanWaitNoLonger!"
				fileEvent.PrintEvent()
				break
			}

			computerClub.AddClientToQueue(fileEvent.EventClient)
		case clientGoneEvent:
			if !computerClub.IsInClub(fileEvent.EventClient) {
				fileEvent.EventErrMessage = "ClientUnknown"
				fileEvent.PrintEvent()
				break
			}

			computerClub.EndClientSession(client.Client{ClientName: fileEvent.EventClient, StartedTime: curEventTime}, price)
		}

		prvEventTime = curEventTime
	}

	// end all sessions and print income by tables
	computerClub.Close(end, price, timeLayout)

	fmt.Println(end.Format(timeLayout))

	// print tables income
	for _, result := range computerClub.GetIncome() {
		fmt.Println(result)
	}
}
