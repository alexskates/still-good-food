package food

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

type NutritionInfo struct {
	Kcal      string
	Fat       string
	Saturates string
	Carbs     string
	Sugars    string
	Fibre     string
	Protein   string
	Salt      string
}

type Recipe struct {
	Name        string
	Ingredients []string
	Steps       []string
	Yield       string
	Difficulty  string
	Preparation string
	Cook        string
	Nutrition   NutritionInfo
}

func Scrape(recipeURL string) (Recipe, error) {

	// Create a new context
	// With a deadline of 500 milliseconds
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

	// Make a request, that will call the google homepage
	req, _ := http.NewRequest(http.MethodGet, recipeURL, nil)
	// Associate the cancellable context we just created to the request
	req = req.WithContext(ctx)

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Recipe{}, err
	}

	// Load the response into a goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Recipe{}, err
	}

	// Get the recipe information from the text using class names
	name := doc.Find(".recipe-header__title").First().Text()

	var ingredients []string
	doc.Find(".ingredients-list__item").Each(func(i int, s *goquery.Selection) {
		s.Find("span").Remove()
		ingredients = append(ingredients, strings.TrimSpace(s.Text()))
	})

	var steps []string
	doc.Find(".method__item").Each(func(i int, s *goquery.Selection) {
		steps = append(steps, strings.TrimSpace(s.Text()))
	})

	yield := extract("recipeYield", doc)

	difficulty := strings.TrimSpace(doc.Find("section.recipe-details__item--skill-level").Text())

	preparationSpan := doc.Find(".recipe-details__cooking-time-prep")
	preparationSpan.Find("strong").Remove()

	preparation := strings.TrimSpace(preparationSpan.Text())

	cookSpan := doc.Find(".recipe-details__cooking-time-cook")
	cookSpan.Find("strong").Remove()

	cook := strings.TrimSpace(cookSpan.Text())

	nutrition := NutritionInfo{
		Kcal:      extract("calories", doc),
		Fat:       extract("fatContent", doc),
		Saturates: extract("saturatedFatContent", doc),
		Sugars:    extract("sugarContent", doc),
		Fibre:     extract("fiberContent", doc),
		Protein:   extract("proteinContent", doc),
		Salt:      extract("sodiumContent", doc),
	}

	r := Recipe{
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
		Yield:       yield,
		Difficulty:  difficulty,
		Preparation: preparation,
		Cook:        cook,
		Nutrition:   nutrition,
	}

	return r, nil
}
func extract(itemprop string, doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find(fmt.Sprintf("span[itemprop='%s']", itemprop)).Text())
}
