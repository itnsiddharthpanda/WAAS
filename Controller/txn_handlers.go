package Controller

import (
	"net/http"
	"waas/Domain"
)

func GenerateCSV(w http.ResponseWriter, r *http.Request) {
	go Domain.GenerateCSV()
}