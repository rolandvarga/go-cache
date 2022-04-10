package main

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type Service struct {
	cache             Cache
	evictionFrequency time.Duration
}

type Cache struct {
	mu      sync.Mutex
	Entries []Entry
}

type Entry struct {
	UUID      string    `json:"uuid"`
	TTL       int       `json:"ttl"` // milliseconds
	Body      string    `json:"body" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *Entry) isExpired() bool {
	return time.Now().UTC().After(e.CreatedAt.Add(time.Duration(e.TTL * int(time.Millisecond))))
}

func (c *Cache) add(new Entry) {
	c.mu.Lock()
	c.Entries = append(c.Entries, new)
	c.mu.Unlock()
}

func (c *Cache) remove(this Entry) {
	for idx, entry := range c.Entries {
		if entry.UUID == this.UUID {
			c.mu.Lock()
			c.Entries = append(c.Entries[:idx], c.Entries[idx+1:]...)
			c.mu.Unlock()
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

	var entry Entry
	err = c.BindJSON(&entry)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		log.Println(err)
		return
	}

	ttl := 10000 // default value
	if entry.TTL > 0 {
		ttl = entry.TTL
	}
	s.cache.add(newEntry(id.String(), ttl, entry.Body))

	c.JSON(200, gin.H{
		"uuid": id.String(),
	})
}

func (s *Service) expireEntries(expiredCh chan int) {
	t := time.NewTicker(s.evictionFrequency * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			count := 0
			for _, entry := range s.cache.Entries {
				if entry.isExpired() {
					s.cache.remove(entry)
					count++
				}
			}
			expiredCh <- count
		}
	}
}

func newEntry(uuid string, ttl int, body string) Entry {
	return Entry{UUID: uuid, TTL: ttl, Body: body, CreatedAt: time.Now().UTC()}
}

func newService() *Service {
	return &Service{cache: Cache{}, evictionFrequency: time.Duration(10)}
}

func main() {
	r := gin.Default()

	svc := newService()

	r.GET("/entries", svc.entriesHandler)
	r.POST("/new", svc.newEntryHandler)

	expiredCh := make(chan int)
	go svc.expireEntries(expiredCh)

	go func() {
		for {
			select {
			case expirecount := <-expiredCh:
				log.Printf("expired '%d' entries", expirecount)
			}
		}
	}()

	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
