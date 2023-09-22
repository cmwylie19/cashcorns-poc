package server

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	pay "github.com/cmwylie19/cashcorns-backend/pkg/payrun"
)

type Server struct {
	Port string
	pay.PayRun
}

func NewServer(port, fileLocation string) *Server {
	return &Server{
		Port: port,
		PayRun: pay.PayRun{
			FileLocation: fileLocation,
		},
	}
}

func (s *Server) Serve() {
	fmt.Println("serve called")

	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	// Start the server on port 8080
	fmt.Printf("Server is listening on port %s...\n", s.Port)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%s", s.Port), "tls.crt", "tls.key", nil)
	//err := http.ListenAndServe(fmt.Sprintf(":%s", s.Port), nil)
	if err != nil {
		panic(err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handle file upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Decode the JSON file
	var data struct {
		Leaderboard []struct {
			Username string `json:"username"`
			Sum      string `json:"sum"`
		} `json:"leaderboard"`
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a CSV file
	csvFile, err := os.Create("output.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer csvFile.Close()

	// Write CSV data
	csvWriter := csv.NewWriter(csvFile)
	_ = csvWriter.Write([]string{"Rippling Emp No", "Employee Name", "Mission Hero", "Project", "Service Item", "Shift ID", "Salary Hours", "Salary Amount", "Dental Insurance (Imputed)", "Expenses Reimbursement", "Life Insurance (Imputed)", "Medical Insurance (Imputed)", "NSO Exercise Income", "Office Expense Reimbursement (Recurring)", "Reimbursements", "Taco Bonus", "Vision Insurance (Imputed)", "Dental Deductions (Dental) Emp. Deduction", "Dental Deductions (Dental) Co. Contribution", "401K (Pre-tax %) (401K) Emp. Deduction", "401K (Pre-tax %) (401K) Co. Contribution", "Vision Deductions (Vision) Emp. Deduction", "Vision Deductions (Vision) Co. Contribution", "Medical Deductions (Medical) Emp. Deduction", "Medical Deductions (Medical) Co. Contribution", "Spend Management - Card Repayment (Post Tax) Emp. Deduction", "Spend Management - Card Repayment (Post Tax) Co. Contribution", "401K (Post-tax %) (Roth 401K) Emp. Deduction", "401K (Post-tax %) (Roth 401K) Co. Contribution", "Life Deductions (Life) Emp. Deduction", "Life Deductions (Life) Co. Contribution", "Long Term Disability Deductions (Long Term Disability) Emp. Deduction", "Long Term Disability Deductions (Long Term Disability) Co. Contribution", "Payroll Correction (Post Tax) Emp. Deduction", "Payroll Correction (Post Tax) Co. Contribution", "HSA (HSA) Emp. Deduction", "HSA (HSA) Co. Contribution", "Short Term Disability Deductions (Short Term Disability) Emp. Deduction", "Short Term Disability Deductions (Short Term Disability) Co. Contribution", "Loan Repayment (Post-Tax) (Post Tax) Emp. Deduction", "Loan Repayment (Post-Tax) (Post Tax) Co. Contribution", "Negative Net Pay Payment Plan (Post Tax) Emp. Deduction", "Negative Net Pay Payment Plan (Post Tax) Co. Contribution", "Michael N Palassis_Case ID: 7105103332", "Paystub Note"})
	for _, item := range data.Leaderboard {
		_, exists := nameMap[item.Username]
		if exists {
			record := []string{nameMap[item.Username], item.Sum, endCommas}
			if err := csvWriter.Write(record); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}
	defer csvWriter.Flush()

	// Check for errors in writing CSV
	if err := csvWriter.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// _ = removeDoubleQuotesFromFile("output.csv")

	// Respond with the download link
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<html>
	<head>
	<style>
	body {
		/* Set the background image URL */
		background-image: url('https://www.defenseunicorns.com/static/home-hero-bg-feff54111a83d6b7c0c38219c80b1333.jpg');
	
		/* Set the background image size */
		background-size: cover;
	
		/* Center the background image horizontally and vertically */
		background-position: center;
	
		/* Make the background image fixed (doesn't scroll with the content) */
		background-attachment: fixed;
	
		/* Specify a background color in case the image doesn't load or cover the entire viewport */
		background-color: #f0f0f0; /* Replace with your desired background color */
		display: flex;
		justify-content: center;
		align-items: center;
	}
	a {
		text-decoration: none;
		color: #fff;
		font-size: 2rem;
		transition: transform 0.3s ease-in-out;
	}
	a:hover {
		transform: scale(1.2) rotate(360deg);
	}
	</style>
	</head>
	
	<body><a href=\"/download\">Download Me 🦄</a></body></html>`))
}

func removeDoubleQuotesFromFile(filePath string) error {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a temporary file to store the modified content
	tempFilePath := filePath + ".temp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Remove double quotes from the line
		modifiedLine := strings.Replace(line, `"`, "", -1)

		// Write the modified line to the temporary file
		_, err := tempFile.WriteString(modifiedLine + "\n")
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Close the temporary file
	tempFile.Close()

	// Rename the temporary file to the original file
	if err := os.Rename(tempFilePath, filePath); err != nil {
		return err
	}

	return nil
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Specify the file you want to send for download
	filePath := "output.csv" // Replace with the path to your file
	_ = removeDoubleQuotesFromFile(filePath)
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the Content-Disposition header to suggest a filename for the download
	w.Header().Set("Content-Disposition", "attachment; filename="+filePath)

	// Set the Content-Type header based on the file's type or specify "application/octet-stream" for binary files
	w.Header().Set("Content-Type", "text/plain") // Replace with the appropriate MIME type

	// Copy the file's content to the HTTP response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
