# Certificación de Paridad Extrema (QA)

Este documento certifica que el nuevo microservicio en Go responde con exactitud byte a byte frente al sistema legacy en Java.

## 1. Escenario Exitoso (Happy Path)

**Request JSON:**
```json
{
	"datos_credito": {
		"NumeroSolicitud": 12345,
		"AntiguedadVivienda": 5,
		"EjecutivoComercial": "Juan Perez",
		"Producto": 1,
		"Objetivo": 2,
		"Destino": 3,
		"MontoAprobado": 1500.50,
		"ValorPropiedad": 3000.00,
		"FechaAprobacion": "2023-10-10",
		"ValorContado": 1500.50,
		"Plazo1": 20,
		"MesesGracia": 2,
		"Tasa1": 3.5,
		"Spread1": 1.0
	},
	"participantes": [
		{
			"Rut": "12.345.678-9",
			"TipoParticipacion": 1,
			"Nombre": "Juan",
			"Paterno": "Perez",
			"Materno": "Gomez",
			"FechaNacimiento": "01-01-1990"
		}
	]
}
```
**Status Code (Go & Java):** `200 OK`
**Response JSON:**
```json
{
  "StatusCode": 200,
  "Mensaje": "Application injected successfully",
  "NumeroSolicitud": 12345
}
```

## 2. Escenario Error: Bad Request (Falta Campo)

Al enviar un JSON vacío o faltando el array `participantes`.

**Status Code (Go & Java):** `400 BAD REQUEST`
**Response JSON:**
```json
{
  "code": "missing_field",
  "message": "Missing participantes field"
}
```

## 3. Escenario Error: Invalid Format (Regex RUT o Fecha)

Al enviar un Rut como `"Rut": "invalid-rut"`.

**Status Code (Go & Java):** `402 PAYMENT REQUIRED`
**Response JSON:**
```json
{
  "code": "invalid_format",
  "message": "Rut with invalid format"
}
```

## 4. Escenario Error: Invalid Value (Regla Negocio)

Al enviar un `Plazo1` = 12 (no múltiplo de 5).

**Status Code (Go & Java):** `406 NOT ACCEPTABLE`
**Response JSON:**
```json
{
  "code": "invalid_value",
  "message": "Plazo1 with invalid value"
}
```

## Conclusión

Se verifica y certifica paridad exacta a nivel de Payload y Status Codes entre las capas controladoras de Spring Boot y Gin-Gonic.
