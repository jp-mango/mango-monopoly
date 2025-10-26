package scraper

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"mango-monopoly/internal/data"
)

const gwinnettCountyKey = "gwinnett"

type gwinnettScraper struct {
	allowedDomains    []string
	userAgent         string
	propertyDetailURL string
}

// NewGwinnettScraper constructs a scraper tailored to Gwinnett county's site.
func NewGwinnettScraper() PropertyScraper {
	return &gwinnettScraper{
		allowedDomains:    []string{"gwinnettassessor.manatron.com"},
		userAgent:         "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		propertyDetailURL: "https://gwinnettassessor.manatron.com/IWantTo/PropertyGISSearch/PropertyDetail.aspx?p=%s",
	}
}

func init() {
	MustRegister(gwinnettCountyKey, NewGwinnettScraper())
}

func (s *gwinnettScraper) ScrapeProperty(ctx context.Context, propID string) (*data.Property, error) {
	propID = strings.TrimSpace(propID)
	if propID == "" {
		return nil, fmt.Errorf("propID cannot be empty")
	}

	target := fmt.Sprintf(s.propertyDetailURL, url.QueryEscape(propID))
	property := &data.Property{
		PropID:    propID,
		County:    "Gwinnett",
		SourceURL: target,
	}

	c := colly.NewCollector(
		colly.AllowedDomains(s.allowedDomains...),
		colly.UserAgent(s.userAgent),
		colly.MaxDepth(1),
	)

	var visitErr error
	c.OnError(func(_ *colly.Response, err error) {
		visitErr = err
	})

	c.OnRequest(func(r *colly.Request) {
		select {
		case <-ctx.Done():
			r.Abort()
		default:
		}
	})

	c.OnHTML("div#lxT1385 table.generalinfo", func(e *colly.HTMLElement) {
		parseGwinnettGeneralInfo(e, property)
	})

	c.OnHTML("div#lxT1386 table#ValueHistory", func(e *colly.HTMLElement) {
		parseGwinnettValueHistory(e, property)
	})

	c.OnHTML("div#lxT1696 table.ui-table", func(e *colly.HTMLElement) {
		parseGwinnettTransferHistory(e, property)
	})

	c.OnHTML("div#1388tabs", func(e *colly.HTMLElement) {
		parseGwinnettImprovements(e, property)
	})

	c.OnHTML("div#lxT1390 table.ui-table", func(e *colly.HTMLElement) {
		parseGwinnettLandDetails(e, property)
	})

	c.OnHTML("div#lxT1391 table.ui-table", func(e *colly.HTMLElement) {
		parseGwinnettLegalDescription(e, property)
	})

	if err := c.Visit(target); err != nil {
		return nil, err
	}
	c.Wait()

	if visitErr != nil {
		return nil, visitErr
	}

	return property, ctx.Err()
}

func parseGwinnettGeneralInfo(e *colly.HTMLElement, p *data.Property) {
	ownersCell := e.DOM.Find("td.ui-widget-content.center").First()
	if ownersCell.Length() > 0 {
		lines := extractLines(ownersCell.Text())
		var owners []string
		var addressLines []string

		for _, line := range lines {
			if containsDigit(line) && len(addressLines) == 0 {
				addressLines = append(addressLines, line)
				continue
			}

			if len(addressLines) > 0 {
				addressLines = append(addressLines, line)
				continue
			}

			owners = append(owners, line)
		}

		if len(owners) > 0 {
			p.Owners = owners
		}

		if len(addressLines) > 0 {
			p.OwnerAddress = addressLines[0]
		}

		if len(addressLines) > 1 {
			city, state, zip := splitCityStateZip(addressLines[1])
			if city != "" {
				p.City = city
			}
			if state != "" {
				p.State = state
			}
			if zip != "" {
				p.ZIP = zip
			}
		}
	}

	e.DOM.Find("tr").Each(func(_ int, row *goquery.Selection) {
		label := collapseSpaces(row.Find("th").First().Text())
		value := collapseSpaces(row.Find("td").Last().Text())

		switch label {
		case "Property ID":
			if value != "" {
				p.PropID = value
			}
		case "Alternate ID":
			p.AltID = value
		case "Address":
			if value != "" {
				p.SitusAddress = value
			}
		case "Property Class":
			p.Class = value
		case "Neighborhood":
			p.Neighborhood = value
		case "Deed Acres":
			if acres := parseFloat64(value); acres > 0 {
				p.Acres = acres
			}
		}
	})
}

func parseGwinnettValueHistory(e *colly.HTMLElement, p *data.Property) {
	headerCells := e.DOM.Find("tr").First().Find("th")
	if headerCells.Length() >= 2 {
		if year, err := strconv.Atoi(strings.TrimSpace(headerCells.Eq(1).Text())); err == nil {
			p.ValueYear = int32(year)
		}
	}

	e.DOM.Find("tr").Each(func(_ int, row *goquery.Selection) {
		label := collapseSpaces(row.Find("th").First().Text())
		valCell := row.Find("td").First()
		if valCell.Length() == 0 {
			return
		}

		text := collapseSpaces(valCell.Text())
		raw := valCell.AttrOr("title", text)

		switch label {
		case "Reason":
			p.ValueReason = text
		case "Land Val":
			p.LandValue = parseInt64(raw)
		case "Imp Val":
			p.ImprovementValue = parseInt64(raw)
		case "Total Appr":
			p.TotalAppraised = parseInt64(raw)
		case "Total Assd":
			p.TotalAssessed = parseInt64(raw)
		}
	})
}

