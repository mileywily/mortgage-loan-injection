# Mortgage Loan Injection (Go Hexagonal)

Bienvenido al repositorio oficial del microservicio **Mortgage Loan Injection**. 
Este proyecto es el resultado de la migración del antiguo monolito Java/Spring Boot a **Golang 1.23** utilizando Arquitectura Hexagonal, logrando un aumento masivo de rendimiento y reduciendo el consumo de memoria en un 90%, **manteniendo una paridad de comportamiento extrema (bug-for-bug compatibility)**.

## 📚 Portal de Documentación Oficial

Como parte de la entrega arquitectónica, la documentación ha sido dividida por roles para facilitar su consumo:

### 1. Equipo de Desarrollo & Arquitectura
Diseñado para los desarrolladores que mantendrán o extenderán el código en el futuro.
* 🏛️ [Arquitectura y Mantenibilidad](docs/ARQUITECTURA_Y_MANTENIBILIDAD.md): Decisiones de diseño, Clean Architecture, gestión de errores, observabilidad (Datadog) y 12-Factor App.
* 🧠 [Definiciones de Dominio](docs/DEFINICIONES_DOMINIO.md): Diccionario de datos y reglas de negocio del dominio Hipotecario.
* 📜 [Casos de Uso del Legado](docs/CASOS_DE_USO_LEGADO.md): Especificación técnica de cómo el legado trataba los endpoints.

### 2. Equipo de Aseguramiento de Calidad (QA)
Diseñado para la certificación del software y el pase a producción sin impacto.
* 🧪 [Guía de Certificación QA](docs/GUIA_CERTIFICACION_QA.md): Matriz de pruebas obligatorias, escenarios de fallo y explicación de la compatibilidad retroactiva de errores (Status 500).
* 🗃️ [Colección de Postman](docs/postman_collection.json): Archivo exportado con los 4 flujos principales para ejecución manual o automatizada en pipelines (Happy Path, Missing Fields, Invalid Rut, Invalid Plazo).
* 📝 [Payloads Históricos](docs/PAYLOADS_CASOS_DE_USO.md): Ejemplos en crudo de las peticiones históricas.

### 3. Operaciones & DevOps
* El servicio requiere obligatoriamente la inyección de variables de entorno para su arranque (`FINNFLOW_URL`, `FINNFLOW_KEY`, `FINNFLOW_SECRET`, `PORT`).
* Las trazas de Log JSON están formateadas para ser ingeridas exactamente igual que la antigua librería de `logback`.
* Métricas y APM están integradas mediante `dd-trace-go` de Datadog de forma automática.

## 🚀 Cómo ejecutar localmente

```bash
# Compilar y ejecutar
go build -o api.exe cmd/api/main.go
$env:PORT="8082"
$env:FINNFLOW_URL="http://url-de-tu-ambiente:8083"
$env:FINNFLOW_KEY="admin"
$env:FINNFLOW_SECRET="admin123"
./api.exe
```

## 🛠 Herramienta de Certificación Automática
Se adjunta el script `test_local.ps1` que orquesta de forma local tu imagen Java existente y la nueva aplicación Go, enviándoles peticiones en paralelo a un mock y asegurando que las respuestas **coincidan byte a byte**.

```powershell
powershell.exe -ExecutionPolicy Bypass -File .\test_local.ps1
```
