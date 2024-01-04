package main

import (
	"github.com/boltdb/bolt"
	"fmt"
	"bytes"
	"encoding/gob"
	"strings"

)

//Initialize database
func createDB() *bolt.DB {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}

	//Creates a bucket - a set of key value pairs.
	db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucket([]byte("testBucket"))
		if e != nil {
			fmt.Println(e)
			
			return e
		}
		return nil
	})

	return db

}


//Can double as an update function.
func writeDB(t Task) {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testBucket"))
		idBS, dBS := []byte{byte(t.id)}, encode(t.title, t.description, t.status)
		fmt.Println(idBS)
		e := b.Put(idBS, dBS)
		return e
	})
}

func getTask(id int)  Task{
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
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

func getTasks(f Filter) []Task {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}

	var result []Task
	viewErr := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testBucket"))
		b.ForEach(func(k, v []byte) error {
			//Logic where we check if it matches the filter
			d := decode(v)
			t := Task{
				title: d[0],
				description: d[1],
				status: d[2],
			}

			if strings.Contains(t.title, f.title){
				result = append(result, t)
				return nil
			} else if strings.Contains(t.description, f.description){
				result = append(result, t)
				return nil
			} else if strings.Contains(t.status, f.status){
				result = append(result, t)
				return nil
			}

			return fmt.Errorf("404: Nothing found.")
		})
		return nil
	})

	if viewErr != nil {
		fmt.Println(viewErr)
	}

	return result
}

func deleteDB(id int) {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
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
	e := enc.Encode([]string{title, description, status})
	fmt.Println(e)
	// if  e != nil {
	// 	log.Fatal(e)
	// 	return
	// }

	return buf.Bytes()
}

func decode(bs []byte) []string {
	buf := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(buf)

	var td []string
	e := dec.Decode(&td)
	fmt.Println(e)

	return td
}