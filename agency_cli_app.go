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

type Service struct {
	reader      *csv.Reader
	writer      *csv.Writer
	csv         *os.File
	recordCount int
}

func main() {
	fmt.Println("===Agency CLI App===")

	initializeCSVFile()

	csvFile, err := os.OpenFile("./data/data.csv", os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	svc := Service{
		reader: csv.NewReader(csvFile),
		writer: csv.NewWriter(csvFile),
		csv:    csvFile,
	}

	region := flag.String("region", "no region", "region to use")
	command := flag.String("command", "no command", "the command to execute")
	flag.Parse()

	for {
		err := svc.readCSVFile()
		if err != nil {
			fmt.Println(err)
			return
		}

		svc.runCommand(*command, *region)

		fmt.Println("please enter another command or exit:")
		fmt.Scan(command)
	}
}

func (svc Service) runCommand(command string, region string) {
	switch command {
	case "list":
		svc.commandList(region)
	case "get":
		svc.commandGet(region)
	case "create":
		svc.commandCreate(region)
	case "edit":
		svc.commandEdit(region)
	case "status":
		svc.commandStatus(region)
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command not valid!")
	}
}

func (svc Service) commandList(region string) {
	var agencies []Data

	_, err := svc.csv.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = svc.reader.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		line, err := svc.reader.Read()
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			fmt.Println(err)
			return
		}

		if line[2] != region {
			id, err := strconv.Atoi(line[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
			worker, err := strconv.Atoi(line[5])
			if err != nil {
				fmt.Println(err)
				continue
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

func (svc Service) commandGet(region string) { // get a specific agency
	var agency Data
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("please enter your id:")
	scanner.Scan()
	num := scanner.Text()

	svc.csv.Seek(0, 0)
	for {
		line, err := svc.reader.Read()
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

func (svc Service) commandCreate(region string) { // create a new agency
	var agency Data
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

	fmt.Println("please enter your region:")
	scanner.Scan()
	agency.Region = scanner.Text()

	fmt.Println("please enter count of your worker:")
	scanner.Scan()
	workerString := scanner.Text()
	worker, err := strconv.Atoi(workerString)
	if err != nil {
		fmt.Println(err)
	}
	agency.Worker = worker

	agency.ID = svc.recordCount
	err = svc.writer.Write([]string{strconv.Itoa(agency.ID), agency.Name, agency.Phone, agency.Address, agency.Region, strconv.Itoa(agency.Worker)})
	if err != nil {
		fmt.Println(err)
	}
	svc.writer.Flush()
	fmt.Println(agency)
}

func (svc Service) commandEdit(region string) {
	var agencies []Data
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter your id:")
	scanner.Scan()
	num := scanner.Text()
	id, err := strconv.Atoi(num)
	if err != nil {
		fmt.Println(err)
		return
	}

	svc.csv.Seek(0, 0)
	svc.reader.Read()
	for {
		line, err := svc.reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			continue
		}

		idReal, err := strconv.Atoi(line[0])
		if err != nil {
			fmt.Println(err)
		}
		worker, err := strconv.Atoi(line[5])
		if err != nil {
			fmt.Println(err)
		}

		agency := Data{
			ID:      idReal,
			Name:    line[1],
			Region:  line[2],
			Phone:   line[3],
			Address: line[4],
			Worker:  worker,
		}
		agencies = append(agencies, agency)
	}

	found := false
	for i, agency := range agencies {
		if agency.ID == id {
			found = true
			fmt.Printf("%+v\nWhich one do you want to edit or exit to save:\n", agency)
			svc.writer.Flush()
			scanner.Scan()
			field := scanner.Text()

			switch field {
			case "Name":
				fmt.Println("Please enter the new name:")
				scanner.Scan()
				newName := scanner.Text()
				agencies[i].Name = newName
			case "Region":
				fmt.Println("Please enter the new region:")
				scanner.Scan()
				newRegion := scanner.Text()
				agencies[i].Region = newRegion
			case "Phone":
				fmt.Println("Please enter the new phone:")
				scanner.Scan()
				newPhone := scanner.Text()
				agencies[i].Phone = newPhone
			case "Address":
				fmt.Println("Please enter the new address:")
				scanner.Scan()
				newAddress := scanner.Text()
				agencies[i].Address = newAddress
			case "Worker":
				fmt.Println("Please enter the new worker count:")
				scanner.Scan()
				newWorkerString := scanner.Text()
				newWorker, err := strconv.Atoi(newWorkerString)
				if err != nil {
					fmt.Println(err)
				}
				agencies[i].Worker = newWorker
			default:
				fmt.Println("Not a valid field!!!")
			}
		}
	}

	if !found {
		fmt.Println("Record not found.")
		return
	}

	svc.csv.Truncate(0)
	svc.writer = csv.NewWriter(svc.csv)
	header := []string{"ID", "Name", "Phone", "Address", "Region", "Worker"}
	if err := svc.writer.Write(header); err != nil {
		fmt.Println(err)
		return
	}
	svc.csv.Seek(0, 0)
	for _, agency := range agencies {
		err := svc.writer.Write([]string{
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
		svc.writer.Flush()
	}
	fmt.Println("Record successfully edited and saved.")
}

func (svc Service) commandStatus(region string) {
	var workerTotal, agencies int
	svc.csv.Seek(0, 0)
	for {
		line, err := svc.reader.Read()
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

func initializeCSVFile() {
	_, err := os.Stat("./data/data.csv")
	if os.IsNotExist(err) {
		file, err := os.Create("./data/data.csv")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"ID", "Name", "Phone", "Address", "Region", "Worker"}
		if err := writer.Write(header); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (svc *Service) readCSVFile() error {
	_, err := svc.csv.Seek(0, 0)
	if err != nil {
		return err
	}

	records, err := svc.reader.ReadAll()
	if err != nil {
		return err
	}

	svc.recordCount = len(records)
	return nil
}
