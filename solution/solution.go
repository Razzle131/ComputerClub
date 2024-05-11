package Solution

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const timeLayout string = "15:04"

var price int = 0

func Solve() {
	inputData := readFile()

	if inputData != nil {
		processFile(inputData)
	}
}

func readFile() []string {
	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name")
		return nil
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Can't read file:", os.Args[1])
		return nil
	}

	return strings.Split(string(data), "\n")
}

func printError(errTime string, errText string) {
	fmt.Println(errTime, 13, errText)
}

func getInitData(input []string) (int, time.Time, time.Time) {
	tableCount, err := strconv.Atoi(input[0])
	if err != nil {
		panic("Error with parsing number of tables, recheck it")
	}

	workTime := strings.Split(input[1], " ")
	if workTime[0] == "24:00" {
		workTime[0] = "00:00"
	}
	if workTime[1] == "24:00" {
		workTime[1] = "00:00"
	}

	start, err := time.Parse(timeLayout, workTime[0])
	if err != nil {
		panic("Error with parsing club start time, recheck it")
	}
	end, err := time.Parse(timeLayout, workTime[1])
	if err != nil {
		panic("Error with parsing club end time, recheck it")
	}

	// if club closes before its start time it means that it is new day for closure definitely
	if end.Compare(start) <= 0 {
		end = end.AddDate(0, 0, 1)
	}

	price, err = strconv.Atoi(input[2])
	if err != nil {
		panic("Error with parsing price, recheck it")
	}

	return tableCount, start, end
}

func processFile(input []string) {
	tableCount, start, end := getInitData(input)

	var computerClub club = club{make([]table, tableCount), []string{}, []string{}}

	var fileStr []string
	var eventTime string
	var eventId string
	var eventClient string

	fmt.Println(start.Format(timeLayout))
	for i := 3; i < len(input); i++ {
		fileStr = strings.Split(input[i], " ")

		eventTime = fileStr[0]
		eventId = fileStr[1]
		eventClient = fileStr[2]
		switch eventId {
		case "1":
			fmt.Println(eventTime, eventId, eventClient)

			if slices.Contains(computerClub.cameClients, eventClient) {
				printError(eventTime, "YouShallNotPass")
				break
			}

			curEventTime, err := time.Parse(timeLayout, eventTime)
			if err != nil {
				fmt.Println("Error with parsing last event time, cant check if the club is open, moving to next event...")
				continue
			}

			if checkIfThisIsNewDay(i, input, curEventTime) {
				curEventTime = curEventTime.AddDate(0, 0, 1)
			}

			if start.Compare(curEventTime) > 0 || curEventTime.Compare(end) > 0 {
				printError(eventTime, "NotOpenYet")
				break
			}

			computerClub.cameClients = append(computerClub.cameClients, eventClient)
		case "2":
			fmt.Println(eventTime, eventId, eventClient, fileStr[3])

			if !slices.Contains(computerClub.cameClients, eventClient) {
				printError(eventTime, "ClientUnknown")
				break
			}

			eventTable, err := strconv.Atoi(fileStr[3])
			if err != nil {
				fmt.Println("Error with parsing last event table id, cant check is it busy or not, moving to next event...")
				continue
			}

			if computerClub.clubClients[eventTable-1].tableClient.clientName != "" {
				printError(eventTime, "PlaceIsBusy")
				break
			}

			curEventTime, err := time.Parse(timeLayout, eventTime)
			if err != nil {
				fmt.Println("Error with parsing last event time, cant start client`s session, moving to next event...")
				continue
			}

			// if it is new day we add 1 day to the event date
			if checkIfThisIsNewDay(i, input, curEventTime) {
				curEventTime = curEventTime.AddDate(0, 0, 1)
			}

			// skipping calculation of the session, if it is after end of the club work day
			if start.Compare(curEventTime) > 0 || curEventTime.Compare(end) > 0 {
				break
			}

			// end previous session (if exists) and start new
			if leavedClientTable := computerClub.findTableByClient(eventClient); leavedClientTable >= 0 {
				computerClub.clubClients[leavedClientTable].endClientSession(curEventTime)
			}
			computerClub.clubClients[eventTable-1].tableClient = client{eventClient, curEventTime}
		case "3":
			fmt.Println(eventTime, eventId, eventClient)

			if computerClub.findTableByClient("") >= 0 {
				printError(eventTime, "ICanWaitNoLonger!")
				break
			}

			if len(computerClub.clubQueue) > tableCount {
				fmt.Println(eventTime, 11, eventClient)
				computerClub.deleteClient(eventClient)
				break
			}

			// if client doesnt appears in common list, we will add him
			if !slices.Contains(computerClub.cameClients, eventClient) {
				computerClub.cameClients = append(computerClub.cameClients, eventClient)
			}
			computerClub.clubQueue.push(eventClient)
		case "4":
			fmt.Println(eventTime, eventId, eventClient)

			if !slices.Contains(computerClub.cameClients, eventClient) {
				printError(eventTime, "ClientUnknown")
				break
			}

			computerClub.deleteClient(eventClient)

			if leavedClientTable := computerClub.findTableByClient(eventClient); leavedClientTable >= 0 {
				curEventTime, err := time.Parse(timeLayout, eventTime)
				if err != nil {
					fmt.Println("Error with parsing last event time, cant start new client session and update table income, moving to next event...")
					continue
				}

				leavedTable := &computerClub.clubClients[leavedClientTable]
				leavedTable.endClientSession(curEventTime)
				if len(computerClub.clubQueue) > 0 {
					newClientName := computerClub.clubQueue.pop()
					fmt.Println(eventTime, 12, newClientName, leavedClientTable+1)

					leavedTable.tableClient = client{newClientName, curEventTime}
				}
			}
		}
	}

	// end all sessions and print income by tables
	slices.Sort(computerClub.cameClients)
	for _, client := range computerClub.cameClients {
		fmt.Println(end.Format(timeLayout), 11, client)
		if leavedClientTable := computerClub.findTableByClient(client); leavedClientTable >= 0 {
			leavedTable := &computerClub.clubClients[leavedClientTable]
			leavedTable.endClientSession(end)
		}
	}

	fmt.Println(end.Format(timeLayout))

	for id, table := range computerClub.clubClients {
		fmt.Println(id+1, table.tableTotalIncome, fmt.Sprintf("%02d:%02d", int(table.tableTotalTime.Hours()), int(table.tableTotalTime.Minutes())%60))
	}
}
