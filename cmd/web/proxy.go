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
		w.Header().Set("X-Cache", "HIT")
		fmt.Println("HIT\t", key)
		responseWrite(w, object)
		return
	}

	w.Header().Set("X-Cache", "MISS")

	fmt.Println("MISS\t", key)
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
		Headers:   resp.Header,
		Response:  body,
		Timestamp: time.Now(),
	}
	storeInCache(r.Method+":"+query, newObject)

	responseWrite(w, newObject)
}

func responseWrite(w http.ResponseWriter, object CacheObject) {
	for header, values := range object.Headers {
		w.Header().Set(header, strings.Join(values, ", "))
	}
	w.WriteHeader(object.Status)
	w.Write(object.Response)
}
