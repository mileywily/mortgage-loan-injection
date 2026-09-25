# Guía Detallada: Despliegue en Nullplatform

## ¿Qué es Nullplatform?
Nullplatform es un Portal de Desarrolladores (Internal Developer Portal / PaaS) que abstrae la complejidad de la infraestructura en la nube (como Kubernetes, AWS, etc.). 
En lugar de escribir complejos manifiestos, te permite definir tu aplicación de forma declarativa. Nullplatform se encarga de crear la infraestructura, manejar los balanceadores de carga, y aprovisionar el despliegue de tus contenedores Docker de forma estandarizada.

---

## 1. Identidad de la Aplicación
El primer aspecto crítico es la identidad. En Nullplatform, cada microservicio debe tener un nombre único que lo vincule con su infraestructura.
- **Nombre de la Aplicación:** `api-bfcl-mortgage-loans-injection`
- **Lenguaje/Framework:** Go (Golang) 1.23
- **Tipo de Carga de Trabajo:** Servicio Web / API REST (Stateless)

---

## 2. Archivo de Configuración (`nullplatform.yml`)
Este es el manifiesto principal que vivirá en la raíz de nuestro repositorio. Nullplatform lo leerá para saber cómo configurar el contenedor. Deberá contener los siguientes aspectos necesarios:

*   **Aplicación (app):** `api-bfcl-mortgage-loans-injection`
*   **Networking (Puertos):** Expondremos el puerto `8081` (el mismo que usa nuestro `Dockerfile` actual).
*   **Healthchecks (Sondas de Salud):**
    *   `livenessProbe`: Para que Nullplatform sepa si el contenedor se bloqueó y necesita reiniciarlo.
    *   `readinessProbe`: Para saber cuándo el servicio en Go está listo para recibir tráfico (usualmente chequeando el puerto TCP o un endpoint `/health`).
*   **Recursos (Compute):** Se asignan límites de Memoria y CPU recomendados para Go (ej. CPU: 100m - 500m, RAM: 128Mi - 256Mi).

---

## 3. Matriz de Variables de Entorno (Parameters & Secrets)
Nullplatform elimina el uso de archivos `.env`. Las variables se deben dar de alta directamente en la consola/UI de Nullplatform o vía CLI bajo el "Scope" de la aplicación. 
Para que el microservicio funcione con **Paridad Extrema**, debes solicitar a los administradores de Nullplatform (o configurarlas tú mismo) las siguientes variables por cada entorno (Dev, QA, Prod):

| Variable / Parámetro | Tipo | Descripción |
| :--- | :--- | :--- |
| `PORT` | Plana | Puerto donde escuchará Go (Obligatorio: `8081`). |
| `FINNFLOW_URL` | Plana | URL del core ASICOM/FinnFlow correspondiente al ambiente (ej. URL del Mock en QA, URL real en Prod). |
| `FINNFLOW_KEY` | Secreto | Credencial (API Key / Client ID) para integrarse al Legacy. |
| `FINNFLOW_SECRET` | Secreto | Contraseña o Hash secreto para la autenticación en FinnFlow. |
| `GIN_MODE` | Plana | Debe ser `release` para ambientes de QA y Producción. |

> ⚠️ **Atención:** En nuestro código en Go, usamos `os.Getenv()`. Nullplatform inyectará estas llaves directamente al sistema operativo del contenedor, por lo que el código no requiere ninguna alteración.

---

## 4. Integración CI/CD (GitHub Actions)
Nuestro pipeline `.github/workflows/ci.yml` sufrirá una reestructuración en su etapa final. Los aspectos necesarios a incorporar son:

1. **Autenticación (Auth):**
   - GitHub Actions necesita credenciales (`NULLPLATFORM_API_KEY`) guardadas en los *GitHub Secrets* para poder hablar con la API de Nullplatform.
2. **Registro de Build (Build Push):**
   - El pipeline compilará la imagen Docker.
   - En lugar de subirla a cualquier lado, usará los GitHub Actions oficiales de Nullplatform (ej. `nullplatform/github-action-build@v1`) para notificar que la versión con el commit actual está lista.
3. **Solicitud de Release (Despliegue):**
   - Se añadirá un paso (ej. `nullplatform/github-action-release@v1`) para disparar automáticamente el despliegue hacia un entorno (ej. `dev` o `qa`) apenas el *Build* termine exitosamente.

---

## 5. Observabilidad (Logs y APM)
- **Logs:** Nullplatform captura la salida estándar (`stdout`) del contenedor por defecto. Como configuramos nuestro servicio en Go con logs estructurados en JSON (`slog`), Nullplatform (y Datadog/ELK por debajo) los parseará automáticamente sin esfuerzo adicional.
- **Trazabilidad (Datadog):** Si Nullplatform tiene integración automática con Datadog APM, inyectará las variables `DD_ENV`, `DD_SERVICE` y `DD_VERSION`. El middleware `gintrace` que ya incluimos en el código tomará estas variables automáticamente.

---

## Resumen del Plan de Ejecución Técnico
Una vez validada esta guía, aplicaremos los siguientes cambios reales en la rama:
1. Crear el manifiesto `nullplatform.yml` con la topología de la app.
2. Modificar `.github/workflows/ci.yml` para agregar los jobs de Nullplatform (Build & Release).
3. Asegurar que el `Dockerfile` cumpla con los estándares esperados.
