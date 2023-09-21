package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func UnmarshallFileToStruct(filePath string, target interface{}) error {
	// Read the file
	file, err := os.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("unable to read file %s. Err: %s", filePath, err.Error())
		return err
	}

	// Unmarshall the file into the target struct
	err = json.Unmarshal(file, &target)
	if err != nil {
		err = fmt.Errorf("unable to unmarshall file %s into struct. Err: %s", filePath, err.Error())
		return err
	}

	return nil

}

// parseDate parses a date string in "YYYY-MM-DD" format and returns a time.Time value.
func parseDate(dateString string) (time.Time, error) {
	return time.Parse("2006-01-02", dateString)
}

// FindClosestPastDate finds the closest past date in the dateStrings array to the given date.
func FindClosestPastDate(dateStrings []string, date time.Time) (string, error) {
	var closestDate string
	var closestDiff time.Duration
	// now := date

	for _, dateString := range dateStrings {
		parsedDate, err := parseDate(dateString)
		if err != nil {
			return "", err
		}

		diff := date.Sub(parsedDate)
		if diff > 0 && (closestDiff == 0 || diff < closestDiff) {
			closestDate = dateString
			closestDiff = diff
		}
	}

	if closestDate != "" {
		return closestDate, nil
	}

	return "", fmt.Errorf("no past dates found in the input")
}
