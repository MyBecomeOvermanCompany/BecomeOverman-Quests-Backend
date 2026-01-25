@echo off
echo ========================================
echo LaTeX Compilation using Docker
echo ========================================
echo.

REM Check if Docker is available
where docker >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker not found!
    echo.
    echo Please install Docker Desktop from:
    echo https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

echo [INFO] Docker found, starting compilation...
echo [INFO] This will download LaTeX image (~3GB) on first run
echo.

REM Get absolute path
for %%I in ("%CD%") do set ABS_PATH=%%~fI

echo [1/4] Running pdflatex (first pass)...
docker run --rm -v "%ABS_PATH%":/workdir -w /workdir texlive/texlive:latest pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo [ERROR] First pass failed
    pause
    exit /b 1
)

if exist "thesis.bib" (
    echo [2/4] Running bibtex...
    docker run --rm -v "%ABS_PATH%":/workdir -w /workdir texlive/texlive:latest bibtex pract-example
)

echo [3/4] Running pdflatex (second pass)...
docker run --rm -v "%ABS_PATH%":/workdir -w /workdir texlive/texlive:latest pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo [ERROR] Second pass failed
    pause
    exit /b 1
)

echo [4/4] Running pdflatex (third pass - final)...
docker run --rm -v "%ABS_PATH%":/workdir -w /workdir texlive/texlive:latest pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo [ERROR] Third pass failed
    pause
    exit /b 1
)

echo.
echo ========================================
echo [SUCCESS] Compilation completed!
echo ========================================
echo.
echo Output file: pract-example.pdf

if exist "pract-example.pdf" (
    start pract-example.pdf
)

pause
