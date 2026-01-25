#!/bin/bash

echo "========================================"
echo "LaTeX Document Compilation Script"
echo "========================================"
echo ""

# Check if pdflatex is available
if ! command -v pdflatex &> /dev/null; then
    echo "[ERROR] pdflatex not found!"
    echo "Please install TeX Live: https://www.tug.org/texlive/"
    exit 1
fi

# Check required files
echo "[1/6] Checking required files..."
if [ ! -f "pract-example.tex" ]; then
    echo "[ERROR] pract-example.tex not found!"
    exit 1
fi
if [ ! -f "SCWorks.cls" ]; then
    echo "[ERROR] SCWorks.cls not found!"
    exit 1
fi
if [ ! -f "thesis.bib" ]; then
    echo "[WARNING] thesis.bib not found - bibliography will be empty"
fi
if [ ! -f "gost780uv.bst" ]; then
    echo "[WARNING] gost780uv.bst not found - bibliography style may not work"
fi
if [ ! -d "../diagrams" ]; then
    echo "[WARNING] diagrams folder not found - images may not be included"
fi
echo "[OK] All required files found"
echo ""

# Check for minted package (requires Python + Pygments)
echo "[2/6] Checking minted package requirements..."
if ! command -v python3 &> /dev/null && ! command -v python &> /dev/null; then
    echo "[WARNING] Python not found - minted code blocks may not work"
    echo "[INFO] Install Python and run: pip install pygments"
else
    if python3 -c "import pygments" 2>/dev/null || python -c "import pygments" 2>/dev/null; then
        echo "[OK] Python and Pygments are available"
    else
        echo "[WARNING] Pygments not installed - minted code blocks may not work"
        echo "[INFO] Run: pip install pygments"
    fi
fi
echo ""

# First pass - generate aux file
echo "[3/6] Running pdflatex (first pass)..."
echo "This may take a while..."
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if [ $? -ne 0 ]; then
    echo ""
    echo "[ERROR] pdflatex failed on first pass"
    echo "Check pract-example.log for details"
    exit 1
fi
echo "[OK] First pass completed"
echo ""

# Run bibtex for bibliography
if [ -f "thesis.bib" ]; then
    echo "[4/6] Running bibtex..."
    bibtex pract-example
    if [ $? -ne 0 ]; then
        echo "[WARNING] bibtex may have failed, continuing anyway..."
    else
        echo "[OK] Bibliography processed"
    fi
else
    echo "[4/6] Skipping bibtex (thesis.bib not found)"
fi
echo ""

# Second pass - include bibliography
echo "[5/6] Running pdflatex (second pass)..."
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if [ $? -ne 0 ]; then
    echo ""
    echo "[ERROR] pdflatex failed on second pass"
    echo "Check pract-example.log for details"
    exit 1
fi
echo "[OK] Second pass completed"
echo ""

# Third pass - resolve all references
echo "[6/6] Running pdflatex (third pass - final)..."
pdflatex -interaction=nonstopmode -shell-escape pract-example.tex
if [ $? -ne 0 ]; then
    echo ""
    echo "[ERROR] pdflatex failed on third pass"
    echo "Check pract-example.log for details"
    exit 1
fi
echo "[OK] Third pass completed"
echo ""

echo "========================================"
echo "[SUCCESS] Compilation completed!"
echo "========================================"
echo ""
echo "Output file: pract-example.pdf"
echo "Location: $(pwd)/pract-example.pdf"
echo ""

# Ask if user wants to clean auxiliary files
read -p "Clean auxiliary files (*.aux, *.log, *.out, *.toc, *.bbl, *.blg)? (y/n): " CLEAN
if [ "$CLEAN" = "y" ] || [ "$CLEAN" = "Y" ]; then
    echo "Cleaning up..."
    rm -f *.aux *.log *.out *.toc *.bbl *.blg *.synctex.gz
    echo "[OK] Cleaned up auxiliary files"
else
    echo "Keeping auxiliary files for debugging"
fi

echo ""
if [ -f "pract-example.pdf" ]; then
    echo "Opening PDF..."
    if command -v xdg-open &> /dev/null; then
        xdg-open pract-example.pdf
    elif command -v open &> /dev/null; then
        open pract-example.pdf
    fi
else
    echo "[WARNING] PDF file not found!"
fi
