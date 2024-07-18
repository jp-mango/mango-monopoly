import tabula

gwinnettPastSalesPDF = "./tax-auction/Gwinnett/pdf/Gwinnett-Past-Sales.pdf"
gwinnettUpcomingSalesPDF = "./tax-auction/Gwinnett/pdf/Gwinnett-Upcoming-Sales.pdf"

gwinnettPastSalesCSV = "./tax-auction/Gwinnett/csv/Gwinnett-Past-Sales.csv"
gwinnettUpcomingSalesCSV = "./tax-auction/Gwinnett/csv/Gwinnett-Upcoming-Sales.csv"

tabula.convert_into(
    gwinnettPastSalesPDF, gwinnettPastSalesCSV, output_format="csv", pages="all"
)
tabula.convert_into(
    gwinnettUpcomingSalesPDF, gwinnettUpcomingSalesCSV, output_format="csv", pages="all"
)