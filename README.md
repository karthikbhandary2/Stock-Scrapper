# Stock Scrapper

A Go-based web scraper that extracts stock data from Yahoo Finance and exports it to CSV format.

## Project Structure

```
Stock-Scrapper/
├── README.md              # This file
├── main.go               # Main scraper application
├── go.mod                # Go module dependencies
├── go.sum                # Go module checksums
├── stocks.csv            # Output CSV file (generated)
└── .git/                 # Git repository
```

## What This Code Does

The Stock Scrapper:

1. **Scrapes Yahoo Finance** for stock data from 18 major companies
2. **Extracts** company names, current prices, and percentage changes
3. **Exports** data to `stocks.csv` in a structured format
4. **Handles** rate limiting and proper HTTP headers to avoid being blocked

### Target Stocks
- MSFT (Microsoft)
- IBM (International Business Machines)
- GE (General Electric)
- UNP (Union Pacific)
- COST (Costco)
- MCD (McDonald's)
- V (Visa)
- WMT (Walmart)
- DIS (Disney)
- MMM (3M)
- INTC (Intel)
- AXP (American Express)
- AAPL (Apple)
- BA (Boeing)
- CSCO (Cisco)
- GS (Goldman Sachs)
- JPM (JPMorgan Chase)
- CRM (Salesforce)

## Setup Instructions

### Prerequisites

1. **Go Programming Language** (version 1.19 or higher)
   ```bash
   # Check if Go is installed
   go version
   
   # If not installed, download from https://golang.org/dl/
   ```

2. **Internet Connection** (required for scraping Yahoo Finance)

### Installation

1. **Clone or download the project**
   ```bash
   git clone <repository-url>
   cd Stock-Scrapper
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

   This will install:
   - `github.com/gocolly/colly` - Web scraping framework

### Running the Application

1. **Execute the main scraper**
   ```bash
   go run main.go
   ```

2. **Expected output**
   ```
   Visiting https://finance.yahoo.com/quote/MSFT/
   Company: Microsoft Corporation (MSFT)
   Price: 420.50
   Change: (+1.25%)
   ...
   Found 18 stocks
   Data written to stocks.csv
   ```

3. **Check the results**
   ```bash
   cat stocks.csv
   ```

## Output Format

The generated `stocks.csv` contains:

| Column  | Description                    | Example                        |
|---------|--------------------------------|--------------------------------|
| company | Full company name with ticker  | Microsoft Corporation (MSFT)   |
| price   | Current stock price            | 420.50                         |
| change  | Percentage change              | (+1.25%)                       |

## Configuration

### Modify Target Stocks

Edit the `ticker` slice in `main.go`:

```go
ticker := []string{
    "MSFT",
    "AAPL",
    // Add more stock symbols here
}
```

### Adjust Scraping Speed

Modify the delay in `main.go`:

```go
c.Limit(&colly.LimitRule{
    DomainGlob:  "*",
    Parallelism: 1,
    Delay:       3 * time.Second,  // Change this value
})
```

## Troubleshooting

### Common Issues

1. **Empty CSV file**
   - Ensure internet connection is active
   - Yahoo Finance may be blocking requests - try increasing delay

2. **"Service Unavailable" errors**
   - Yahoo Finance is rate limiting - increase delay between requests
   - Check if Yahoo Finance is accessible in your region

3. **Incorrect stock prices**
   - Yahoo Finance structure may have changed
   - Consider using a dedicated financial API instead

## Technical Details

### Dependencies

- **Colly v1.2.0**: Web scraping framework
  - Handles HTTP requests and HTML parsing
  - Provides rate limiting and error handling
  - Supports CSS selectors for data extraction

### Key Features

- **Rate Limiting**: 3-second delays between requests
- **Browser Simulation**: Proper User-Agent headers
- **Error Handling**: Graceful handling of failed requests
- **CSV Export**: Clean, structured data output
- **Concurrent Safety**: Single-threaded scraping to avoid blocks

### Data Extraction Strategy

1. **Company Names**: Extracted from `<h1>` elements
2. **Stock Prices**: Found via `fin-streamer[data-field='regularMarketPrice']`
3. **Price Changes**: Located using `fin-streamer[data-field='regularMarketChangePercent']`

## Limitations

- **Yahoo Finance Dependency**: Relies on Yahoo Finance's HTML structure
- **Rate Limiting**: Takes ~1 minute to scrape all 18 stocks
- **Data Accuracy**: Yahoo Finance may show multiple prices per page
- **Regional Restrictions**: May not work in all countries

## Legal Notice

This tool is for educational purposes. Ensure compliance with:
- Yahoo Finance's Terms of Service
- Local web scraping regulations
- Rate limiting to avoid overloading servers

