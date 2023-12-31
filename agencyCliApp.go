package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Data struct {
	ID      int
	Name    string
	Region  string
	Phone   string
	Address string
	Worker  int
}

func main() {
	fmt.Println("===Agency Cli App===")

	_, err := os.Stat("./data/data.csv")
	if os.IsNotExist(err) {
		file, err := os.Create("./data/data.csv")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		header := []string{"ID", "Name", "Phone", "Address", "Region", "Worker"}
		err = writer.Write(header)
		if err != nil {
			fmt.Println(err)
			return
		}
		writer.Flush()
	}

	csvFile, err := os.OpenFile("./data/data.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	writer := csv.NewWriter(csvFile)

	region := flag.String("region", "no region", "region to use")
	command := flag.String("command", "no command", "the command to execute")
	flag.Parse()

	for {
		runCommand(*command, *region, reader, writer)

		fmt.Println("please enter another command and region:")
		fmt.Scan(command, region)
	}
}

func runCommand(command string, region string, reader *csv.Reader, writer *csv.Writer) {
	switch command {
	case "list":
		commandList(region, reader)
	case "get":
		commandGet(region, reader)
	case "create":
		commandCreate(region, reader, writer)
	case "edit":
		commandEdit(region, reader, writer)
	case "status":
		commandStatus(region, reader)
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command not valid!")
	}
}

func commandList(region string, reader *csv.Reader) { // list all agencies in a region
	var agencies []Data
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		if line[2] != region {
			id, err := strconv.Atoi(line[0])
			if err != nil {
				fmt.Println(err)
			}
			worker, err := strconv.Atoi(line[5])
			if err != nil {
				fmt.Println(err)
			}
			agencies = append(agencies, Data{
				ID:      id,
				Name:    line[1],
				Region:  line[2],
				Phone:   line[3],
				Address: line[4],
				Worker:  worker,
			})
		}

	}
	fmt.Println(agencies)
}

func commandGet(region string, reader *csv.Reader) { // get a specific agency
	var agency Data
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("please enter your id:")
	scanner.Scan()
	num := scanner.Text()

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		} else if line[0] == num {
			id, err := strconv.Atoi(num)
			if err != nil {
				fmt.Println(err)
			}
			worker, err := strconv.Atoi(line[5])
			if err != nil {
				fmt.Println(err)
			}
			agency = Data{
				ID:      id,
				Name:    line[1],
				Region:  line[2],
				Phone:   line[3],
				Address: line[4],
				Worker:  worker,
			}
			fmt.Println(agency)
		}
	}
}

func commandCreate(region string, reader *csv.Reader, writer *csv.Writer) { // create a new agency
	var agency Data
	var tmp int
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("please enter your name:")
	scanner.Scan()
	agency.Name = scanner.Text()

	fmt.Println("please enter your phone:")
	scanner.Scan()
	agency.Phone = scanner.Text()

	fmt.Println("please enter your address:")
	scanner.Scan()
	agency.Address = scanner.Text()

	fmt.Println("please enter count of your worker:")
	scanner.Scan()
	workerString := scanner.Text()
	worker, err := strconv.Atoi(workerString)
	if err != nil {
		fmt.Println(err)
	}
	agency.Worker = worker

	agency.Region = region

	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	tmp = len(lines)

	agency.ID = tmp
	err = writer.Write([]string{strconv.Itoa(agency.ID), agency.Name, agency.Phone, agency.Address, agency.Region, strconv.Itoa(agency.Worker)})
	if err != nil {
		fmt.Println(err)
	}
	writer.Flush()
	fmt.Println(agency)
}

func commandEdit(region string, reader *csv.Reader, writer *csv.Writer) {
	var agency Data
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter your id:")
	scanner.Scan()
	num := scanner.Text()
	id, err := strconv.Atoi(num)
	if err != nil {
		fmt.Println(err)
	}

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		} else if line[0] == num {
			worker, err := strconv.Atoi(line[5])
			if err != nil {
				fmt.Println(err)
			}
			agency = Data{
				ID:      id,
				Name:    line[1],
				Region:  line[2],
				Phone:   line[3],
				Address: line[4],
				Worker:  worker,
			}
			for {
				fmt.Printf("%+v\nWhich one do you want to edit or exit to save:\n", agency)
				writer.Flush()
				scanner.Scan()
				field := scanner.Text()

				switch field {
				case "Name":
					fmt.Println("Please enter the new name:")
					scanner.Scan()
					newName := scanner.Text()

					agency.Name = newName
				case "Region":
					fmt.Println("Please enter the new region:")
					scanner.Scan()
					newRegion := scanner.Text()

					agency.Region = newRegion
				case "Phone":
					fmt.Println("Please enter the new phone:")
					scanner.Scan()
					newPhone := scanner.Text()

					agency.Phone = newPhone
				case "Address":
					fmt.Println("Please enter the new address:")
					scanner.Scan()
					newAddress := scanner.Text()

					agency.Address = newAddress
				case "Worker":
					fmt.Println("Please enter the new worker count:")
					scanner.Scan()
					newWorkerString := scanner.Text()
					newWorker, err := strconv.Atoi(newWorkerString)
					if err != nil {
						fmt.Println(err)
					}

					agency.Worker = newWorker
				case "exit":
					err := writer.Write([]string{
						strconv.Itoa(agency.ID),
						agency.Name,
						agency.Region,
						agency.Phone,
						agency.Address,
						strconv.Itoa(agency.Worker),
					})
					if err != nil {
						fmt.Println(err)
					}
					writer.Flush()
					return
				default:
					fmt.Println("Not a valid field!!!")
				}
			}
		}
	}
}

func commandStatus(region string, reader *csv.Reader) {
	var workerTotal, agencies int
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		if line[4] == region {
			worker, err := strconv.Atoi(line[5])
			if err != nil {
				fmt.Println(err)
			}
			workerTotal += worker
			agencies++
		}
	}
	fmt.Printf("Total Agency in %s: %d\nTotal Worker in %s: %d\n", region, agencies, region, workerTotal)
}
