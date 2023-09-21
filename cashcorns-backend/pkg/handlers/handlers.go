package handlers

import (
	"fmt"
	"net/http"
)

func LastPayHandler(w http.ResponseWriter, r *http.Request) {
	// Read the dates file

}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the about page.")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Contact us at example@example.com.")
}
