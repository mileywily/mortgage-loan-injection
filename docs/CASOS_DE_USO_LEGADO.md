# Casos de Uso y Matriz de Errores Legado

## Endpoints

- **POST /v1/bfcl/mortgage-loan/injections**
  - **Descripción**: Inyecta una solicitud de crédito hipotecario en Salesforce/Finnflow.

## Validaciones de Negocio (InjectionService.java)

1. `Participantes`: no debe ser nulo o vacío. Error: `IllegalArgumentException("Missing participantes field")` -> HTTP 400.
2. `MontoAprobado` no debe ser mayor a `ValorPropiedad`. Error: `IllegalArgumentException("MontoAprobado with invalid value")` -> HTTP 406.
3. `Plazo1` debe estar entre 5 y 30 años, y ser múltiplo de 5. Error: `IllegalArgumentException("Plazo1 with invalid value")` -> HTTP 406.
4. `Tasa1` debe estar entre 0.1 y 50.0. Error: `IllegalArgumentException("Tasa1 with invalid value")` -> HTTP 406.

## Matriz de Errores (GlobalExceptionHandler.java)

El sistema legacy mapea los errores de la siguiente forma no estándar:

| Excepción / Condición | Status Code | Error Code (JSON) | Mensaje (JSON) |
| :--- | :--- | :--- | :--- |
| `AccessDeniedException` | 401 UNAUTHORIZED | `token_not_valid` | "Token is not valid" |
| `HttpMessageNotReadableException` (Sin Body) | 400 BAD REQUEST | `missing_field` | "Request body is required" |
| Campo inválido - `invalid format` | 402 PAYMENT REQUIRED | `invalid_format` | "{fieldName} with invalid format" |
| Campo inválido - `invalid value` | 406 NOT ACCEPTABLE | `invalid_value` | "{fieldName} with invalid value" |
| Falla validación - Otro (ej. nulo) | 400 BAD REQUEST | `missing_field` | "Missing {fieldName} field" |
| Error Interno (`Exception`) | 500 INTERNAL SERVER ERROR| `internal_error` | "Internal server error" |

## Respuestas del Servicio a Terceros (ThirdPartyApiService)

Respuestas directas del método `processInjection` que extrae errores de FinnFlow:
- Token inválido (`token_not_valid`, `Token is invalid`): "Token de autenticación no válido para el sistema FinnFlow."
- Error Auth (`Authentication failed`, `401`): "Error de autenticación: credenciales inválidas..." o "Error de autenticación con el sistema..."
- Data inválida (`400`): "Datos de solicitud no válidos para el sistema FinnFlow."
- Acceso denegado (`403`): "Acceso denegado al sistema FinnFlow."
- No encontrado (`404`): "Servicio no encontrado en el sistema FinnFlow."
- Timeout / Conectividad (`timeout`, `connect`): "Timeout en la comunicación con FinnFlow."
- SSL / Cert (`certificate`, `SSL`, `PKIX`): "Error de conectividad SSL con FinnFlow."
- Por defecto: "Error interno del sistema FinnFlow."
