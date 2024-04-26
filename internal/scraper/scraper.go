package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Recipe struct {
	title       string
	ingredients []ingredient
}

type ingredient struct {
	quantity string
	unit     string
	name     string
}

// TODO: create router func that takes in url and returns recipe site ID
// TODO: create recipe and ingredients and sites tables to keep track of all these things
// TODO: create scraper func that takes in scraper function and url and returns *recipe
// TODO: create recipes model
// TODO: create lists model

func ScrapeDelishSite(url string) string {
	scraper := colly.NewCollector()
	recipe := new(Recipe)

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

	var title string
	scraper.OnHTML("h1.article-heading", func(h *colly.HTMLElement) {
		title = h.Text
	})

	err := scraper.Visit(url)
	if err != nil {
		fmt.Printf("Error visiting: %s\n", err.Error())
		return ""
	}

	recipe.title = title
	recipe.ingredients = ingredients

	return fmt.Sprintf("INGREDIENTS: %+v\n", recipe)
}
