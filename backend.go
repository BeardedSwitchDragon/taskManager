package main

import (
	"github.com/boltdb/bolt"
	"fmt"
)

//Initialize database
func createDB() *bolt.DB {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
	}

	//Creates a bucket - a set of key value pairs.
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("testBucket"))
		if err != nil {
			fmt.Println(err)
			
			return err
		}
		return nil
	})

	return db

}

func writeDB(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testBucket"))
	})
}

