package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Map brands to bike types
func getBikeType(title string) string {
	title = strings.ToLower(title)

	// Mapping of brands and models to their likely bike types
	brandBikeTypes := map[string]string{
		"trek":                 "Road Bike",
		"schwinn":              "Cruiser Bike",
		"wheel speed":          "Electric Bike",
		"rad power":            "Electric Bike",
		"igo":                  "Electric Bike",
		"emmo":                 "Electric Bike",
		"decathlon":            "Hybrid Bike",
		"norco":                "Mountain Bike",
		"raleigh":              "Hybrid Bike",
		"brodie":               "Hybrid Bike",
		"supercycle":           "Mountain Bike",
		"gt bicycles":          "Mountain Bike",
		"aostirmotor":          "Electric Bike",
		"rocky mountain":       "Mountain Bike",
		"specialized":          "Hybrid Bike",
		"kona":                 "Hybrid Bike",
		"giant":                "Road Bike",
		"devinci":              "Mountain Bike",
		"ccm":                  "Hybrid Bike",
		"northrock":            "Mountain Bike",
		"diamondback":          "Mountain Bike",
		"opus":                 "Hybrid Bike",
		"fx":                   "Road Bike",
		"storm":                "Mountain Bike",
		"dew":                  "Hybrid Bike",
		"hardrock":             "Mountain Bike",
		"aggressor":            "Mountain Bike",
		"checkpoint":           "Gravel Bike",
		"diadora":              "Hybrid Bike",
		"ghost":                "Mountain Bike",
		"kato":                 "Mountain Bike",
		"eastern":              "BMX Bike",
		"gotrax":               "Electric Scooter",
		"felt":                 "Road Bike",
		"cinder cone":          "Mountain Bike",
		"fuji":                 "Road Bike",
		"talon":                "Mountain Bike",
		"nakamura":             "Mountain Bike",
		"cannondale":           "Road Bike",
		"ctm":                  "Mountain Bike",
		"reebok":               "Hybrid Bike",
		"sonar":                "Hybrid Bike",
		"enduro":               "Mountain Bike",
		"marin":                "Mountain Bike",
		"mongoose":             "Mountain Bike",
		"rize":                 "Electric Bike",
		"jamis":                "Road Bike",
		"s-works":              "Road Bike",
		"marlin":               "Mountain Bike",
		"stagger":              "Road Bike",
		"sanctuary":            "Cruiser Bike",
		"honey stinger":        "Mountain Bike", // Stinga
		"aquila":               "Road Bike",     // B-Drive
		"santa cruz":           "Mountain Bike", // Heckler
		"surface 604":          "Electric Bike", // Surface 604
		"khs bicycles":         "Mountain Bike", // Vitamin A, Winslow
		"evo":                  "Hybrid Bike",   // Evo
		"scott":                "Road Bike",     // Foil Contessa
		"subrosa":              "BMX Bike",      // Novus
		"gary fisher":          "Mountain Bike", // Wahoo
		"colnago":              "Road Bike",     // Air
		"escape":               "Hybrid Bike",   // Eternity
		"miyata":               "Road Bike",     // Miyata (general)
		"cult":                 "BMX Bike",      // Gateway
		"haro":                 "BMX Bike",      // Demon
		"iron horse bicycles":  "Mountain Bike", // 3.0
		"linus":                "City Bike",     // Rambler
		"sekine":               "Road Bike",     // Sekine (vintage, general)
		"lemond racing cycles": "Road Bike",     // Tourmalet
		"pro":                  "Road Bike",     // General for Pro
	}

	// Check title for known brands and return corresponding bike type
	for brand, bikeType := range brandBikeTypes {
		if strings.Contains(title, brand) {
			return bikeType
		}
	}

	// Default if no brand match is found
	return "Other"
}

func main() {
	// Open the CSV file for reading
	inputFile, err := os.Open("updated_bike_types_final.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer inputFile.Close()

	// Create a CSV reader
	reader := csv.NewReader(inputFile)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Update bike types based on the title in each record
	for i, record := range records {
		if i == 0 {
			// Skip header row
			continue
		}

		// Assume "Title" is in the first column and "Bike Type" is in the last column
		title := record[0]
		// Update the bike type if it is "Other" or empty
		if record[len(record)-1] == "Other" {
			record[len(record)-1] = getBikeType(title)
		}
	}

	// Create a new CSV file for output
	outputFile, err := os.Create("updated_bike_types_final_final.csv")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Write the updated records to the new CSV file
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Write all records back to the new CSV file
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error writing record to file:", err)
			return
		}
	}

	fmt.Println("Bike types updated successfully and saved to 'updated_bikes.csv'")
}
