# Definiciones de Dominio y Contratos JSON (Legado)

Este documento contiene las definiciones de entidades y payloads JSON esperados y retornados por el sistema legado de Java, para garantizar paridad exacta en la migración a Go.

## 1. Solicitud de Inyección (InyeccionRequest)

Estructura raíz del JSON de petición:
```json
{
  "datos_credito": { ... },
  "participantes": [ { ... } ],
  "propiedades": [ { ... } ]
}
```

### 1.1 DatosCredito
Campos esperados dentro de `datos_credito`:
- `NumeroSolicitud` (Integer, obligatorio)
- `AntiguedadVivienda` (Integer, obligatorio)
- `EjecutivoComercial` (String, máx 20 chars, obligatorio)
- `Producto` (Integer, obligatorio)
- `Objetivo` (Integer, obligatorio)
- `Destino` (Integer, obligatorio)
- `MontoAprobado` (Float/Double, min 0.0001, obligatorio)
- `ValorPropiedad` (Float/Double, min 0.01, obligatorio)
- `FechaAprobacion` (String, obligatorio)
- `ValorContado` (Float/Double, min 0.01, obligatorio)
- `Plazo1` (Integer, obligatorio)
- `Plazo2` (Integer, opcional)
- `MesesGracia` (Integer, obligatorio)
- `Tasa1` (Float/Double, min 0.01, obligatorio)
- `Spread1` (Float/Double, min 0.0, obligatorio)

### 1.2 Participante
Campos esperados dentro de la lista `participantes` (mínimo 1):
- `Rut` (String, obligatorio, patrón regex: `^[0-9]{1,2}\.[0-9]{3}\.[0-9]{3}-[0-9kK]$`)
- `TipoParticipacion` (Integer, obligatorio)
- `Nombre` (String, máx 100 chars, obligatorio)
- `Paterno` (String, máx 100 chars, obligatorio)
- `Materno` (String, máx 100 chars, obligatorio)
- `FechaNacimiento` (String, obligatorio, patrón regex: `^[0-9]{2}-[0-9]{2}-[0-9]{4}$`)
- `Email` (String, formato email, máx 100 chars, opcional)

### 1.3 Propiedad
Campos esperados dentro de la lista `propiedades` (opcional):
- `TipoInmueble` (Integer, obligatorio)
- `Antiguedad` (Integer, obligatorio)
- `Direccion` (String, máx 100 chars, obligatorio)
- `Numero` (Integer, obligatorio)
- `Depto` (String, máx 20 chars, opcional)
- `Comuna` (Integer, obligatorio)

## 2. Respuesta de Inyección (InyeccionResponse)

Estructura del JSON de respuesta en casos de éxito y algunos errores reportados por FinnFlow:
```json
{
  "StatusCode": 200,
  "Mensaje": "Mensaje descriptivo",
  "NumeroSolicitud": 12345
}
```
Campos de la respuesta:
- `StatusCode` (Integer)
- `Mensaje` (String)
- `NumeroSolicitud` (Integer)
