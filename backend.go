package main

import (
	"github.com/boltdb/bolt"
	"fmt"
	"bytes"
	"encoding/gob"

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



func writeDB(t Task) {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testBucket"))
		idBS, dBS := []byte{byte(t.id)}, encode(t.title, t.description, t.status)
		fmt.Println(idBS)
		err := b.Put(idBS, dBS)
		return err
	})
}


//converts string array/slice into byte slice, in order to be passed into writeDB().
//Alternative approach involves nesting buckets within the database, but I prefer this method.
func encode(title string, description string, status string) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode([]string{title, description, status})
	fmt.Println(err)
	// if  err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	return buf.Bytes()
}