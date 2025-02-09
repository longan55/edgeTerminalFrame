package global

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"go.uber.org/zap"
)

var (
	//Bucket
	BucketFilenames = "filenames"
	BucketDevices   = "cnc_devices"
	//
	BucketDevice    = "device"
	BucketReport    = "report"
	BucketProtocol  = "protocol"
	BucketModel     = "model"
	BucketModelItem = "modelItem"
	BucketOther     = "other"
)

const runmode1 = "_only_PLATFORMMODE"

var Mydb *mydb

type mydb struct {
	boltdb *bolt.DB
}

func InitBoltdb() {
	//网关模式
	db, err := bolt.Open("./my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Panic(err)
	}
	Mydb = &mydb{}
	Mydb.boltdb = db

	buckets := []string{
		BucketFilenames,
		BucketDevices,
		BucketDevice,
		BucketReport,
		BucketProtocol,
		BucketModel,
		BucketModelItem,
		BucketOther,
	}

	for _, bucket := range buckets {
		if err := Mydb.boltdb.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}); err != nil {
			Logger.Error("数据桶创建失败", zap.String("bucket", bucket), zap.Error(err))
			continue
		}
		Logger.Info("数据桶创建成功", zap.String("bucket", bucket))
	}

	RegisterQuitTask(func() error {
		if Mydb != nil && Mydb.boltdb != nil {
			// stats := Mydb.boltdb.Stats()
			// if stats.TxStats.Open > 0 {
			// 	log.Println("There are open transactions!")
			// }
			return Mydb.boltdb.Close()
		}
		return errors.New("boltdb为空或已关闭")
	}, "关闭boltdb")
}

func (my mydb) CreateBucketIfNotExists(bucketname string) error {
	return my.boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		return err
	})
}

func (my mydb) Select(bucket, key string) (result []byte, err error) {
	return result, my.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		result = b.Get([]byte(key))
		return nil
	})
}
func (my mydb) SelectAll(bucket string) (result [][]byte, err error) {
	return result, my.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			result = append(result, v)
		}
		return nil
	})
}
func (my mydb) UpdateOne(bucket, key string, value any) (err error) {
	var bytes []byte
	bytes, ok := value.([]byte)
	if !ok {
		bytes, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("保存数据失败:%w", err)
		}
	}
	return my.boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(key), bytes)
	})
}

func (my mydb) DeleteOne(bucket, key string) error {
	return my.boltdb.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(bucket)); b == nil {
			return bolt.ErrBucketNotFound
		} else {
			return b.Delete([]byte(key))
		}
	})
}

func SelectOne(bucket, key string) (result []byte, err error) {
	return result, Mydb.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		result = b.Get([]byte(key))
		return nil
	})
}

func SelectAll(bucket string) (result [][]byte, err error) {
	return result, Mydb.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			result = append(result, v)
		}
		return nil
	})
}

// UpdateOne 创建或更新
func UpdateOne(bucket, key string, value any) (err error) {
	var bytes []byte
	bytes, ok := value.([]byte)
	if !ok {
		bytes, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("保存数据失败:%w", err)
		}
	}
	return Mydb.boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(key), bytes)
	})
}

type DataItem struct {
	Bucket string
	Key    string
	Val    any
}

func DeleteOne(bucket, key string) error {
	return Mydb.boltdb.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(bucket)); b == nil {
			return bolt.ErrBucketNotFound
		} else {
			return b.Delete([]byte(key))
		}
	})
}
