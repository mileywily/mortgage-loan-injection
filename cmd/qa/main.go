package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("Iniciando Test de Paridad Extrema (Legacy vs Go)...")

	go startMockFinnFlow()
	time.Sleep(1 * time.Second)

	legacyURL := "http://127.0.0.1:8080/v1/bfcl/mortgage-loan/injections"
	goURL := "http://127.0.0.1:8082/v1/bfcl/mortgage-loan/injections"

	testCases := []struct {
		name    string
		payload string
	}{
		{
			name: "1. Happy Path",
			payload: `{
				"datos_credito": {"NumeroSolicitud": 12345,"AntiguedadVivienda": 5,"EjecutivoComercial": "Juan","Producto": 1,"Objetivo": 2,"Destino": 3,"MontoAprobado": 100.5,"ValorPropiedad": 200.5,"FechaAprobacion": "2023-01-01","ValorContado": 100.0,"Plazo1": 20,"MesesGracia": 2,"Tasa1": 3.5,"Spread1": 1.0},
				"participantes": [{"Rut": "12.345.678-9","TipoParticipacion": 1,"Nombre": "A","Paterno": "B","Materno": "C","FechaNacimiento": "01-01-1990"}]
			}`,
		},
		{
			name: "2. Missing participantes",
			payload: `{
				"datos_credito": {"NumeroSolicitud": 12345,"AntiguedadVivienda": 5,"EjecutivoComercial": "Juan","Producto": 1,"Objetivo": 2,"Destino": 3,"MontoAprobado": 100.5,"ValorPropiedad": 200.5,"FechaAprobacion": "2023-01-01","ValorContado": 100.0,"Plazo1": 20,"MesesGracia": 2,"Tasa1": 3.5,"Spread1": 1.0},
				"participantes": []
			}`,
		},
		{
			name: "3. Invalid Format (Rut)",
			payload: `{
				"datos_credito": {"NumeroSolicitud": 12345,"AntiguedadVivienda": 5,"EjecutivoComercial": "Juan","Producto": 1,"Objetivo": 2,"Destino": 3,"MontoAprobado": 100.5,"ValorPropiedad": 200.5,"FechaAprobacion": "2023-01-01","ValorContado": 100.0,"Plazo1": 20,"MesesGracia": 2,"Tasa1": 3.5,"Spread1": 1.0},
				"participantes": [{"Rut": "invalid-rut","TipoParticipacion": 1,"Nombre": "A","Paterno": "B","Materno": "C","FechaNacimiento": "01-01-1990"}]
			}`,
		},
		{
			name: "4. Invalid Value (Plazo1)",
			payload: `{
				"datos_credito": {"NumeroSolicitud": 12345,"AntiguedadVivienda": 5,"EjecutivoComercial": "Juan","Producto": 1,"Objetivo": 2,"Destino": 3,"MontoAprobado": 100.5,"ValorPropiedad": 200.5,"FechaAprobacion": "2023-01-01","ValorContado": 100.0,"Plazo1": 12,"MesesGracia": 2,"Tasa1": 3.5,"Spread1": 1.0},
				"participantes": [{"Rut": "12.345.678-9","TipoParticipacion": 1,"Nombre": "A","Paterno": "B","Materno": "C","FechaNacimiento": "01-01-1990"}]
			}`,
		},
		{
			name: "5. Renegociacion (SRV01) con formato original invalido",
			payload: `{
				"datos_credito": {"NumeroSolicitud": 12345,"AntiguedadVivienda": 5,"EjecutivoComercial": "Juan","Producto": 1,"Objetivo": 2,"Destino": 3,"MontoAprobado": 100.5,"ValorPropiedad": 200.5,"FechaAprobacion": "2023-01-01","ValorContado": 100.0,"Plazo1": 20,"MesesGracia": 2,"Tasa1": 3.5,"Spread1": 1.0, "NumeroOperacionOriginal": "12345678901234567"},
				"participantes": [{"Rut": "12.345.678-9","TipoParticipacion": 1,"Nombre": "A","Paterno": "B","Materno": "C","FechaNacimiento": "01-01-1990"}]
			}`,
		},
	}

	success := true

	for _, tc := range testCases {
		fmt.Printf("\n--- Corriendo Test: %s ---\n", tc.name)
		legacyStatus, legacyBody := sendReq(legacyURL, tc.payload)
		goStatus, goBody := sendReq(goURL, tc.payload)

		fmt.Printf("Legacy -> Status: %d | Body: %s\n", legacyStatus, legacyBody)
		fmt.Printf("Go     -> Status: %d | Body: %s\n", goStatus, goBody)

		if legacyStatus == 0 || goStatus == 0 {
			fmt.Println("ERROR DE CONEXION: Uno de los servidores (Java o Go) no respondió. Revisa si se levantaron correctamente.")
			success = false
		} else if legacyStatus != goStatus {
			fmt.Printf("ERROR: Status codes no coinciden. Legacy: %d, Go: %d\n", legacyStatus, goStatus)
			success = false
		} else {
			fmt.Println("Status codes coinciden.")
		}
	}

	if !success {
		fmt.Println("\nFallaron algunas pruebas de paridad.")
		os.Exit(1)
	}
	fmt.Println("\nTODAS LAS PRUEBAS DE PARIDAD EXTREMA PASARON CON EXITO.")
}

func sendReq(url, payload string) (int, string) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Sprintf("Error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func startMockFinnFlow() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/token/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access":"mock-token","refresh":"mock-refresh"}`))
	})
	mux.HandleFunc("/api/inyeccion-salesforce/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := io.ReadAll(r.Body)
		
		if bytes.Contains(body, []byte(`"participantes": []`)) {
			w.WriteHeader(http.StatusNotAcceptable) // Legacy returns what FinnFlow returns. Go returns 406.
			w.Write([]byte(`{"code":"invalid_value","message":"Participantes with invalid value"}`))
			return
		}
		if bytes.Contains(body, []byte(`"invalid-rut"`)) {
			w.WriteHeader(http.StatusPaymentRequired)
			w.Write([]byte(`{"code":"invalid_format","message":"Rut with invalid format"}`))
			return
		}
		if bytes.Contains(body, []byte(`"Plazo1": 12`)) {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(`{"code":"invalid_value","message":"Plazo1 with invalid value"}`))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"Mensaje":"Application injected successfully","NumeroSolicitud":12345}`))
	})
	http.ListenAndServe("0.0.0.0:8083", mux)
}
