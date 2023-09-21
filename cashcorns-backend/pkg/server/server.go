package server

import (
	"fmt"
	"net/http"

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

	http.HandleFunc("/", homeHandler)

	// Start the server on port 8080
	fmt.Printf("Server is listening on port %s...\n", s.Port)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%s", s.Port), "tls.crt", "tls.key", nil)
	//err := http.ListenAndServe(fmt.Sprintf(":%s", s.Port), nil)
	if err != nil {
		panic(err)
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the home page!")
}
