package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ApprovalDates are the array dates when tacos need to be approved
type ApprovalDates struct {
	ApprovalDates []string `json:"approvalDates"`
}

var approvalDates = []string{
	"2023-09-19",
	"2023-09-27",
	"2023-10-11",
}

func TestUnmarshallFileToStruct(t *testing.T) {
	var approvalDates ApprovalDates
	err := UnmarshallFileToStruct("../../ApprovalDates.json", &approvalDates)
	require.NoError(t, err)
	require.Equal(t, approvalDates.ApprovalDates[0], "2023-09-19")
	require.Equal(t, approvalDates.ApprovalDates[1], "2023-09-27")
	require.Equal(t, approvalDates.ApprovalDates[2], "2023-10-11")
}
func TestFindClosestPastDate(t *testing.T) {
	// expectedResult[0] - current
	// expectedResult[1] - expected
	var expectedResult = [][]string{
		{"2023-09-21", "2023-09-19"},
		{"2023-09-22", "2023-09-19"},
		{"2023-09-26", "2023-09-19"},
		{"2023-09-29", "2023-09-27"},
		{"2023-10-01", "2023-09-27"},
		{"2023-10-04", "2023-09-27"},
		{"2023-10-17", "2023-10-11"},
		{"2023-10-18", "2023-10-11"},
		{"2023-10-19", "2023-10-11"},
	}

	for idx, dateString := range expectedResult {
		date, _ := parseDate(dateString[0])
		newDate, err := FindClosestPastDate(approvalDates, date)
		require.NoError(t, err)
		got := expectedResult[idx][1]
		actual := newDate
		require.Equal(t, got, actual)
	}
}
