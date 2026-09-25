# Tutorial Paso a Paso: Interfaz de Nullplatform (UI)

Como es tu primera vez usando Nullplatform, he diseñado este tutorial como una guía "clic a clic". Sigue estas instrucciones exactas dentro de la página web de Nullplatform para configurar todo.

---

## PASO 1: Ingresar a tu Aplicación
1. Inicia sesión en la consola web de **Nullplatform**.
2. En el menú lateral izquierdo, haz clic en **Applications** (o Applications / Services).
3. Usa el buscador y escribe el nombre de tu app: `api-bfcl-mortgage-loans-injection`.
4. Haz clic sobre ella para entrar al Panel Principal (Dashboard) de tu microservicio.

---

## PASO 2: Crear los Parámetros y Secretos
Las variables de entorno no se inyectan solas; debes crearlas para cada "Scope" (ambiente).
1. Dentro del menú de tu aplicación, busca la pestaña que dice **Configuration** o **Parameters**.
2. Verás un menú desplegable para elegir el **Scope** (asegúrate de seleccionar el que corresponda, por ejemplo, `qa`).
3. Haz clic en el botón **+ Create Parameter** o **Add Variable**.
4. Agrega una por una las siguientes variables:
   
   **Variable 1 (Puerto):**
   - **Key/Name:** `PORT`
   - **Type:** `Plaintext` o `String`
   - **Value:** `8081`
   - Guarda los cambios.

   **Variable 2 (URL de Finnflow):**
   - **Key/Name:** `FINNFLOW_URL`
   - **Type:** `Plaintext`
   - **Value:** *(Pega aquí la URL del ambiente Legacy, ej: https://qa.finnflow.internal)*
   - Guarda los cambios.

   **Variable 3 (Credenciales Seguras):**
   - **Key/Name:** `FINNFLOW_SECRET`
   - **Type:** Cambia el tipo a **Secret** o **Encrypted** (¡Muy importante! Esto oculta la contraseña para siempre).
   - **Value:** *(Pega aquí el token o password real)*
   - Guarda los cambios.

*(Repite este mismo proceso si necesitas configurar el Scope de `prod`)*.

---

## PASO 3: Generar la API Key para GitHub
Para que GitHub Actions pueda compilar y avisarle a Nullplatform, necesita una llave de acceso.
1. En el menú lateral izquierdo (fuera de tu aplicación, a nivel general), ve a **Settings** o **Access Control**.
2. Haz clic en la sección **API Keys** o **Service Accounts**.
3. Haz clic en **Generate New API Key**.
4. Ponle un nombre descriptivo, por ejemplo: `github-actions-mortgage-injection`.
5. **¡Copia el valor que aparecerá en pantalla!** (Solo se muestra una vez).
6. Ve a tu repositorio en **GitHub** > **Settings** > **Secrets and variables** > **Actions** > **New repository secret**.
7. Nombra el secreto como `NULLPLATFORM_API_KEY` y pega el valor que copiaste.

---

## PASO 4: Validar el Ingress (Ruta Web)
1. Vuelve a tu aplicación en Nullplatform (`api-bfcl-mortgage-loans-injection`).
2. Ve a la pestaña **Networking** o **Ingress**.
3. Revisa que el dominio expuesto tenga sentido (ej. `api-bfcl-mortgage-loans.bancofalabella.cl`).
4. Revisa que exista un **Healthcheck** configurado apuntando a la ruta `/health` y al puerto `8081`. Si no existe, créalo.

---

## PASO 5: ¡Magia! (Ver tu Despliegue)
Una vez que hayas completado los 4 pasos anteriores en la pantalla de Nullplatform, ya no necesitas hacer nada más allí.

1. Ve a GitHub y haz *Merge* de tu Pull Request.
2. Si vuelves a Nullplatform, ve a la pestaña **Builds**. Verás que mágicamente aparece una nueva compilación cargando.
3. Cuando termine, ve a la pestaña **Releases** o **Deployments** y verás cómo tu contenedor se levanta y se pone en color **Verde (Healthy)**.
