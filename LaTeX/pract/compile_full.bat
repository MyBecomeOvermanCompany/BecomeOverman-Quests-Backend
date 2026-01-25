@echo off
chcp 65001 >nul
echo ========================================
echo LaTeX Document Compilation Script
echo ========================================
echo.

REM Check if pdflatex is available
where pdflatex >nul 2>&1
if errorlevel 1 (
    echo [ERROR] pdflatex not found in PATH
    echo.
    echo Trying to find MiKTeX installation...
    if exist "C:\Program Files\MiKTeX\miktex\bin\x64\pdflatex.exe" (
        set "MIKTEX_BIN=C:\Program Files\MiKTeX\miktex\bin\x64"
        set "PATH=%MIKTEX_BIN%;%PATH%"
        echo [OK] Found MiKTeX at %MIKTEX_BIN%
    ) else if exist "C:\Program Files (x86)\MiKTeX\miktex\bin\x64\pdflatex.exe" (
        set "MIKTEX_BIN=C:\Program Files (x86)\MiKTeX\miktex\bin\x64"
        set "PATH=%MIKTEX_BIN%;%PATH%"
        echo [OK] Found MiKTeX at %MIKTEX_BIN%
    ) else (
        echo [ERROR] MiKTeX not found!
        echo.
        echo Please install MiKTeX from: https://miktex.org/download
        echo Or add MiKTeX to your PATH
        pause
        exit /b 1
    )
)

REM Check required files
echo [1/6] Checking required files...
if not exist "pract-example.tex" (
    echo [ERROR] pract-example.tex not found!
    pause
    exit /b 1
)
if not exist "SCWorks.cls" (
    echo [ERROR] SCWorks.cls not found!
    pause
    exit /b 1
)
if not exist "thesis.bib" (
    echo [WARNING] thesis.bib not found - bibliography will be empty
)
if not exist "gost780uv.bst" (
    echo [WARNING] gost780uv.bst not found - bibliography style may not work
)
if not exist "..\diagrams" (
    echo [WARNING] diagrams folder not found - images may not be included
)
echo [OK] All required files found
echo.

REM Check for minted package (requires Python + Pygments)
echo [2/6] Checking minted package requirements...
where python >nul 2>&1
if errorlevel 1 (
    echo [WARNING] Python not found - minted code blocks may not work
    echo [INFO] Install Python and run: pip install pygments
) else (
    python -c "import pygments" >nul 2>&1
    if errorlevel 1 (
        echo [WARNING] Pygments not installed - minted code blocks may not work
        echo [INFO] Run: pip install pygments
    ) else (
        echo [OK] Python and Pygments are available
    )
)
echo.

REM First pass - generate aux file
echo [3/6] Running pdflatex (first pass)...
echo This may take a while...
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo.
    echo [ERROR] pdflatex failed on first pass
    echo Check pract-example.log for details
    pause
    exit /b 1
)
echo [OK] First pass completed
echo.

REM Run bibtex for bibliography
if exist "thesis.bib" (
    echo [4/6] Running bibtex...
    bibtex pract-example
    if errorlevel 1 (
        echo [WARNING] bibtex may have failed, continuing anyway...
    ) else (
        echo [OK] Bibliography processed
    )
) else (
    echo [4/6] Skipping bibtex (thesis.bib not found)
)
echo.

REM Second pass - include bibliography
echo [5/6] Running pdflatex (second pass)...
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo.
    echo [ERROR] pdflatex failed on second pass
    echo Check pract-example.log for details
    pause
    exit /b 1
)
echo [OK] Second pass completed
echo.

REM Third pass - resolve all references
echo [6/6] Running pdflatex (third pass - final)...
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if errorlevel 1 (
    echo.
    echo [ERROR] pdflatex failed on third pass
    echo Check pract-example.log for details
    pause
    exit /b 1
)
echo [OK] Third pass completed
echo.

echo ========================================
echo [SUCCESS] Compilation completed!
echo ========================================
echo.
echo Output file: pract-example.pdf
echo Location: %CD%\pract-example.pdf
echo.

REM Ask if user wants to clean auxiliary files
set /p CLEAN="Clean auxiliary files (*.aux, *.log, *.out, *.toc, *.bbl, *.blg)? (y/n): "
if /i "%CLEAN%"=="y" (
    echo Cleaning up...
    del *.aux *.log *.out *.toc *.bbl *.blg *.synctex.gz 2>nul
    echo [OK] Cleaned up auxiliary files
) else (
    echo Keeping auxiliary files for debugging
)

echo.
echo Opening PDF...
if exist "pract-example.pdf" (
    start pract-example.pdf
) else (
    echo [WARNING] PDF file not found!
)

pause
