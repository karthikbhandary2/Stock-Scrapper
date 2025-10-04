package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/gocolly/colly"
)

type Stock struct {
	company string
	price   string
	change  string
}

func main() {
	ticker := []string{
		"MSFT",
		"IBM", 
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"WMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAPL",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
	}
	stocks := []Stock{}
	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       3 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Connection", "keep-alive")
		fmt.Println("Visiting", r.URL)
	})
	
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	// Store current ticker being processed
	currentTicker := ""
	
	c.OnHTML("body", func(e *colly.HTMLElement) {
		stock := Stock{}
		
		// Get company name
		e.ForEach("h1", func(i int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if stock.company == "" && text != "Yahoo Finance" && text != "" && !strings.Contains(text, "Yahoo") {
				stock.company = text
			}
		})
		
		// Collect all prices and find the most reasonable one for a stock
		var prices []string
		e.ForEach("fin-streamer[data-field='regularMarketPrice']", func(i int, el *colly.HTMLElement) {
			price := strings.TrimSpace(el.Text)
			prices = append(prices, price)
		})
		
		// Strategy: Find a price that makes sense for the stock
		// Most major stocks are between $20-$600
		for _, priceStr := range prices {
			cleanPrice := strings.ReplaceAll(priceStr, ",", "")
			if price, err := strconv.ParseFloat(cleanPrice, 64); err == nil {
				// Adjust range based on known stock ranges
				if currentTicker == "MSFT" && price >= 300 && price <= 600 {
					stock.price = priceStr
					break
				} else if currentTicker == "AAPL" && price >= 150 && price <= 300 {
					stock.price = priceStr
					break
				} else if price >= 20 && price <= 800 {
					stock.price = priceStr
					break
				}
			}
		}
		
		// Fallback: if no reasonable price found, try to get the price from the main quote area
		if stock.price == "" {
			// Look for the price in a span with specific styling (main prices are usually larger)
			e.ForEach("span[data-field='regularMarketPrice']", func(i int, el *colly.HTMLElement) {
				if stock.price == "" {
					stock.price = strings.TrimSpace(el.Text)
				}
			})
		}
		
		// Another fallback: use the first reasonable price
		if stock.price == "" && len(prices) > 0 {
			for _, priceStr := range prices {
				cleanPrice := strings.ReplaceAll(priceStr, ",", "")
				if price, err := strconv.ParseFloat(cleanPrice, 64); err == nil {
					if price >= 1 && price <= 10000 { // Very broad range
						stock.price = priceStr
						break
					}
				}
			}
		}
		
		// Get change percentage
		e.ForEach("fin-streamer[data-field='regularMarketChangePercent']", func(i int, el *colly.HTMLElement) {
			if stock.change == "" {
				stock.change = strings.TrimSpace(el.Text)
			}
		})

		if stock.company != "" && stock.price != "" {
			fmt.Printf("Company: %s\n", stock.company)
			fmt.Printf("Price: %s\n", stock.price)
			fmt.Printf("Change: %s\n", stock.change)
			stocks = append(stocks, stock)
		}
	})

	for _, t := range ticker {
		currentTicker = t
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}
	c.Wait()

	fmt.Printf("Found %d stocks\n", len(stocks))
	
	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	headers := []string{"company", "price", "change"}

	writer.Write(headers)
	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
		}
		writer.Write(record)
	}
	writer.Flush()
	
	fmt.Println("Data written to stocks.csv")
}
