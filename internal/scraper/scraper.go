package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

type ingredient struct {
	quantity string
	unit     string
	name     string
}

func ScrapeDelishSite(url string, scraper *colly.Collector) string {
	fmt.Println(url)
	err := scraper.Visit(url)
	if err != nil {
		fmt.Printf("Error visiting: %s\n", err.Error())
		return ""
	}

	scraper.OnError(func(_ *colly.Response, err error) {
		fmt.Println(err.Error())
	})

	scraper.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL)
	})

	ingredients := make([]ingredient, 0)
	scraper.OnHTML("li.mntl-structured-ingredients__list-item", func(h *colly.HTMLElement) {
		ingr := ingredient{}
		ingr.quantity = h.ChildText("span[data-ingredient-quantity]")
		ingr.unit = h.ChildText("span[data-ingredient-unit]")
		ingr.name = h.ChildText("span[data-ingredient-name]")

		ingredients = append(ingredients, ingr)
	})
	scraper.OnHTML("html", func(h *colly.HTMLElement) {
		fmt.Println(h.ChildText("body"))
	})

	return fmt.Sprintf("INGREDIENTS: %d\n", len(ingredients))
}
