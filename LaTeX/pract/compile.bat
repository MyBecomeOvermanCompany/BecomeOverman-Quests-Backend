@echo off
echo Compiling LaTeX document...
echo.

REM Try to find MiKTeX if not in PATH
where pdflatex >nul 2>&1
if errorlevel 1 (
    echo MiKTeX not found in PATH, trying to locate it...
    if exist "C:\Program Files\MiKTeX\miktex\bin\x64\pdflatex.exe" (
        set "MIKTEX_BIN=C:\Program Files\MiKTeX\miktex\bin\x64"
        set "PATH=%MIKTEX_BIN%;%PATH%"
        echo Found MiKTeX at %MIKTEX_BIN%
    ) else if exist "C:\Program Files (x86)\MiKTeX\miktex\bin\x64\pdflatex.exe" (
        set "MIKTEX_BIN=C:\Program Files (x86)\MiKTeX\miktex\bin\x64"
        set "PATH=%MIKTEX_BIN%;%PATH%"
        echo Found MiKTeX at %MIKTEX_BIN%
    ) else (
        echo ERROR: Could not find MiKTeX installation.
        echo Please add MiKTeX to your PATH or install it.
        pause
        exit /b 1
    )
)

REM First pass - generate aux file
echo [1/4] Running pdflatex (first pass)...
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo ERROR: pdflatex failed on first pass
    pause
    exit /b 1
)

REM Run bibtex for bibliography
echo [2/4] Running bibtex...
bibtex pract-example
if errorlevel 1 (
    echo WARNING: bibtex may have failed, continuing anyway...
)

REM Second pass - include bibliography
echo [3/4] Running pdflatex (second pass)...
pdflatex -interaction=nonstopmode pract-example.tex
if errorlevel 1 (
    echo ERROR: pdflatex failed on second pass
    pause
    exit /b 1
)

REM Third pass - resolve all references
echo [4/4] Running pdflatex (third pass)...
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo ERROR: pdflatex failed on third pass
    pause
    exit /b 1
)

echo.
echo ========================================
echo Compilation completed successfully!
echo Output file: pract-example.pdf
echo ========================================
echo.

REM Clean up auxiliary files (optional)
echo Cleaning up auxiliary files...
del *.aux *.log *.out *.toc *.bbl *.blg 2>nul

pause
