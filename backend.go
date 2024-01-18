package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
	// "strconv"
)

// Initialize database
func createDB() *bolt.DB {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}
	defer db.Close()
	//Creates a bucket - a set of key value pairs.
	db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte("testBucket"))
		if e != nil {
			fmt.Println(e)

			return e
		}
		return nil
	})

	return db

}

// Can double as an update function.
func writeDB(t Task) {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("testBucket"))
		IdBS, dBS := []byte{byte(t.Id)}, encode(t.Title, t.Description, t.Status)
		fmt.Println("!!!!!!", IdBS)
		e := b.Put(IdBS, dBS)
		return e
	})
}

func getTaskDB(Id int) (Task, error) {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}
	defer db.Close()

	//Creates pointer so task can be referred to in anonymous function
	var t Task
	tPointer := &t

	viewErr := db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("testBucket"))
		v := b.Get([]byte{byte(Id)})
		fmt.Printf("Key: %d, Value: %v\n", Id, v)
		if v == nil {
			return errors.New("Key not found")
		}

		data := decode(v)
		if len(data) <= 0 {
			return errors.New("404")
		}
		(*tPointer).Title = data[0]
		(*tPointer).Description = data[1]
		(*tPointer).Status = data[2]

		return nil
	})

	t.Id = Id

	if viewErr != nil {
		return t, viewErr
	}

	return t, viewErr
}

func getTasksDB(f Filter) []Task {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}

	defer db.Close()

	var result []Task
	viewErr := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("testBucket"))
		if err != nil {
			panic(err)
		}
		matchFound := false
		b.ForEach(func(k, v []byte) error {
			//Logic where we check if it matches the filter
			fmt.Println("key: ", k)
			id := int(k[0])
			fmt.Println("the id", id)
			d := decode(v)
			t := Task{
				Id:          id,
				Title:       d[0],
				Description: d[1],
				Status:      d[2],
			}

			if strings.Contains(t.Title, f.Title) || strings.Contains(t.Description, f.Description) || strings.Contains(t.Status, f.Status) {
				fmt.Println("hello world", t.Title)
				result = append(result, t)
				matchFound = true
			} else {
				fmt.Println(f.Title, t.Title)
			}
			return nil

		})

		if !matchFound {
			return fmt.Errorf("404: Nothing found.")
		}
		return nil
	})

	if viewErr != nil {
		fmt.Println(viewErr)
	}
	fmt.Println(result)

	return result
}

func deleteTaskDB(Id int) error {
	db, e := bolt.Open("test.db", 0600, nil)
	if e != nil {
		fmt.Println(e)
	}
	defer db.Close()
	deleteErr := db.Update(func(tx *bolt.Tx) error {
		//Deletes bucket item given Id parameter in parent function
		b, e := tx.CreateBucketIfNotExists([]byte("testBucket"))
		b.Delete([]byte{byte(Id)})
		return e
	})

	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

// converts string array/slice into byte slice, in order to be passed into writeDB().
// Alternative approach involves nesting buckets within the database, but I prefer this method.
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
