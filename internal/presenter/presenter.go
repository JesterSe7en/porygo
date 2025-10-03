// Package presenter handles formatting and printing the scraped data.
package presenter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/JesterSe7en/scrapego/internal/scraper"
)

type Presenter interface {
	Write(data any) error
}

type JSONPresenter struct {
	writer io.Writer
}

func NewJSONPresenter(w io.Writer) *JSONPresenter {
	return &JSONPresenter{writer: w}
}

func (p *JSONPresenter) Write(data any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result to JSON: %w", err)
	}
	_, err = fmt.Fprintln(p.writer, string(b))
	return err
}

type TextPresenter struct {
	writer io.Writer
}

func NewTextPresenter(w io.Writer) *TextPresenter {
	return &TextPresenter{writer: w}
}

func (p *TextPresenter) Write(data any) error {
	scrapedData, ok := data.(scraper.ScrapedData)
	if !ok {
		return fmt.Errorf("unexpected type for text presenter: expected ScrapedData, got %T", data)
	}

	var sb strings.Builder

	// --- Metadata ---
	sb.WriteString("--- Metadata ---\n")
	sb.WriteString(fmt.Sprintf("URL:          %s\n", scrapedData.URL))
	sb.WriteString(fmt.Sprintf("Status:       %d\n", scrapedData.Status))
	sb.WriteString(fmt.Sprintf("Content-Type: %s\n", scrapedData.ContentType))
	sb.WriteString(fmt.Sprintf("Size:         %d bytes\n", scrapedData.Size))
	sb.WriteString(fmt.Sprintf("Response Time: %s\n", scrapedData.ResponseTime))

	// --- Extracted Data ---
	if len(scrapedData.Extracted) > 0 {
		sb.WriteString("\n--- Extracted by CSS Selectors ---\n")
		for selector, items := range scrapedData.Extracted {
			sb.WriteString(fmt.Sprintf("Selector: %s\n", selector))
			if len(items) == 0 {
				sb.WriteString("  (No results found)\n")
				continue
			}
			for _, item := range items {
				// Indent and handle multi-line text
				indentedItem := "  - " + strings.ReplaceAll(item, "\n", "\n    ")
				sb.WriteString(indentedItem + "\n")
			}
		}
	}

	if len(scrapedData.Matches) > 0 {
		sb.WriteString("\n--- Matched by Regex Patterns ---\n")
		for pattern, items := range scrapedData.Matches {
			sb.WriteString(fmt.Sprintf("Pattern: %s\n", pattern))
			if len(items) == 0 {
				sb.WriteString("  (No matches found)\n")
				continue
			}
			for _, item := range items {
				sb.WriteString(fmt.Sprintf("  - %s\n", item))
			}
		}
	}

	_, err := fmt.Fprintln(p.writer, sb.String())
	return err
}
