from os import error
from datetime import datetime
import tabula


def main():
    gwinnettUpcomingSales = (
        "https://www.gwinnetttaxcommissioner.com/documents/d/egov/decembertaxsalelist"
    )
    pdfConvert(gwinnettUpcomingSales, "Gwinnett")


def pdfConvert(link: str, county: str) -> error:
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")

    filepath = f"./scraper/{county}/{county}UpcomingSales{timestamp}.csv"
    tabula.io.convert_into(link, filepath, output_format="csv", pages="all")


if __name__ == "__main__":
    main()
