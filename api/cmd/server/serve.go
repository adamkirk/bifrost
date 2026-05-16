package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func handleServe(_ *cobra.Command, _ []string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	return http.ListenAndServe(":8080", nil)
}
