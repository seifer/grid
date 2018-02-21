package main

import (
	"net/http"

	"github.com/seifer/grid/grid"
)

func main() {
	http.Handle("/fetch", &grid.Handler{})
	http.ListenAndServe(":8080", nil)
}
