import fitz  # PyMuPDF
import csv

# Paths to the input PDF and output CSV files
pauldingUpcomingSalesPDF = "./tax-auction/Paulding/pdf/Paulding-Upcoming-Sales.pdf"
pauldingUpcomingSalesCSV = "./tax-auction/Paulding/csv/Paulding-Upcoming-Sales.csv"

# Open the PDF file
doc = fitz.open(pauldingUpcomingSalesPDF)

# Extract text from each page
all_text = ""
for page in doc:
    all_text += page.get_text()

# Process the extracted text into rows
lines = all_text.split("\n")
rows = [line.strip() for line in lines if line.strip().startswith("R")]

# Write the rows to a CSV file
with open(pauldingUpcomingSalesCSV, "w", newline="") as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(["ID"])
    for row in rows:
        writer.writerow([row])

print(f"Processed {len(rows)} rows into {pauldingUpcomingSalesCSV}")
