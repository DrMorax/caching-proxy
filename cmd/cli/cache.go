package main

import "time"

type CacheObject struct {
	Status    int
	Headers   map[string][]string
	Response  []byte
	Timestamp time.Time
}

var cache = make(map[string]CacheObject)

func getFromCache(url string) (CacheObject, bool) {
	object, found := cache[url]
	return object, found
}

func storeInCache(url string, object CacheObject) {
	cache[url] = CacheObject{
		Status:    object.Status,
		Headers:   object.Headers,
		Response:  object.Response,
		Timestamp: time.Now(),
	}
}
