package locale

import (
	"fmt"
	"testing"
)

func BenchmarkDetector(b *testing.B) {
	detector := NewDetector()
	testURLs := []string{
		"https://example.com/en/about",
		"https://es.example.com/productos",
		"https://example.com/page?lang=fr",
		"https://example.com/it/chi-siamo",
		"https://example.com/unique-endpoint",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		url := testURLs[i%len(testURLs)]
		_, _ = detector.Detect(url)
	}
}

func BenchmarkDetectorPathPrefix(b *testing.B) {
	detector := NewDetector()
	url := "https://example.com/en/about/page/123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.Detect(url)
	}
}

func BenchmarkDetectorSubdomain(b *testing.B) {
	detector := NewDetector()
	url := "https://en.example.com/about"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.Detect(url)
	}
}

func BenchmarkDetectorQueryParam(b *testing.B) {
	detector := NewDetector()
	url := "https://example.com/page?lang=en&foo=bar&baz=qux"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.Detect(url)
	}
}

func BenchmarkTranslationMatcher(b *testing.B) {
	matcher := NewTranslationMatcher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher.AreTranslations("about", "sobre-nosotros")
	}
}

func BenchmarkTranslationMatcherGetCanonical(b *testing.B) {
	matcher := NewTranslationMatcher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher.GetCanonical("sobre-nosotros")
	}
}

func BenchmarkGrouper(b *testing.B) {
	urls := []string{
		"https://example.com/about",
		"https://example.com/en/about",
		"https://example.com/es/sobre-nosotros",
		"https://example.com/it/chi-siamo",
		"https://example.com/products",
		"https://example.com/en/products",
		"https://example.com/es/productos",
		"https://example.com/contact",
		"https://example.com/en/contact",
		"https://example.com/es/contacto",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grouper := NewGrouper([]string{"en"})
		for _, url := range urls {
			_ = grouper.Add(url)
		}
		_ = grouper.GetBestURLs()
	}
}

func BenchmarkGrouperLargeScale(b *testing.B) {
	// Generate URLs
	locales := []string{"en", "es", "fr", "de", "it", "pt", "ja", "zh"}
	paths := []string{"about", "products", "contact", "services", "help", "privacy", "terms"}

	urls := make([]string, 0, len(locales)*len(paths))
	for _, locale := range locales {
		for _, path := range paths {
			urls = append(urls, fmt.Sprintf("https://example.com/%s/%s", locale, path))
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grouper := NewGrouper([]string{"en"})
		for _, url := range urls {
			_ = grouper.Add(url)
		}
		_ = grouper.GetBestURLs()
	}
}

func BenchmarkGrouperAdd(b *testing.B) {
	grouper := NewGrouper([]string{"en"})
	url := "https://example.com/en/about"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		grouper = NewGrouper([]string{"en"})
		b.StartTimer()
		_ = grouper.Add(url)
	}
}

func BenchmarkGrouperShouldGroup(b *testing.B) {
	grouper := NewGrouper([]string{"en"})
	url1 := "https://example.com/en/about"
	url2 := "https://example.com/es/sobre-nosotros"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = grouper.ShouldGroup(url1, url2)
	}
}

func BenchmarkScorer(b *testing.B) {
	scorer := NewScorer([]string{"en", "es", "fr"})
	detector := NewDetector()

	url := "https://example.com/en/about"
	localized, _ := detector.Detect(url)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = scorer.Score(localized, true)
	}
}

// Benchmark realistic workflow
func BenchmarkRealisticWorkflow(b *testing.B) {
	urls := []string{
		"https://example.com/about",
		"https://example.com/en/about",
		"https://example.com/es/sobre-nosotros",
		"https://example.com/it/chi-siamo",
		"https://example.com/fr/a-propos",
		"https://example.com/de/uber-uns",
		"https://example.com/products",
		"https://example.com/en/products",
		"https://example.com/es/productos",
		"https://example.com/it/prodotti",
		"https://example.com/contact",
		"https://example.com/en/contact",
		"https://example.com/es/contacto",
		"https://en.example.com/help",
		"https://es.example.com/ayuda",
		"https://example.com/page?lang=en",
		"https://example.com/page?lang=es",
		"https://example.com/unique-endpoint-1",
		"https://example.com/unique-endpoint-2",
		"https://example.com/api/v1/users",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grouper := NewGrouper([]string{"en"})
		for _, url := range urls {
			_ = grouper.Add(url)
		}
		results := grouper.GetBestURLs()
		_ = results
	}
}

// Benchmark memory allocation
func BenchmarkMemoryAllocation(b *testing.B) {
	urls := []string{
		"https://example.com/en/about",
		"https://example.com/es/sobre-nosotros",
		"https://example.com/it/chi-siamo",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		grouper := NewGrouper([]string{"en"})
		for _, url := range urls {
			_ = grouper.Add(url)
		}
	}
}

// Comparative benchmarks
func BenchmarkDetectorWithVsWithoutLocale(b *testing.B) {
	detector := NewDetector()

	b.Run("WithLocale", func(b *testing.B) {
		url := "https://example.com/en/about/page"
		for i := 0; i < b.N; i++ {
			_, _ = detector.Detect(url)
		}
	})

	b.Run("WithoutLocale", func(b *testing.B) {
		url := "https://example.com/about/page"
		for i := 0; i < b.N; i++ {
			_, _ = detector.Detect(url)
		}
	})
}
