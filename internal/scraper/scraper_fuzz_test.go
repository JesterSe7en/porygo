package scraper

import (
	"testing"
)

func FuzzApplySelectors(f *testing.F) {
	// Seed corpus with the HTML from the unit test
	html := `
	<html>
		<body>
			<h1>Hello, World!</h1>
			<a href="/page1">Page 1</a>
			<a href="/page2">Page 2</a>
			<a href="http://www.google.com">Google</a>
		</body>
	</html>
	`
	f.Add([]byte(html))

	s := &Scraper{}
	selectors := []string{"h1", "a@href"}

	f.Fuzz(func(t *testing.T, data []byte) {
		// The fuzz engine will generate random 'data'.
		// The goal is to ensure that applySelectors does not panic on any input.
		// We don't need to check the correctness of the results, just that it runs without crashing.
		s.applySelectors(data, selectors)
	})
}
