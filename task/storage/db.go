package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/st0zy/gophercises/task/pkg/adding"
	"github.com/st0zy/gophercises/task/pkg/listing"
)

func openDatabase() (*bolt.DB, error) {

	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: time.Second * 1})
	if err != nil {
		panic(err)
	}
	return db, nil

}

func Init() (*bolt.DB, error) {
	db, err := openDatabase()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		return err
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AddTask(task adding.Task) error {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: time.Second * 5})

	if err != nil {
		return errors.New("failed to open db connection")
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return errors.New("failed to retrieve bucket")
		}
		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		task := Task{
			Id:        id,
			Name:      task.Name,
			Completed: false,
		}
		var buffer bytes.Buffer
		json.NewEncoder(&buffer).Encode(task)
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, id)
		b.Put(bs, buffer.Bytes())
		return err
	})

	return err

}

func GetTasks() ([]listing.Task, error) {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: time.Second * 5})
	if err != nil {
		return nil, errors.New("failed to open db connection")
	}
	defer db.Close()

	var tasks []listing.Task
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return errors.New("failed to retrieve bucket")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task listing.Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				fmt.Println(err)
				continue
			}
			tasks = append(tasks, task)
		}
		return nil

	})
	return tasks, nil
}

func DoTask(taskId uint64) error {

	db, err := bolt.Open("my.db", 0666, &bolt.Options{Timeout: time.Second * 3})
	if err != nil {
		return errors.New("failed to open db connection")
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return errors.New("failed to locate the tasks bucket")
		}
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, taskId)
		var task listing.Task
		err := json.Unmarshal(b.Get(bs), &task)
		if err != nil {
			return err
		}
		task.Completed = true
		marshalledTask, _ := json.Marshal(task)
		err = b.Put(bs, marshalledTask)
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}
