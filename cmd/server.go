package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type Service struct {
	cache Cache
}

type Cache struct {
	Entries []Entry
}

type Entry struct {
	UUID string `json:"uuid"`
	TTL  int    `json:"ttl"`
	Body string `json:"body"`
}

func (s *Service) add(new Entry) {
	s.cache.Entries = append(s.cache.Entries, new)
}

func (s *Service) remove(this Entry) {
	for idx, entry := range s.cache.Entries {
		if entry.UUID == this.UUID {
			s.cache.Entries = append(s.cache.Entries[:idx], s.cache.Entries[idx+1:]...)
		}
	}
}

func (s *Service) entriesHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"entries": s.cache.Entries,
	})
}

func (s *Service) newEntryHandler(c *gin.Context) {
	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	s.add(Entry{UUID: id.String(), TTL: 1000, Body: ""})

	c.JSON(200, gin.H{
		"uuid": id.String(),
	})
}

func newService() *Service {
	return &Service{cache: Cache{}}
}

func main() {
	r := gin.Default()

	svc := newService()

	// go doEvery(10*time.Second, evictExpired)

	r.GET("/entries", svc.entriesHandler)
	r.POST("/new", svc.newEntryHandler)

	r.Run("localhost:8080")
}

// ----------------------

// func doEvery(d time.Duration, f func(time.Time)) {
// 	for x := range time.Tick(d) {
// 		f(x)
// 	}
// }

// func evictExpired(now time.Time) {
// 	fmt.Printf("evicting expired entries at %v", now)
// 	for idx, entry := range localCache.Entries {
// 		if entry.TTL == 0 {
// 			localCache.Entries = append(localCache.Entries[:idx], localCache.Entries[idx+1:]...)
// 		}
// 	}
// }
