package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Endpoint falso para entregar el token
	http.HandleFunc("/api/token/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access":"mock-token","refresh":"mock-refresh"}`))
	})

	// Endpoint falso para recibir la inyección
	http.HandleFunc("/api/inyeccion-salesforce/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"StatusCode":200,"Mensaje":"Inyeccion simulada exitosamente (Mock)","NumeroSolicitud":99999}`))
	})

	fmt.Println("===================================================")
	fmt.Println("🟢 Servidor Mock (Dummy) de FinnFlow iniciado")
	fmt.Println("🟢 Escuchando en: http://localhost:8083")
	fmt.Println("===================================================")
	
	server := &http.Server{
		Addr:              ":8083",
		ReadHeaderTimeout: 3 * time.Second,
	}
	
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error iniciando el mock:", err)
	}
}
