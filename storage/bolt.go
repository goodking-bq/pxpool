package storage

import bolt "go.etcd.io/bbolt"

type Bolt struct {
	db   *bolt.DB
	path string
}
