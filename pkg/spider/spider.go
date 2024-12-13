package spider

/*
Defines a Spider, which is a class that allows for the crawling of webpages and
the parsing of data into objects of a specific type. To use a Spider, a struct
implementing `Crawler` is defined and passed when creating the spider object.
*/
type Spider[T any] struct {
}
