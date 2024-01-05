package main

import (
	"github.com/boltdb/bolt"
	"fmt"
	"bytes"
	"encoding/gob"
	"strings"
	"strconv"

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
		IdBS, dBS := []byte{byte(t.Id)}, encode(t.Title, t.Description, t.Status)
		fmt.Println(IdBS)
		e := b.Put(IdBS, dBS)
		return e
	})
}

func getTask(Id int)  Task{
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}

	//Creates pointer so task can be referred to in anonymous function
	var t Task
	tPointer := &t

	viewErr := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testBucket"))
		v := b.Get([]byte{byte(Id)})
		data := decode(v)

		(*tPointer).Title = data[0]
		(*tPointer).Description = data[1]
		(*tPointer).Status = data[2]

		return nil
	})

	t.Id = Id

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

	defer db.Close()

	var result []Task
	viewErr := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testBucket"))
		b.ForEach(func(k, v []byte) error {
			//Logic where we check if it matches the filter
			id, _ := strconv.Atoi(string(k))
			d := decode(v)
			t := Task{
				Id: id,
				Title: d[0],
				Description: d[1],
				Status: d[2],
			}

			if strings.Contains(t.Title, f.Title) || strings.Contains(t.Description, f.Description) || strings.Contains(t.Status, f.Status){
				fmt.Println(result)
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
	fmt.Println(result)

	return result
}

func deleteTask(Id int) error {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}
	deleteErr := db.Update(func(tx *bolt.Tx) error {
		//Deletes bucket item given Id parameter in parent function
		return tx.Bucket([]byte("testBucket")).Delete([]byte{byte(Id)})
	})

	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

//converts string array/slice into byte slice, in order to be passed into writeDB().
//Alternative approach involves nesting buckets within the database, but I prefer this method.
func encode(Title string, Description string, Status string) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	e := enc.Encode([]string{Title, Description, Status})
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