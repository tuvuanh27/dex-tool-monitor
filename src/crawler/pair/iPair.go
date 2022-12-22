package pair

type IPairCrawler interface {
	GetFactoryAddress() string
	GetToken() (string, string)
}
