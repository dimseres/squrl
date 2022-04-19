package database

// WatchLinks Что смотреть и к-во наблюдателей
// когда наблюдатели ноль стереть запись
type WatchLinks struct {
	Url        string
	WatchCount int
}

// ParsedData Хранит загруженные данные
type ParsedData struct {
	Url  string
	Data CrawlData
}

type CrawlData struct {
	Id               string
	GroupSubscribers int
}
