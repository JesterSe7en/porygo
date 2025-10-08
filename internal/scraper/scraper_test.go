package scraper

import (
	"testing"
)

func TestApplySelectors(t *testing.T) {
	html := `
	<html>
		<body>
			<h1>Hello, World!</h1>
			<a href="/page1">Page 1</a>
			<a href="/page2">Page 2</a>
			<a href="https://www.google.com">Google</a>
		</body>
	</html>
	`
	s := &Scraper{}

	t.Run("Test with single selector", func(t *testing.T) {
		selectors := []string{"h1"}
		expectedResults := map[string][]string{
			"h1": {"Hello, World!"},
		}
		expectedTexts := []string{"Hello, World!"}

		results, texts := s.applySelectors([]byte(html), selectors)

		if len(results) != len(expectedResults) {
			t.Errorf("Expected %d results, but got %d", len(expectedResults), len(results))
		}

		for selector, expectedValues := range expectedResults {
			values, ok := results[selector]
			if !ok {
				t.Errorf("Expected selector %s not found in results", selector)
				continue
			}

			if len(values) != len(expectedValues) {
				t.Errorf("Expected %d values for selector %s, but got %d", len(expectedValues), selector, len(values))
				continue
			}

			for i, expectedValue := range expectedValues {
				if values[i] != expectedValue {
					t.Errorf("Expected value %s, but got %s", expectedValue, values[i])
				}
			}
		}

		if len(texts) != len(expectedTexts) {
			t.Errorf("Expected %d texts, but got %d", len(expectedTexts), len(texts))
		}
	})

	t.Run("Test with multiple selectors", func(t *testing.T) {
		selectors := []string{"h1", "a@href"}
		expectedResults := map[string][]string{
			"h1":     {"Hello, World!"},
			"a@href": {"/page1", "/page2", "https://www.google.com"},
		}
		expectedTexts := []string{"Hello, World!", "/page1", "/page2", "https://www.google.com"}

		results, texts := s.applySelectors([]byte(html), selectors)

		if len(results) != len(expectedResults) {
			t.Errorf("Expected %d results, but got %d", len(expectedResults), len(results))
		}

		for selector, expectedValues := range expectedResults {
			values, ok := results[selector]
			if !ok {
				t.Errorf("Expected selector %s not found in results", selector)
				continue
			}

			if len(values) != len(expectedValues) {
				t.Errorf("Expected %d values for selector %s, but got %d", len(expectedValues), selector, len(values))
				continue
			}

			for i, expectedValue := range expectedValues {
				if values[i] != expectedValue {
					t.Errorf("Expected value %s, but got %s", expectedValue, values[i])
				}
			}
		}

		if len(texts) != len(expectedTexts) {
			t.Errorf("Expected %d texts, but got %d", len(expectedTexts), len(texts))
		}
	})
}

func TestApplyRegexPatterns(t *testing.T) {
	texts := []string{"hello world", "hello there", "general kenobi"}
	s := &Scraper{}

	t.Run("Test with single pattern", func(t *testing.T) {
		patterns := []string{"hello"}
		expectedResults := map[string][]string{
			"hello": {"hello", "hello"},
		}

		results := s.applyRegexPatterns(texts, patterns)

		if len(results) != len(expectedResults) {
			t.Errorf("Expected %d results, but got %d", len(expectedResults), len(results))
		}

		for pattern, expectedValues := range expectedResults {
			values, ok := results[pattern]
			if !ok {
				t.Errorf("Expected pattern %s not found in results", pattern)
				continue
			}

			if len(values) != len(expectedValues) {
				t.Errorf("Expected %d values for pattern %s, but got %d", len(expectedValues), pattern, len(values))
				continue
			}

			for i, expectedValue := range expectedValues {
				if values[i] != expectedValue {
					t.Errorf("Expected value %s, but got %s", expectedValue, values[i])
				}
			}
		}
	})

	t.Run("Test with multiple patterns", func(t *testing.T) {
		patterns := []string{"hello", "kenobi"}
		expectedResults := map[string][]string{
			"hello":  {"hello", "hello"},
			"kenobi": {"kenobi"},
		}

		results := s.applyRegexPatterns(texts, patterns)

		if len(results) != len(expectedResults) {
			t.Errorf("Expected %d results, but got %d", len(expectedResults), len(results))
		}

		for pattern, expectedValues := range expectedResults {
			values, ok := results[pattern]
			if !ok {
				t.Errorf("Expected pattern %s not found in results", pattern)
				continue
			}

			if len(values) != len(expectedValues) {
				t.Errorf("Expected %d values for pattern %s, but got %d", len(expectedValues), pattern, len(values))
				continue
			}

			for i, expectedValue := range expectedValues {
				if values[i] != expectedValue {
					t.Errorf("Expected value %s, but got %s", expectedValue, values[i])
				}
			}
		}
	})
}
