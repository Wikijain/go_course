package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

)

// readPostgreSQLCSV reads a PostgreSQL CSV file and returns a list of keys.
func readPostgreSQLCSV(filename string) [][]string {
	keys := [][]string{}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read() // Skip header

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if len(row) > 0 {
			keys = append(keys, []string{row[2], row[1]})
		}
	}
	log.Printf("Read %d keys from PostgreSQL CSV", len(keys))
	//log.Println(keys)
	return keys
}

// readCassandraCSV reads a Cassandra CSV file and returns a list of JSON data and array data.
func readCassandraCSV(filename string) [][]interface{} {
	data := [][]interface{}{}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read() // Skip header

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if len(row) > 1 {
			jsonData := row[0]
			var arrayData []interface{}
			err := json.Unmarshal([]byte(row[1]), &arrayData)
			if err != nil {
				log.Printf("Failed to parse JSON data: %s", jsonData)
				continue
			}
			data = append(data, []interface{}{jsonData, arrayData})
		}
	}

	log.Printf("Read %d records from Cassandra CSV", len(data))
	//log.Println(data)
	return data
}

// main is the entry point of the program.
func main() {
	log.Println("Starting data comparison process...")

	postgresqlCSV := "output_postgresql.csv"
	cassandraCSV := "output_cassandra.csv"

	postgresqlKeys := readPostgreSQLCSV(postgresqlCSV)
	cassandraData := readCassandraCSV(cassandraCSV)
	// Create or open the CSV file
	file, err := os.Create("output_log.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all data is written

	// Write the header row (optional)
	header := []string{"errorCode", "errorMsg"}
	writer.Write(header)

	// create data rows
	data := [][]string{}

	for _, postgresqlKey := range postgresqlKeys {
		postgresqlID := postgresqlKey[0]
		postgresqlColumn2 := postgresqlKey[1]

		matchFoundID := false
		matchFoundCategory := false
		for _, cassandraRecord := range cassandraData {
			jsonData := cassandraRecord[0].(string)
			arrayData := cassandraRecord[1].([]interface{})
			if strings.Contains(jsonData, postgresqlID) {
				matchFoundID = true
			}
			var cassandraJSON map[string]interface{}
			err := json.Unmarshal([]byte(jsonData), &cassandraJSON)
			if err != nil {
				log.Printf("Failed to parse JSON data: %s", jsonData)
				continue
			}
			cassandraID, ok := cassandraJSON["id"].(string)
			if !ok {
				log.Printf("Invalid ID format in Cassandra record: %v", cassandraJSON)
				continue
			}
			if postgresqlID == cassandraID {
				if arrayContainsString(arrayData, postgresqlColumn2) {
					matchFoundCategory = true
				}
			}
		}

		if !matchFoundID && !matchFoundCategory {
			log.Printf("No match found for PostgreSQL id '%s',and Category '%s'", postgresqlID, postgresqlColumn2)
			errorMsg := fmt.Sprintf("No match found for PostgreSQL id '%s',and Category '%s'", postgresqlID, postgresqlColumn2)
			myslice := []string{"3", errorMsg}
			data = append(data, myslice)
		} else if !matchFoundID {
			log.Printf("No match found for PostgreSQL id '%s'", postgresqlID)
			errorMsg := fmt.Sprintf("No match found for PostgreSQL id '%s'", postgresqlID)
			myslice := []string{"1", errorMsg}
			data = append(data, myslice)
		} else if !matchFoundCategory {
			log.Printf("No match found for  Category '%s'", postgresqlColumn2)
			errorMsg := fmt.Sprintf("No match found for  Category '%s'", postgresqlColumn2)
			myslice := []string{"2", errorMsg}
			data = append(data, myslice)
		}
	}

	log.Println("Data comparison process completed.")
	for _, row := range data {
		err := writer.Write(row)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// arrayContainsString checks if a string is contained in an array of interfaces.
func arrayContainsString(array []interface{}, str string) bool {
	for _, elem := range array {
		s := fmt.Sprintf("%v", elem)
		if strings.Contains(s, str) {
			return true
		}
	}
	return false
}
