package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Entry struct {
	ID  string `json:"id"`
	TTL int    `json:"ttl"`
}

type Cache struct {
	entries []Entry
}

var localCache = Cache{}

func entriesHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"entries": localCache.entries,
	})
}

func newEntryHandler(c *gin.Context) {
	var entry Entry

	if err := json.NewDecoder(c.Request.Body).Decode(&entry); err != nil {
		c.JSON(400, gin.H{
			"error": "bad request",
		})
		return
	}
	localCache.entries = append(localCache.entries, entry)
	c.JSON(200, gin.H{
		"message": "success",
	})
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func evictExpired(now time.Time) {
	fmt.Printf("evicting expired entries at %v", now)
	for idx, entry := range localCache.entries {
		if entry.TTL == 0 {
			localCache.entries = append(localCache.entries[:idx], localCache.entries[idx+1:]...)
		}
	}
}

func main() {
	r := gin.Default()

	go doEvery(10*time.Second, evictExpired)

	r.GET("/entries", entriesHandler)
	r.POST("/new", newEntryHandler)

	r.Run("localhost:8080")
}