func parseGwinnettTransferHistory(e *colly.HTMLElement, p *data.Property) {
	row := e.DOM.Find("tr").Eq(1)
	if row.Length() == 0 {
		return
	}

	cells := row.Find("td")
	if cells.Length() < 9 {
		return
	}

	p.TransferBook = collapseSpaces(cells.Eq(0).Text())
	p.TransferPage = collapseSpaces(cells.Eq(1).Text())

	if dateStr := collapseSpaces(cells.Eq(2).Text()); dateStr != "" {
		if dt, err := time.Parse("1/2/2006", dateStr); err == nil {
			p.TransferDate = dt
		}
	}

	p.Grantor = collapseSpaces(cells.Eq(3).Text())
	p.Grantee = collapseSpaces(cells.Eq(4).Text())

	if deed := strings.TrimSpace(cells.Eq(5).Find("a").AttrOr("title", "")); deed != "" {
		p.DeedType = deed
	} else {
		p.DeedType = collapseSpaces(cells.Eq(5).Text())
	}

	p.VacantLand = strings.EqualFold(collapseSpaces(cells.Eq(7).Text()), "Yes")
	p.SalePrice = parseInt64(cells.Eq(8).Text())
}

func parseGwinnettImprovements(e *colly.HTMLElement, p *data.Property) {
	e.DOM.Find("#1388gallerywrap table.ui-table tr").Each(func(_ int, row *goquery.Selection) {
		label := collapseSpaces(row.Find("th").Text())
		value := collapseSpaces(row.Find("td").Text())

		switch label {
		case "Address":
			if value != "" {
				p.SitusAddress = value
			}
		case "Type":
			p.Type = value
		case "Grade":
			p.Grade = value
		case "Year Built":
			if year := parseInt(value); year > 0 {
				p.YearBuilt = int16(year)
			}
		}
	})

	e.DOM.Find("#1388Attributes table#Attribute tbody tr").Each(func(_ int, row *goquery.Selection) {
		attr := collapseSpaces(row.Find("td").Eq(1).Text())
		detail := collapseSpaces(row.Find("td").Eq(2).Text())

		switch attr {
		case "Type":
			p.Type = detail
		case "Occupancy":
			p.Occupancy = detail
		case "Roof Structure":
			p.RoofStruct = detail
		case "Roof Cover":
			p.RoofCover = detail
		case "Heating":
			p.Heating = detail
		case "A/C":
			p.AC = detail
		case "Stories":
			p.Stories = parseFloat32(detail)
		case "Bedrooms":
			if bedrooms := parseInt(detail); bedrooms > 0 {
				p.Bedrooms = int16(bedrooms)
			}
		case "Bathrooms":
			p.Bathrooms = parseFloat32(detail)
		case "Exterior Wall":
			p.ExteriorWall = detail
		case "Interior Flooring":
			p.InteriorFloor = detail
		}
	})

	if row := e.DOM.Find("#1388Areas table#Area tbody tr").First(); row.Length() > 0 {
		p.FloorAreaCode = collapseSpaces(row.Find("td").Eq(0).Text())
		p.FloorDesc = collapseSpaces(row.Find("td").Eq(1).Text())
		p.GrossAreaSF = parseInt32(row.Find("td").Eq(2).Text())
		p.FinishedSF = parseInt32(row.Find("td").Eq(3).Text())
		p.FloorConstr = collapseSpaces(row.Find("td").Eq(4).Text())
	}

	if row := e.DOM.Find("#1388Features table#Feature tbody tr").First(); row.Length() > 0 {
		p.ExtFeatCode = collapseSpaces(row.Find("td").Eq(0).Text())
		p.ExtFeatDesc = collapseSpaces(row.Find("td").Eq(1).Text())
		p.ExtFeatSizeSF = parseInt32(row.Find("td").Eq(2).Text())
		p.ExtFeatConstr = collapseSpaces(row.Find("td").Eq(3).Text())
	}

	e.DOM.Find("#1388gallerywrap img").Each(func(_ int, el *goquery.Selection) {
		if src, ok := el.Attr("src"); ok {
			src = strings.TrimSpace(src)
			if src != "" {
				p.Photos = appendUnique(p.Photos, src)
			}
		}
	})
}

func parseGwinnettLandDetails(e *colly.HTMLElement, p *data.Property) {
	row := e.DOM.Find("tr").Eq(1)
	if row.Length() == 0 {
		return
	}

	cells := row.Find("td")
	if cells.Length() < 5 {
		return
	}

	p.LandPrimaryUse = collapseSpaces(cells.Eq(0).Text())
	p.LandType = collapseSpaces(cells.Eq(1).Text())
	if acres := parseFloat64(cells.Eq(2).Text()); acres > 0 {
		p.Acres = acres
	}
	p.EffFrontage = parseInt32(cells.Eq(3).Text())
	p.EffDepth = parseInt32(cells.Eq(4).Text())
}

func parseGwinnettLegalDescription(e *colly.HTMLElement, p *data.Property) {
	var descriptions []string
	e.DOM.Find("tr").Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return
		}

		text := collapseSpaces(row.Find("td").Eq(1).Text())
		if text != "" {
			descriptions = append(descriptions, text)
		}
	})

	if len(descriptions) > 0 {
		p.LegalDescription = descriptions
	}
}
