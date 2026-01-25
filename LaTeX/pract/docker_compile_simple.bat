@echo off
chcp 65001 >nul
echo ========================================
echo LaTeX Compilation via Docker
echo (No LaTeX installation needed!)
echo ========================================
echo.

REM Check Docker
where docker >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker not found!
    echo.
    echo Please install Docker Desktop from:
    echo https://www.docker.com/products/docker-desktop
    echo.
    pause
    exit /b 1
)

echo [INFO] Docker found!
echo [INFO] First run will download LaTeX image (~3GB)
echo [INFO] This may take 10-15 minutes...
echo.

REM Get current directory (absolute path for Windows)
cd /d "%~dp0"
for %%I in ("%CD%") do set WORKDIR=%%~fI

echo [INFO] Working directory: %WORKDIR%
echo.

REM Get diagrams folder path
for %%I in ("%WORKDIR%\..\diagrams") do set DIAGRAMS_DIR=%%~fI

REM Compile using Docker
echo [1/1] Compiling in Docker container...
echo [INFO] This may take 10-15 minutes on first run (downloading LaTeX image)
echo.

docker run --rm ^
    -v "%WORKDIR%":/workdir ^
    -v "%DIAGRAMS_DIR%":/workdir/../diagrams ^
    -w /workdir ^
    texlive/texlive:latest ^
    bash -c "pdflatex -interaction=nonstopmode -shell-escape pract-example.tex && (test -f thesis.bib && bibtex pract-example || true) && pdflatex -interaction=nonstopmode -shell-escape pract-example.tex && pdflatex -interaction=nonstopmode -shell-escape pract-example.tex"

if errorlevel 1 (
    echo.
    echo [ERROR] Compilation failed!
    echo Check the error messages above.
    pause
    exit /b 1
)

echo.
echo ========================================
echo [SUCCESS] PDF created: pract-example.pdf
echo ========================================
echo.

if exist "pract-example.pdf" (
    echo Opening PDF...
    start pract-example.pdf
) else (
    echo [WARNING] PDF file not found!
)

pause
