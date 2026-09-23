package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	payload := `{
		"datos_credito": {
			"NumeroSolicitud": 12345,
			"AntiguedadVivienda": 5,
			"EjecutivoComercial": "Juan",
			"Producto": 1,
			"Objetivo": 2,
			"Destino": 3,
			"MontoAprobado": 100.5,
			"ValorPropiedad": 200.5,
			"FechaAprobacion": "2023-01-01",
			"ValorContado": 100.0,
			"Plazo1": 20,
			"MesesGracia": 2,
			"Tasa1": 3.5,
			"Spread1": 1.0,

			"NumeroOperacionOriginal": "OPE-RENEG-001",
			"MontoAbono": 19.2000,
			"IncluirGastosOperacionales": 1,
			"TipoGarantia": 2,
			"NumeroRenegociacion": 1,
			"SubsidioOriginal": 1,
			"ReduccionInteresesYGC": 10
		},
		"participantes": [
			{
				"Rut": "12.345.678-9",
				"TipoParticipacion": 1,
				"Nombre": "A",
				"Paterno": "B",
				"Materno": "C",
				"FechaNacimiento": "01-01-1990"
			}
		]
	}`

	fmt.Println("=== 1. PAYLOAD ENVIADO A GO (SRV01 - Renegociacion) ===")
	fmt.Println(payload)

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8082/v1/bfcl/mortgage-loan/injections", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	
	if err != nil {
		fmt.Println("Error conectando a Go:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("\n=== 2. RESPUESTA DEL SERVIDOR GO (Status %d) ===\n", resp.StatusCode)
	fmt.Println(string(body))
}
