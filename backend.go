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


//Can double as an update function.
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

func readDB(id int)  Task{
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
	}

	//Creates pointer so task can be referred to in anonymous function
	var t Task
	tPointer := &t

	viewErr := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testBucket"))
		v := b.Get([]byte{byte(id)})
		data := decode(v)

		(*tPointer).title = data[0]
		(*tPointer).description = data[1]
		(*tPointer).status = data[2]

		return nil
	})

	t.id = id

	if viewErr != nil {
		fmt.Println(viewErr)
	}
	

	return t
}

func deleteDB(id int) {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	deleteErr := db.Update(func(tx *bolt.Tx) error {
		//Deletes bucket item given id parameter in parent function
		return tx.Bucket([]byte("testBucket")).Delete([]byte{byte(id)})
	})

	if deleteErr != nil {
		fmt.Println(deleteErr)
	}
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

func decode(bs []byte) []string {
	buf := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(buf)

	var td []string
	err := dec.Decode(&td)
	fmt.Println(err)

	return td
}