package cleaner

import (
	"context"
	"log"
)

type Cleaner struct {
	ctx     context.Context
	cleanCh <-chan string
	s       Storage
}

type Storage interface {
	Delete(filename string) error
}

func NewCleaner(ctx context.Context, c <-chan string, s Storage) *Cleaner {
	return &Cleaner{
		ctx:     ctx,
		cleanCh: c,
		s:       s,
	}
}

func (c *Cleaner) Start() {
	log.Println("Starting cleaner")

	for {
		select {
		case <-c.ctx.Done():
			log.Println("Stop cleaner by done channel")
			break
		case filename := <-c.cleanCh:
			log.Println("Delete old img")
			if err := c.s.Delete(filename); err != nil {
				log.Println("Delete file err ", err)
			}
		}
	}
}
