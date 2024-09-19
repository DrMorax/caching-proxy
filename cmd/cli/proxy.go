package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func Proxy(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("url")
	key := r.Method + ":" + query

	object, found := getFromCache(key)
	if found {
		fmt.Printf("Request '%s' is cached\n", key)
		responseWrite(w, r, object, "HIT")
		return
	}

	fmt.Println("Cache Not Present for key : ", key)
	req, err := http.NewRequest(r.Method, query, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to make request to url '%s' \nError: %s", req.URL, err)
		http.Error(w, "Failed to make request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to parse the response from the targe url", http.StatusInternalServerError)
		return
	}

	newObject := CacheObject{
		Status:    resp.StatusCode,
		Headers:   req.Header,
		Response:  body,
		Timestamp: time.Now(),
	}
	storeInCache(r.Method+":"+query, newObject)

	responseWrite(w, r, newObject, "MISS")

	for key, _ := range cache {
		fmt.Println(key)
	}
}

func responseWrite(w http.ResponseWriter, r *http.Request, object CacheObject, result string) {
	r.Header.Set("X-Cache", result)
	for header, values := range object.Headers {
		r.Header.Set(header, strings.Join(values, "; "))
	}
	w.WriteHeader(object.Status)
	w.Write(object.Response)
}
