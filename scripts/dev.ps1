# TicketFlow Development Server Runner
$ErrorActionPreference = "Stop"

$env:APP_ENV = "development"
$env:PORT = "8080"
$env:LOG_LEVEL = "debug"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host " Starting TicketFlow API (Development)  " -ForegroundColor Cyan
Write-Host " URL: http://localhost:$env:PORT       " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

go run ./cmd/api
