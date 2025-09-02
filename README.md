# Scrapgo – Concurrent Web Scraper CLI

Scrapgo is a command-line tool written in Go that demonstrates how to scrape data from web pages concurrently using goroutines and channels.

For demonstration purposes, the scraper currently fetches the **page title** (`<title>` tag) from each URL provided.

---

## How It Works

1. **Input**
   - You pass one or more URLs to the CLI.
   - Example:
     ```bash
     scrapgo https://example.com https://golang.org
     ```

2. **Concurrency**
   - Scrapgo uses a worker pool pattern with goroutines.
   - Multiple URLs are scraped in parallel, making the process much faster than sequential scraping.

3. **Scraping Logic**
   - For each URL:
     - An HTTP GET request is made.
     - The response body is parsed with [goquery](https://github.com/PuerkitoBio/goquery).
     - The `<title>` element is extracted.

   Example snippet:
   ```go
   resp, err := http.Get(url)
   if err != nil {
       return "", err
   }
   defer resp.Body.Close()

   doc, err := goquery.NewDocumentFromReader(resp.Body)
   if err != nil {
       return "", err
   }

   title := doc.Find("title").Text()
   ```

4. **Output**
   - Results are printed to the console or written to a file (JSON/CSV support can be added).
   - Example output:
     ```text
     https://example.com -> Example Domain
     https://golang.org  -> The Go Programming Language
     ```

---

## Example Command

```bash
scrapgo https://example.com https://golang.org
```

Output:
```text
Scraping 2 URLs with 5 workers...
https://example.com -> Example Domain
https://golang.org  -> The Go Programming Language
```

---

## Why This Project?

- **Resume Project**: Showcases concurrency, networking, and CLI design in Go.
- **Extendable**: You can expand it to scrape more than titles (e.g., headings, metadata, links).
- **Teaches Patterns**: Worker pool, error handling, and structured output are demonstrated.

---
