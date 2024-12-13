package spider

/*
Defines the basic layout of a crawler plugin for the `Spider`API. A crawler has
the following methods:
- `Crawl(url)`: fetches the raw HTML that is to be parsed
- `Aggregate()`: parses the raw HTML into a desired type
*/
type Crawler[T any] interface {
	// Fetches the raw HTML that is to be parsed.
	Crawl(url string) (string, error)

	// Parses the raw HTML into a desired type.
	Aggregate() (*T, error)
}
