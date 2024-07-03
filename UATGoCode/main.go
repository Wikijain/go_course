package main

import (
	"encoding/csv"
	"encoding/json"
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

	for _, postgresqlKey := range postgresqlKeys {
		postgresqlID := postgresqlKey[0]
		postgresqlColumn2 := postgresqlKey[1]

		matchFound := false

		for _, cassandraRecord := range cassandraData {
			jsonData := cassandraRecord[0].(string)
			//arrayData := cassandraRecord[1].([]interface{})
			if strings.Contains(jsonData, postgresqlID) {
				//log.Printf("Found '%s' \n", postgresqlID)
				matchFound = true
				break
			}
		}

		if !matchFound {
			log.Printf("No match found for PostgreSQL id '%s', '%s'", postgresqlID, postgresqlColumn2)
		}
	}

	log.Println("Data comparison process completed.")

}
