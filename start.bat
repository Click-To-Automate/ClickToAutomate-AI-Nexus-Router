@echo off
title AI Router Dev Environment

echo Stopping any existing instances on port 20128 (Backend)...
FOR /F "tokens=5" %%a IN ('netstat -a -n -o ^| findstr :20128') DO (
    if NOT "%%a"=="0" taskkill /F /PID %%a 2>nul
)

echo Stopping any existing instances on port 5173 (Frontend)...
FOR /F "tokens=5" %%a IN ('netstat -a -n -o ^| findstr :5173') DO (
    if NOT "%%a"=="0" taskkill /F /PID %%a 2>nul
)

echo Starting Frontend (Vite) on port 5173...
start "Frontend (Vite)" cmd /k "cd frontend && npm run dev"

echo Starting Backend (Go/Air) on port 20128...
:: We use go run to guarantee it works even if air is not in your system PATH
start "Backend (Go/Air)" cmd /k "cd backend && go run github.com/air-verse/air@latest"

echo Both environments have been started in new windows!
pause
