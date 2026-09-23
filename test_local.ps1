$ErrorActionPreference = 'SilentlyContinue'

Write-Host "============================================="
Write-Host "   TEST DE PARIDAD (LOCAL CON PRE-IMAGEN)"
Write-Host "============================================="

# Funciones auxiliares
function Kill-Port {
    param([int]$Port)
    $pidToKill = (Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue).OwningProcess
    if ($pidToKill) {
        Write-Host "Deteniendo proceso $pidToKill ocupando el puerto $Port..."
        Stop-Process -Id $pidToKill -Force -ErrorAction SilentlyContinue
    }
}

function Wait-Http {
    param([int]$Port, [string]$Name)
    Write-Host "Esperando a que $Name responda peticiones HTTP en el puerto $Port..."
    $maxTries = 60
    $tries = 0
    while ($tries -lt $maxTries) {
        try {
            $response = Invoke-WebRequest -Uri "http://127.0.0.1:$Port/v1/bfcl/mortgage-loan/injections" -Method Post -Body "{}" -ContentType "application/json" -ErrorAction Stop
            Write-Host "$Name esta LISTO!"
            return
        } catch {
            if ($_.Exception.Response) {
                Write-Host "$Name esta LISTO! (Respondio HTTP $($_.Exception.Response.StatusCode))"
                return
            }
            $tries++
            Start-Sleep -Seconds 3
        }
    }
    Write-Host "TIMEOUT FATAL: $Name nunca levanto en el puerto $Port."
    exit 1
}

# 1. Liberar puertos
Write-Host "1. Limpiando puertos requeridos (8080, 8082, 8083)..."
Kill-Port 8080
Kill-Port 8082
Kill-Port 8083
Start-Sleep -Seconds 2

# 2. Compilar Go a binario nativo
Write-Host "2. Compilando Go en un binario nativo..."
go build -o api.exe cmd/api/main.go
if (-Not (Test-Path "api.exe")) {
    Write-Host "Error al compilar Go."
    exit 1
}

# 3. Levantando Contenedor Existente (Java) en puerto 8080...
Write-Host "3. Preparando contenedor 'java-app-injection' inyectandole variables de entorno QA..."

# Convertimos su contenedor en una imagen temporal para poder inyectarle las variables
docker commit java-app-injection temp-java-legacy-image > $null

Write-Host "Levantando el Legado en el puerto 8080..."
$javaContainer = docker run -d --rm -p 8080:8080 `
    -e FINNFLOW_URL="http://host.docker.internal:8083" `
    -e FINNFLOW_KEY="admin" `
    -e FINNFLOW_SECRET="admin" `
    temp-java-legacy-image

# 4. Iniciar Nueva App Go (Binario) en puerto 8082
Write-Host "4. Iniciando binario de Go en puerto 8082..."
$env:PORT="8082"
$env:FINNFLOW_URL="http://127.0.0.1:8083"
$env:FINNFLOW_KEY="admin"
$env:FINNFLOW_SECRET="admin"
$goProcess = Start-Process -FilePath ".\api.exe" -RedirectStandardOutput "go.log" -RedirectStandardError "go_err.log" -WindowStyle Minimized -PassThru

Write-Host "Esperando inteligentemente a que los servidores arranquen..."
Wait-Http -Port 8080 -Name "Java Legacy"
Wait-Http -Port 8082 -Name "Go Hexagonal"

# 5. Ejecutar Test QA
Write-Host "5. Ejecutando script de paridad extrema..."
Set-Location -Path "C:\HEXAGONAL-INJECTION\mortgage-loan-injection"
go run cmd/qa/main.go

# 6. Limpiar procesos
Write-Host "6. Limpiando procesos en segundo plano..."
docker stop java-app-injection
Stop-Process -Id $goProcess.Id -Force
Kill-Port 8080
Kill-Port 8082
Kill-Port 8083

Write-Host "============================================="
Write-Host "        TEST LOCAL FINALIZADO"
Write-Host "============================================="
