# Guía de Certificación QA - Migración Hipotecaria (Go)

Este documento contiene los lineamientos, matrices de prueba y herramientas automatizadas para que el equipo de Aseguramiento de Calidad (QA) certifique el nuevo microservicio hipotecario escrito en Go, garantizando **cero impacto** frente a los consumidores actuales.

---

## 1. Alcance de la Certificación

El objetivo de QA **no es** probar nuevas funcionalidades, sino certificar una **Paridad Extrema (Bug-for-Bug Compatibility)**. Esto significa que la nueva API en Go debe comportarse de forma **idéntica** a la API de Java frente a escenarios exitosos y frente a errores.

**Puntos críticos a certificar:**
1. Mapeo de Request/Response JSON (Nombres de atributos y tipos de datos idénticos).
2. Manejo explícito de valores nulos (Ej: `"NumeroSolicitud": null`).
3. Status Codes HTTP (Atención especial a validaciones defectuosas del legado que retornaban HTTP 500).

---

## 2. Matriz de Casos de Uso a Certificar

QA debe validar los siguientes escenarios obligatorios:

| ID | Escenario | Condición de Entrada (Payload Modificado) | Resultado Esperado (Status Code) | Mensaje Esperado (Body JSON) |
|---|---|---|---|---|
| **QA-01** | Happy Path | Payload completo y válido. | `200 OK` | `"Mensaje": "Application injected successfully"` |
| **QA-02** | Campo Requerido Faltante | Se omite un campo clave (ej. `NumeroSolicitud` o el arreglo de `participantes` viene vacío). | `500 Internal Server Error` | `"Mensaje": "FinnFlow injection failed: 400 BAD_REQUEST - {\"code\":\"missing_field\",...}"` |
| **QA-03** | Formato de RUT/Email Inválido | `"Rut": "rut-sin-formato"` o `"Email": "correo-malo"`. | `500 Internal Server Error` | `"Mensaje": "FinnFlow injection failed: 402 PAYMENT_REQUIRED - {\"code\":\"invalid_format\",...}"` |
| **QA-04** | Regla de Negocio (Plazo) | `"Plazo1": 17` (El plazo no es múltiplo de 5, menor a 5 o mayor a 30). | `500 Internal Server Error` | `"Mensaje": "FinnFlow injection failed: 406 NOT_ACCEPTABLE - {\"code\":\"invalid_value\",...}"` |
| **QA-05** | Regla de Negocio (Montos) | `"MontoAprobado"` supera al `"ValorPropiedad"`. | `500 Internal Server Error` | `"Mensaje": "FinnFlow injection failed: 406 NOT_ACCEPTABLE - {\"code\":\"invalid_value\",...}"` |

> ⚠️ **Nota de Arquitectura sobre el Error 500:** El equipo de QA notará que los errores de validación de formulario retornan `500` en lugar de `400`. Esto es **esperado** y fue replicado intencionalmente de un fallo del sistema de Java para no quebrar las integraciones de los frontends actuales.

---

## 3. Herramientas de Prueba

QA tiene a su disposición dos vías para ejecutar esta certificación localmente:

### Vía 1: Automatizada Híbrida (Recomendado)
Se ha construido un script de automatización (`test_local.ps1`) que posee inteligencia de entorno. Levanta de forma simultánea el motor de Java (con Docker o Maven Nativo), el de Go y un Simulador Mock de Finnflow.

**Paso a paso:**
1. Abrir **PowerShell** como administrador (para permitir la gestión y liberación de puertos atascados).
2. Navegar a la carpeta del nuevo proyecto: `cd C:\HEXAGONAL-INJECTION\mortgage-loan-injection`
3. Ejecutar el script eludiendo la restricción de seguridad: 
   ```powershell
   powershell.exe -ExecutionPolicy Bypass -File .\test_local.ps1
   ```
4. **Validación Visual:** El script te avisará que usará la imagen preconstruida `java-app-injection` y clonará su estado con `docker commit` para inyectar inteligentemente las variables del Mock. 
5. El sistema de *Polling HTTP* esperará el tiempo que sea necesario (máximo 3 minutos) hasta que Java y Go estén sanos.
6. El QA script informará en consola si los bytes exactos coinciden entre el Legado y Go.
7. QA debe tomar captura de pantalla del mensaje final verde: `TODAS LAS PRUEBAS DE PARIDAD EXTREMA PASARON CON EXITO`.

### Vía 2: Pruebas Exploratorias Manuales (Postman)
1. Importar en Postman el archivo: `docs/postman_collection.json`.
2. Levantar el servicio Go mediante `docker-compose up -d go-app` o `go run cmd/api/main.go` (puerto 8082).
3. Ejecutar la carpeta *"Mortgage Loan Injection (Go & Legacy)"*.
4. Alterar manualmente los JSON (ej. colocar `Tasa1: 60.0`) y validar que la respuesta replique siempre la matriz del apartado 2.

---

## 4. Firmas de Aprobación
Una vez ejecutados los casos de uso, QA debe validar las trazas en APM (Datadog) para asegurar que el `TraceID` viaje correctamente. Al certificar todo, el ticket quedará habilitado para el pase a Producción.
