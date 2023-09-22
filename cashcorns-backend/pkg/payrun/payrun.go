package payrun

import (
	"fmt"

	"github.com/cmwylie19/cashcorns-backend/pkg/helpers"
)

// ApprovalDates are the array dates when tacos need to be approved
type ApprovalDates struct {
	ApprovalDates []string `json:"approvalDates"`
}

// PayRun is the struct that holds the payrun data
type PayRun struct {
	ApprovalDates ApprovalDates `json:"approvalDates"`
	MostRecent    string        `json:"mostRecent"`
	FileLocation  string
}

// Load loads the approvalDates data from the file
func (pr *PayRun) Load() error {

	err := helpers.UnmarshallFileToStruct(pr.FileLocation, &pr.ApprovalDates)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded approval dates from file: %s", pr.FileLocation)

	return nil
}
func (pr *PayRun) GetLastPayDate() string {
	return pr.MostRecent
}
