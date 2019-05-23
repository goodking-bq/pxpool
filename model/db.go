package model

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

// Story 存储
type Story struct {
	db   *bolt.DB
	path string
}

func NewStory(filepath string) *Story {
	db, err := bolt.Open(filepath+"/pxpool.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Story{db: db, path: "."}
}

func GetStory() *Story {
	db, err := bolt.Open("."+"/pxpool.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Story{db: db, path: "."}
}

func (story *Story) Close() {
	story.db.Close()
}

func (story *Story) ShowAll(b string) {
	if err := story.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(b))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key = %s,value = %s\n", k, v)
		}
		return nil
	}); err != nil {
		log.Println(err)
	}
}

// Storage 村粗
