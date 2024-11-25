import tabula


def main():
    file = tabula.read_pdf(
        "https://www.gwinnetttaxcommissioner.com/documents/d/egov/decembertaxsalelist"
    )

    tabula.convert_into(file, "output.csv", output_format="csv", pages="all")


if __name__ == "__main__":
    main()
