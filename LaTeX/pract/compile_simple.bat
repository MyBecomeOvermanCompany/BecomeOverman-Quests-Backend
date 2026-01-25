@echo off
chcp 65001 >nul
echo ========================================
echo LaTeX Compilation - Simple Version
echo ========================================
echo.

REM Check if pdflatex is available
where pdflatex >nul 2>&1
if errorlevel 1 (
    echo [ERROR] pdflatex not found!
    echo.
    echo Please install MiKTeX from: https://miktex.org/download
    echo Or use docker_compile.bat (requires Docker Desktop)
    echo.
    pause
    exit /b 1
)

echo [INFO] Found pdflatex, starting compilation...
echo [INFO] This may take 5-10 minutes on first run
echo.

REM Change to script directory
cd /d "%~dp0"

REM First pass
echo [1/4] pdflatex (first pass)...
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo [ERROR] Compilation failed! Check pract-example.log
    pause
    exit /b 1
)

REM Bibtex
if exist "thesis.bib" (
    echo [2/4] bibtex...
    bibtex pract-example >nul 2>&1
)

REM Second pass
echo [3/4] pdflatex (second pass)...
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo [ERROR] Compilation failed! Check pract-example.log
    pause
    exit /b 1
)

REM Third pass
echo [4/4] pdflatex (final pass)...
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo [ERROR] Compilation failed! Check pract-example.log
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
