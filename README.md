# GoAgencyOps
## Overview
This project is a small command-line interface (CLI) program written in Golang to manage representatives and their information in different regions of a country. The CLI supports the following flags: `command` and `region`.

## Supported Flags
**region**: Specifies the region (e.g., Tehran, Isfahan, Shiraz, New York) for which the command will be executed.

**command**: Specifies the action to be performed. It can take one of the following values:

Command  | Second Header
------------- | -------------
list  | Lists the representatives in the specified region.
get  |  Retrieves detailed information about a specific representative in the given region. Requires the user to provide the representative's ID.
create | Creates a new representative with information such as name, address, phone number, membership date, and the number of employees.
edit  | Edits the information of an existing representative. Requires the user to input the ID of the representative through standard input.
status |  Displays the number of representatives and total employees in the specified region.

## Usage 
```
go run agencyCliApp.go --command <command> --region <region>
```
or
```
./cliApp --command <command> --region <region>
```
  



