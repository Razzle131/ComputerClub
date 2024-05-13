package Solution

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const timeLayout string = "15:04"

func Solve() {
	inputData, err := readFile(os.Args)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	processFile(inputData)
}

var errMissingParam error = errors.New("missing parameter, provide file name")
var errReadFile error = errors.New("cant read given file")

func readFile(args []string) ([]string, error) {
	if len(args) < 2 {
		return nil, errMissingParam
	}

	data, err := os.ReadFile(args[1])
	if err != nil {
		return nil, errReadFile
	}

	return strings.Split(string(data), "\n"), nil
}

func printError(errTime string, errText string) {
	fmt.Println(errTime, 13, errText)
}

func getInitData(initData []string) (tableNum int, startTime time.Time, endTime time.Time, priceForHour int, err error) {
	if len(initData) < 3 {
		return 0, time.Time{}, time.Time{}, 0, errors.New("not enought initial data in file")
	}

	tableCount, err := strconv.Atoi(initData[0])
	if err != nil {
		return 0, time.Time{}, time.Time{}, 0, errors.New("error with parsing number of tables (first line in file), recheck it")
	}
	if tableCount <= 0 {
		return 0, time.Time{}, time.Time{}, 0, errors.New("number of tables must be positive (first line in file)")
	}

	workTime := strings.Split(initData[1], " ")
	if len(workTime) != 2 {
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

	price, err := strconv.Atoi(initData[2])
	if err != nil {
		return 0, time.Time{}, time.Time{}, 0, errors.New("error with parsing price (third line in file), recheck it")
	}
	if price < 0 {
		return 0, time.Time{}, time.Time{}, 0, errors.New("price must be non-negative (third line in file)")
	}

	return tableCount, start, end, price, nil
}

func processFile(input []string) {
	tableCount, start, end, price, err := getInitData(input[:3])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var computerClub club = club{make([]table, tableCount), []string{}, []string{}}

	var fileStr []string
	var curEventTime time.Time
	var prvEventTime time.Time
	var eventTimeStr string
	var eventId string
	var eventClient string

	fmt.Println(strings.Split(input[1], " ")[0])
	for i := 3; i < len(input); i++ {
		fileStr = strings.Split(input[i], " ")
		if len(fileStr) < 3 {
			fmt.Printf("bad data on file line %v\n", i+1)
			continue
		}

		curEventTime, err = time.Parse(timeLayout, fileStr[0])
		if err != nil {
			for j := 0; j < len(fileStr); j++ {
				fmt.Printf("%v ", fileStr[j])
			}
			fmt.Printf("\nError with parsing time on file line %v, moving to next event...\n", i+1)
			continue
		}

		eventTimeStr = fileStr[0]
		eventId = fileStr[1]
		eventClient = fileStr[2]

		// if this is not the first event, compares previous event time to current event time and adds shift in days if it is needed
		shiftDate(prvEventTime, &curEventTime)

		// if club should be closed before current event, we close it now
		if curEventTime.Compare(end) > 0 {
			computerClub.close(strings.Split(input[1], " ")[1], end, price)
		}

		switch eventId {
		case "1":
			fmt.Println(eventTimeStr, eventId, eventClient)

			if slices.Contains(computerClub.cameClients, eventClient) {
				printError(eventTimeStr, "YouShallNotPass")
				break
			}

			if start.Compare(curEventTime) > 0 || curEventTime.Compare(end) > 0 {
				printError(eventTimeStr, "NotOpenYet")
				break
			}

			computerClub.cameClients = append(computerClub.cameClients, eventClient)
		case "2":
			fmt.Println(eventTimeStr, eventId, eventClient, fileStr[3])

			if !slices.Contains(computerClub.cameClients, eventClient) {
				printError(eventTimeStr, "ClientUnknown")
				break
			}

			eventTable, err := strconv.Atoi(fileStr[3])
			if err != nil {
				fmt.Println("Error with parsing last event table id, cant check is it busy or not, moving to next event...")
				continue
			}

			if computerClub.clubClients[eventTable-1].tableClient.clientName != "" {
				printError(eventTimeStr, "PlaceIsBusy")
				break
			}

			// end previous session (if exists) and start new
			if leavedClientTable := computerClub.findTableByClient(eventClient); leavedClientTable >= 0 {
				computerClub.clubClients[leavedClientTable].endClientSession(curEventTime, price)
			}
			computerClub.clubClients[eventTable-1].tableClient = client{eventClient, curEventTime}
		case "3":
			fmt.Println(eventTimeStr, eventId, eventClient)

			if computerClub.findTableByClient("") >= 0 {
				printError(eventTimeStr, "ICanWaitNoLonger!")
				break
			}

			if len(computerClub.clubQueue) > tableCount {
				fmt.Println(eventTimeStr, 11, eventClient)
				computerClub.deleteClient(eventClient)
				break
			}

			// if client doesnt appears in common list, we will add him
			if !slices.Contains(computerClub.cameClients, eventClient) {
				computerClub.cameClients = append(computerClub.cameClients, eventClient)
			}
			computerClub.clubQueue.push(eventClient)
		case "4":
			fmt.Println(eventTimeStr, eventId, eventClient)

			if !slices.Contains(computerClub.cameClients, eventClient) {
				printError(eventTimeStr, "ClientUnknown")
				break
			}

			computerClub.deleteClient(eventClient)

			if leavedClientTable := computerClub.findTableByClient(eventClient); leavedClientTable >= 0 {
				leavedTable := &computerClub.clubClients[leavedClientTable]
				leavedTable.endClientSession(curEventTime, price)
				if len(computerClub.clubQueue) > 0 {
					newClientName := computerClub.clubQueue.pop()
					fmt.Println(eventTimeStr, 12, newClientName, leavedClientTable+1)

					leavedTable.tableClient = client{newClientName, curEventTime}
				}
			}
		}

		prvEventTime = curEventTime
	}

	// end all sessions and print income by tables
	computerClub.close(strings.Split(input[1], " ")[1], end, price)

	fmt.Println(strings.Split(input[1], " ")[1])

	// print tables income
	for id, table := range computerClub.clubClients {
		fmt.Println(id+1, table.tableTotalIncome, fmt.Sprintf("%02d:%02d", int(table.tableTotalTime.Hours()), int(table.tableTotalTime.Minutes())%60))
	}
}
