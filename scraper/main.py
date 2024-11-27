import glob
import sys
from datetime import datetime
import tabula
import os


def main():
    if len(sys.argv) < 3:
        print("Usage: python main.py <link> <county>")
        sys.exit(1)

    link = sys.argv[1]
    county = sys.argv[2]

    pdfConvert(link, county)


def pdfConvert(link: str, county: str):
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")

    directory = f"./scraper/{county}"
    if not os.path.exists(directory):
        os.makedirs(directory)

    files = glob.glob(os.path.join(directory, "*"))  # Match files inside the directory

    for f in files:
        if os.path.isfile(f):  # Check if it's a file
            os.remove(f)

    filepath = os.path.join(directory, f"{county}UpcomingSales{timestamp}.csv")
    tabula.convert_into(link, filepath, output_format="csv", pages="all")

    print(f"PDF downloaded and converted to CSV at {filepath}")


if __name__ == "__main__":
    main()
