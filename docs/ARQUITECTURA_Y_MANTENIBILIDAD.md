# Arquitectura y Guía de Mantenibilidad
**Microservicio:** `mortgage-loan-injection`
**Stack:** Go (Golang), Gin Gonic, Log/Slog, Datadog APM
**Patrón Arquitectónico:** Arquitectura Hexagonal (Ports and Adapters)

---

## 1. Visión Arquitectónica

Se ha migrado el microservicio legado (desarrollado en Java/Spring Boot) hacia una plataforma de alto rendimiento en Go. El objetivo principal ha sido asegurar **paridad extrema (100% retrocompatibilidad)** con los consumidores actuales, al mismo tiempo que se reducen agresivamente el footprint de memoria y los tiempos de arranque.

Se adoptó la **Arquitectura Hexagonal** para asegurar que la lógica de negocio (Dominio) sea el núcleo agnóstico del sistema, permitiendo intercambiar infraestructuras (Frameworks HTTP, CRMs, APIs externas) sin afectar las reglas hipotecarias.

### 1.1 Estructura de Capas
El proyecto se divide de acuerdo al patrón de separación de responsabilidades:

- `/internal/core/domain`: Estructuras puras de negocio. Agnósticas a JSON, bases de datos o HTTP.
- `/internal/core/usecases`: Validaciones complejas y flujos de negocio (Ports).
- `/internal/app/api`: Capa de entrada (Inbound Adapters). Contiene los controladores HTTP (Gin), Middlewares y DTOs con sus tags de validación estricta (`validator/v10`).
- `/internal/infra`: Capa de salida (Outbound Adapters). Implementaciones técnicas, como el cliente REST para conectarse con FinnFlow/Salesforce.

---

## 2. Atributos de Calidad

### 2.1 Escalabilidad y Rendimiento
El sistema aprovecha la compilación binaria de Go. No requiere calentar una máquina virtual (JVM).
* **Consumo de Memoria:** Se reduce de ~300MB (Spring Boot) a un estimado de **~15MB**.
* **Arranque:** Menor a 5ms.
* **Manejo de Concurrencia:** Cada petición HTTP es despachada en su propia `goroutine`, maximizando la escalabilidad horizontal en Kubernetes y minimizando la saturación del CPU durante operaciones I/O (ej. llamando a FinnFlow).

### 2.2 Agnosticismo
El framework `Gin` y la API de terceros `FinnFlow` son vistos como simples plugins.
Si en el futuro Banco Falabella decide abandonar Salesforce, solo se debe construir un nuevo archivo en `internal/infra/client` que implemente la interfaz `FinnFlowClient`. El resto del servicio (`Handlers` y `UseCases`) permanecerá inalterado.

### 2.3 12-Factor y Observabilidad
- **Configuración:** Cero variables harcodeadas. Se exige la inyección de `FINNFLOW_URL`, `FINNFLOW_KEY`, `FINNFLOW_SECRET`. El sistema ejecutará un `os.Exit(1)` (Fail-Fast) si las variables maestras no están presentes en el booteo.
- **Trazabilidad:** Integración nativa con `dd-trace-go`. Las métricas de latencia de APM y los `TraceIDs` se inyectan automáticamente en cada Request.
- **Logs Compatibles:** El paquete `log/slog` fue sobreescrito con `ReplaceAttr` para emular el `Logback JSON Layout` de Java, inyectando los atributos de clase y thread (`"logger"`, `"thread"`) para que ElasticSearch no detecte la transición.

---

## 3. Manejo de Paridad y "Bugs" Legados
Durante la migración se descubrieron comportamientos asíncronos en el casteo de errores de Java:
- Excepciones del cliente HTTP de Spring fallaban al parsear y terminaban enviando un HTTP `500 Internal Server Error` con el payload de validación adentro.
- **Acuerdo de Arquitectura:** En Go, los handlers capturan los fallos de dominio (`402`, `406`, `400`) y **emulan intencionalmente la caída 500** devolviendo la frase exacta `FinnFlow injection failed...` para no romper los clientes actuales que ya manejan la respuesta de esta peculiar forma.

---

## 4. Guía Operativa: ¿Cómo agregar un nuevo campo al JSON?

Ante la necesidad del negocio de incorporar una nueva variable (Ej: `"TasaSeguro"`), el equipo de desarrollo debe seguir obligatoriamente este pipeline de 4 pasos (De afuera hacia adentro, y hacia afuera nuevamente):

### Paso 1: Actualizar el Contrato (Capa DTO)
Modificar `internal/app/api/dto/injection.go`. Añadir el campo y asignar la validación declarativa (ej. mayor a cero).
```go
TasaSeguro *float64 `json:"TasaSeguro" binding:"required,gt=0"`
```
*(Cualquier violación a esta regla causará un rechazo automático sin necesidad de escribir validadores manuales).*

### Paso 2: Actualizar el Modelo Puro (Capa Domain)
Modificar `internal/core/domain/injection.go`. Agregar el campo crudo sin tags técnicos.
```go
TasaSeguro float64
```

### Paso 3: Traducir los datos (Capa Mapper)
Modificar `internal/app/api/mapper/injection_mapper.go`. Realizar el traspaso seguro desde el DTO hacia la entidad de dominio.
```go
TasaSeguro: *dtoReq.TasaSeguro,
```

### Paso 4: Lógica Matemática (Capa UseCase - Opcional)
Si la variable tiene dependencias cruzadas (Ej: TasaSeguro no puede superar Tasa1), aplicar la regla en `internal/core/usecases/injection_usecase.go` dentro del método `validateBusinessRules`.

---

## 5. Estrategia de Testing (Certificación Híbrida)
Todo cambio en el pipeline descrito arriba debe ser avalado por la suite de pruebas.
1. **Unit/Integration Tests (`go test ./...`):** Certifican que las inyecciones de dependencias (Mocks) y la aserción de estados funcionen.
2. **QA Parity Extrema (`powershell.exe -ExecutionPolicy Bypass -File .\test_local.ps1`):** Script de PowerShell que limpia los puertos y posee inteligencia para orquestar el monolito Java de forma híbrida: si Docker Desktop no está encendido, hace fallback usando `mvn spring-boot:run` nativo. Luego levanta el binario de Go e inyecta payloads a un Simulador Mock de pruebas. Certifica que la respuesta de Go sea **idéntica byte a byte** a la respuesta histórica de la aplicación Java y aborta con error si no hay conexión real.
