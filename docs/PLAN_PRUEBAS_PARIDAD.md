# Plan de Pruebas y Validación de Paridad
**Microservicio:** `api-bfcl-mortgage-loans-injection`  
**Versión Go:** `v1.2.0`  
**Endpoint principal:** `POST /v1/bfcl/mortgage-loan/injections`  
**Versión del documento:** 1.0  

---

## Índice
1. [Ambiente de Pruebas](#1-ambiente-de-pruebas)
2. [Estructura del Request](#2-estructura-del-request)
3. [Estructura del Response](#3-estructura-del-response)
4. [Casos de Uso - Capa de Validación HTTP](#4-casos-de-uso---capa-de-validación-http)
5. [Casos de Uso - Reglas de Negocio](#5-casos-de-uso---reglas-de-negocio)
6. [Casos de Uso - Integración FinnFlow](#6-casos-de-uso---integración-finnflow)
7. [Casos de Uso - Campos Opcionales](#7-casos-de-uso---campos-opcionales)
8. [Tabla Resumen de Paridad](#8-tabla-resumen-de-paridad)
9. [Checklist Final](#9-checklist-final)

---

## 1. Ambiente de Pruebas

### Local con Mock (Sin credenciales)
| Componente | URL |
|---|---|
| API Go | `http://localhost:8081` |
| Mock FinnFlow | `http://localhost:8083` |

**Encender ambiente local:**
```bash
# Pestaña 1 - Mock
go run cmd/mock/main.go

# Pestaña 2 (Mac/Linux)
PORT=8081 FINNFLOW_URL="http://127.0.0.1:8083" FINNFLOW_KEY="dummy" FINNFLOW_SECRET="dummy" go run cmd/api/main.go

# Pestaña 2 (Windows PowerShell)
$env:PORT="8081"; $env:FINNFLOW_URL="http://127.0.0.1:8083"; $env:FINNFLOW_KEY="dummy"; $env:FINNFLOW_SECRET="dummy"; go run cmd/api/main.go
```

### QA Real (Con credenciales Nullplatform)
| Variable | Descripción |
|---|---|
| `PORT` | Puerto de escucha (ej. `8081`) |
| `FINNFLOW_URL` | URL base de FinnFlow QA |
| `FINNFLOW_KEY` | Usuario de autenticación FinnFlow |
| `FINNFLOW_SECRET` | Contraseña de autenticación FinnFlow |

---

## 2. Estructura del Request

```json
{
  "datos_credito": {
    "NumeroSolicitud": 192269,       // int - REQUERIDO
    "AntiguedadVivienda": 1,         // int - REQUERIDO
    "EjecutivoComercial": "PSOTO",   // string - REQUERIDO (máx. 20 chars)
    "Producto": 1,                   // int - REQUERIDO
    "Objetivo": 1,                   // int - REQUERIDO
    "Destino": 1,                    // int - REQUERIDO
    "FechaAprobacion": "2024-01-15", // string - REQUERIDO (formato YYYY-MM-DD)
    "MontoAprobado": 1688.09,        // float - REQUERIDO, > 0
    "ValorPropiedad": 2300.00,       // float - REQUERIDO, > 0
    "ValorContado": 261.91,          // float - REQUERIDO, > 0
    "Plazo1": 30,                    // int - REQUERIDO (5, 10, 15, 20, 25 o 30)
    "MesesGracia": 0,                // int - REQUERIDO
    "Tasa1": 3.99,                   // float - REQUERIDO (0.1 a 50.0)
    "Spread1": 0.0,                  // float - REQUERIDO, >= 0

    // Campos opcionales (Solo para Renegociación)
    "NumeroOperacionOriginal": "1234567890123456", // string (máx. 16 chars)
    "MontoAbono": 500.00,            // float
    "IncluirGastosOperacionales": 1, // int
    "TipoGarantia": 1,              // int
    "NumeroRenegociacion": 1,        // int
    "SubsidioOriginal": 0,           // int
    "ReduccionInteresesYGC": 0       // int
  },
  "participantes": [
    {
      "Rut": "18.015.051-K",         // string - REQUERIDO (Formato: XX.XXX.XXX-X)
      "TipoParticipacion": 1,        // int - REQUERIDO
      "Nombre": "Pedro",             // string - REQUERIDO (máx. 100 chars)
      "Paterno": "Perez",            // string - REQUERIDO (máx. 100 chars)
      "Materno": "Perez",            // string - REQUERIDO (máx. 100 chars)
      "FechaNacimiento": "04-07-1992", // string - REQUERIDO (Formato: DD-MM-YYYY)
      "Email": "p@email.com"         // string - OPCIONAL (email válido, máx. 100)
    }
  ],
  "propiedades": [                   // OPCIONAL
    {
      "TipoInmueble": 1,             // int - REQUERIDO si hay propiedades
      "Antiguedad": 10,              // int - REQUERIDO si hay propiedades
      "Direccion": "Av. Falabella",  // string - REQUERIDO (máx. 100 chars)
      "Numero": 1234,                // int - REQUERIDO si hay propiedades
      "Depto": "201",                // string - OPCIONAL (máx. 20 chars)
      "Comuna": 13101                // int - REQUERIDO si hay propiedades
    }
  ]
}
```

---

## 3. Estructura del Response

### Éxito (HTTP 200)
```json
{
  "StatusCode": 200,
  "Mensaje": "Application injected successfully",
  "NumeroSolicitud": 192269
}
```

### Error de Validación HTTP (HTTP 400)
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.Participantes[0].Rut' Error:Field validation for 'Rut' failed on the 'customRutPattern' tag"
}
```

### Error de Negocio / FinnFlow (HTTP 500)
```json
{
  "StatusCode": 500,
  "Mensaje": "FinnFlow injection failed: 406 NOT_ACCEPTABLE - {\"code\":\"invalid_value\",\"message\":\"Plazo1 with invalid value\"}",
  "NumeroSolicitud": null
}
```

### Token Inválido (HTTP 401)
```json
{
  "StatusCode": 401,
  "Mensaje": "Token is not valid",
  "NumeroSolicitud": null
}
```

---

## 4. Casos de Uso - Capa de Validación HTTP

> Estos errores los genera el framework Go **antes** de llegar al negocio. El legado Java hacía lo mismo vía `@Valid` de Spring.  
> **Resultado esperado: HTTP 400**

---

### TC-01: Body vacío o mal formado

**Descripción:** Se envía un JSON completamente vacío o texto plano.

**Request:**
```http
POST /v1/bfcl/mortgage-loan/injections
Content-Type: application/json

{}
```

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.DatosCredito' Error:Field validation for 'DatosCredito' failed on the 'required' tag"
}
```

---

### TC-02: Falta el bloque `datos_credito`

**Request:**
```json
{
  "participantes": [
    {
      "Rut": "18.015.051-K",
      "TipoParticipacion": 1,
      "Nombre": "Pedro",
      "Paterno": "Perez",
      "Materno": "Perez",
      "FechaNacimiento": "04-07-1992"
    }
  ]
}
```

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.DatosCredito' Error:Field validation for 'DatosCredito' failed on the 'required' tag"
}
```

---

### TC-03: Falta el bloque `participantes`

**Request:**
```json
{
  "datos_credito": {
    "NumeroSolicitud": 192269,
    "AntiguedadVivienda": 1,
    "EjecutivoComercial": "PSOTO",
    "Producto": 1,
    "Objetivo": 1,
    "Destino": 1,
    "FechaAprobacion": "2024-01-15",
    "MontoAprobado": 1688.09,
    "ValorPropiedad": 2300.00,
    "ValorContado": 261.91,
    "Plazo1": 30,
    "MesesGracia": 0,
    "Tasa1": 3.99,
    "Spread1": 0.0
  }
}
```

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.Participantes' Error:Field validation for 'Participantes' failed on the 'required' tag"
}
```

---

### TC-04: Array `participantes` vacío `[]`

**Request:** mismo de TC-03 pero con:
```json
"participantes": []
```

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.Participantes' Error:Field validation for 'Participantes' failed on the 'min' tag"
}
```

---

### TC-05: Rut con formato inválido

**Descripción:** El Rut debe cumplir el patrón `XX.XXX.XXX-X` (guiones y puntos obligatorios).

**Casos inválidos:**
- `"18015051K"` (sin puntos ni guión)
- `"18015051-K"` (sin puntos)
- `"18.015.051K"` (sin guión)
- `"1234-5"` (muy corto)

**Request (solo cambiar el Rut en un payload completo):**
```json
"Rut": "18015051K"
```

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.Participantes[0].Rut' Error:Field validation for 'Rut' failed on the 'customRutPattern' tag"
}
```

---

### TC-06: Fecha de Nacimiento con formato inválido

**Descripción:** La fecha debe estar en formato `DD-MM-YYYY`.

**Casos inválidos:**
- `"1992-07-04"` (formato ISO, incorrecto)
- `"04/07/1992"` (barras en lugar de guiones)
- `"4-7-92"` (dígitos sin ceros)

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.Participantes[0].FechaNacimiento' Error:Field validation for 'FechaNacimiento' failed on the 'customFechaPattern' tag"
}
```

---

### TC-07: `MontoAprobado` en cero o negativo

**Descripción:** El monto debe ser estrictamente mayor a 0 (`gt=0`).

**Request:** `"MontoAprobado": 0` o `"MontoAprobado": -100`

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.DatosCredito.MontoAprobado' Error:Field validation for 'MontoAprobado' failed on the 'gt' tag"
}
```

---

### TC-08: `EjecutivoComercial` supera 20 caracteres

**Request:** `"EjecutivoComercial": "NOMBREMUYLARGOPARAELCAMPO"`

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.DatosCredito.EjecutivoComercial' Error:Field validation for 'EjecutivoComercial' failed on the 'max' tag"
}
```

---

### TC-09: `Email` del participante con formato inválido

**Descripción:** El Email es opcional, pero si se envía, debe ser un email válido.

**Request:** `"Email": "esto-no-es-un-email"`

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.Participantes[0].Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

---

## 5. Casos de Uso - Reglas de Negocio

> Estas validaciones las ejecuta el **Caso de Uso** (capa de dominio). El legado Java las tenía en la clase `InjectionService`.  
> **Resultado esperado: HTTP 500** (Paridad con el legado)

---

### TC-10: `MontoAprobado` mayor que `ValorPropiedad`

**Descripción:** El monto a prestar no puede superar el valor del bien hipotecado.

**Request:**
```json
"MontoAprobado": 3000.00,
"ValorPropiedad": 2000.00
```

**Response esperado (HTTP 500):**
```json
{
  "StatusCode": 500,
  "Mensaje": "MontoAprobado with invalid value",
  "NumeroSolicitud": null
}
```

---

### TC-11: `Plazo1` fuera del rango permitido

**Descripción:** Solo se aceptan plazos de 5, 10, 15, 20, 25 o 30 años.

**Casos inválidos:** `3`, `7`, `12`, `17`, `22`, `35`

**Request:** `"Plazo1": 17`

**Response esperado (HTTP 500):**
```json
{
  "StatusCode": 500,
  "Mensaje": "Plazo1 with invalid value",
  "NumeroSolicitud": null
}
```

**Casos válidos:** `5`, `10`, `15`, `20`, `25`, `30` ✅

---

### TC-12: `Tasa1` fuera del rango permitido

**Descripción:** La tasa debe estar entre 0.1% y 50%.

**Casos inválidos:** `0.05`, `0.0`, `-1.0`, `51.0`, `100.0`

**Request:** `"Tasa1": 0.05`

**Response esperado (HTTP 500):**
```json
{
  "StatusCode": 500,
  "Mensaje": "Tasa1 with invalid value",
  "NumeroSolicitud": null
}
```

---

## 6. Casos de Uso - Integración FinnFlow

> Estos casos se prueban conectados a QA real (o con el Mock ampliado). Simulan las respuestas de error de FinnFlow que el legado Java retransmitía al cliente.

---

### TC-13: Happy Path - Inyección exitosa

**Request (payload mínimo válido):**
```json
{
  "datos_credito": {
    "NumeroSolicitud": 192269,
    "AntiguedadVivienda": 1,
    "EjecutivoComercial": "PSOTO",
    "Producto": 1,
    "Objetivo": 1,
    "Destino": 1,
    "FechaAprobacion": "2024-01-15",
    "MontoAprobado": 1688.09,
    "ValorPropiedad": 2300.00,
    "ValorContado": 261.91,
    "Plazo1": 30,
    "MesesGracia": 0,
    "Tasa1": 3.99,
    "Spread1": 0.0
  },
  "participantes": [
    {
      "Rut": "18.015.051-K",
      "TipoParticipacion": 1,
      "Nombre": "Pedro",
      "Paterno": "Perez",
      "Materno": "Perez",
      "FechaNacimiento": "04-07-1992"
    }
  ]
}
```

**Response esperado (HTTP 200):**
```json
{
  "StatusCode": 200,
  "Mensaje": "Application injected successfully",
  "NumeroSolicitud": 192269
}
```

---

### TC-14: FinnFlow rechaza la inyección (Error 406)

**Descripción:** FinnFlow devuelve `NOT_ACCEPTABLE`. El legado re-empaquetaba este error en un 500. Nuestra app Go replica eso.

**Response esperado (HTTP 500):**
```json
{
  "StatusCode": 500,
  "Mensaje": "FinnFlow injection failed: 406 NOT_ACCEPTABLE - {\"code\":\"invalid_value\",\"message\":\"...\"}",
  "NumeroSolicitud": null
}
```

---

### TC-15: FinnFlow rechaza el formato (Error 402)

**Response esperado (HTTP 500):**
```json
{
  "StatusCode": 500,
  "Mensaje": "FinnFlow injection failed: 402 PAYMENT_REQUIRED - {\"code\":\"invalid_format\",\"message\":\"...\"}",
  "NumeroSolicitud": null
}
```

---

### TC-16: Credenciales FinnFlow inválidas (Token inválido)

**Descripción:** Las credenciales `FINNFLOW_KEY`/`FINNFLOW_SECRET` son incorrectas. FinnFlow devuelve `token_not_valid`.

**Response esperado (HTTP 401):**
```json
{
  "StatusCode": 401,
  "Mensaje": "Token is not valid",
  "NumeroSolicitud": null
}
```

---

## 7. Casos de Uso - Campos Opcionales

---

### TC-17: Inyección de Renegociación (Con `NumeroOperacionOriginal`)

**Descripción:** Para operaciones de renegociación se envía el número de operación original.

**Request adicional en `datos_credito`:**
```json
"NumeroOperacionOriginal": "1234567890123456"
```

**Response esperado (HTTP 200):** igual que TC-13.

---

### TC-18: `NumeroOperacionOriginal` supera 16 caracteres

**Request:** `"NumeroOperacionOriginal": "12345678901234567"` (17 chars)

**Response esperado (HTTP 400):**
```json
{
  "code": "invalid_format",
  "message": "Key: 'InyeccionRequestDTO.DatosCredito.NumeroOperacionOriginal' Error:Field validation for 'NumeroOperacionOriginal' failed on the 'max' tag"
}
```

---

### TC-19: Happy Path con múltiples participantes

**Descripción:** La aplicación acepta múltiples participantes (deudor + codeudor, etc.).

**Request:** agregar un segundo objeto en el array `participantes` con Rut diferente.

**Response esperado (HTTP 200):** igual que TC-13.

---

### TC-20: Happy Path con bloque de `propiedades`

**Request:** agregar el bloque `propiedades` al payload de TC-13.
```json
"propiedades": [
  {
    "TipoInmueble": 1,
    "Antiguedad": 10,
    "Direccion": "Av. Falabella 123",
    "Numero": 1234,
    "Depto": "201",
    "Comuna": 13101
  }
]
```

**Response esperado (HTTP 200):** igual que TC-13.

---

### TC-21: Campos opcionales ausentes (Backward Compatibility)

**Descripción:** Clientes que omiten `Email`, `Depto`, `NumeroOperacionOriginal` y otros campos opcionales no deben ser afectados. Los campos `null` no se envían hacia FinnFlow.

**Response esperado (HTTP 200):** igual que TC-13 ✅

---

## 8. Tabla Resumen de Paridad

| # | Nombre del Caso | HTTP Legado | HTTP Go | Paridad |
|---|---|:---:|:---:|:---:|
| TC-01 | Body vacío | 400 | 400 | ✅ |
| TC-02 | Sin `datos_credito` | 400 | 400 | ✅ |
| TC-03 | Sin `participantes` | 400 | 400 | ✅ |
| TC-04 | `participantes: []` | 400 | 400 | ✅ |
| TC-05 | Rut inválido | 400 | 400 | ✅ |
| TC-06 | Fecha nacimiento inválida | 400 | 400 | ✅ |
| TC-07 | MontoAprobado en cero | 400 | 400 | ✅ |
| TC-08 | EjecutivoComercial > 20 chars | 400 | 400 | ✅ |
| TC-09 | Email inválido | 400 | 400 | ✅ |
| TC-10 | Monto > Valor Propiedad | 500 | 500 | ✅ |
| TC-11 | Plazo inválido | 500 | 500 | ✅ |
| TC-12 | Tasa inválida | 500 | 500 | ✅ |
| TC-13 | Happy Path mínimo | 200 | 200 | ✅ |
| TC-14 | FinnFlow rechaza (406) | 500 | 500 | ✅ |
| TC-15 | FinnFlow rechaza (402) | 500 | 500 | ✅ |
| TC-16 | Credenciales inválidas | 401 | 401 | ✅ |
| TC-17 | Renegociación válida | 200 | 200 | ✅ |
| TC-18 | NumeroOperacion > 16 chars | 400 | 400 | ✅ |
| TC-19 | Múltiples participantes | 200 | 200 | ✅ |
| TC-20 | Con bloque propiedades | 200 | 200 | ✅ |
| TC-21 | Campos opcionales ausentes | 200 | 200 | ✅ |

---

## 9. Checklist Final

### Ambiente Local Mock ✅
- [ ] Mock corre en `localhost:8083`
- [ ] API corre en `localhost:8081`
- [ ] `GET /health` devuelve `{"status":"UP"}` (HTTP 200)
- [ ] TC-13 devuelve HTTP 200 con Mock

### Validaciones HTTP (TC-01 a TC-09)
- [ ] TC-01 - Body vacío → 400
- [ ] TC-02 - Sin datos_credito → 400
- [ ] TC-03 - Sin participantes → 400
- [ ] TC-04 - Participantes vacío → 400
- [ ] TC-05 - Rut inválido → 400
- [ ] TC-06 - Fecha nacimiento inválida → 400
- [ ] TC-07 - Monto en cero → 400
- [ ] TC-08 - Ejecutivo > 20 chars → 400
- [ ] TC-09 - Email inválido → 400

### Reglas de Negocio (TC-10 a TC-12)
- [ ] TC-10 - Monto > Propiedad → 500
- [ ] TC-11 - Plazo inválido → 500
- [ ] TC-12 - Tasa inválida → 500

### Integración FinnFlow QA Real (TC-13 a TC-16)
- [ ] TC-13 - Happy Path → 200
- [ ] TC-16 - Credenciales inválidas → 401

### Opcionales y Compatibilidad (TC-17 a TC-21)
- [ ] TC-17 - Renegociación → 200
- [ ] TC-18 - NumeroOperacion largo → 400
- [ ] TC-19 - Múltiples participantes → 200
- [ ] TC-20 - Con propiedades → 200
- [ ] TC-21 - Campos opcionales ausentes → 200
