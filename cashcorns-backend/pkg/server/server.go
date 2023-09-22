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

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

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

	csvFile, err := os.Create("output.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer csvFile.Close()

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

	if err := csvWriter.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<html>
	<head>
	<title>🦄</title>
	<style>
	body {
	
		background-image: url('https://www.defenseunicorns.com/static/home-hero-bg-feff54111a83d6b7c0c38219c80b1333.jpg');
	
		background-size: cover;
		font-family: sans-serif;
		background-position: center;
	
		background-attachment: fixed;
	

		background-color: #f0f0f0; 
		display: flex;
		justify-content: center;
		align-items: center;
	}
	a {
		text-decoration: none;
		color: red;
		font-size: 2rem;
		transition: transform 1.5s ease-in-out;
	}
	a:hover {
		transform: scale(1.2) rotate(360deg);
	}
	</style>
	</head>
	
	<body><a href="/download">Download Me 🦄</a></body></html>`))
}

func removeDoubleQuotesFromFile(filePath string) error {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tempFilePath := filePath + ".temp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		modifiedLine := strings.Replace(line, `"`, "", -1)

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

	filePath := "output.csv"
	_ = removeDoubleQuotesFromFile(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filePath)

	w.Header().Set("Content-Type", "text/csv")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
