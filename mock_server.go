package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/token/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access":"mock-token","refresh":"mock-refresh"}`))
	})
	mux.HandleFunc("/api/inyeccion-salesforce/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := io.ReadAll(r.Body)
		
		fmt.Println("\n=== 3. PAYLOAD RECIBIDO POR FINNFLOW (Mock) ===")
		fmt.Println(string(body))
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"Mensaje":"Application injected successfully","NumeroSolicitud":12345}`))
	})
	fmt.Println("Mock FinnFlow escuchando en 8083...")
	http.ListenAndServe("0.0.0.0:8083", mux)
}
