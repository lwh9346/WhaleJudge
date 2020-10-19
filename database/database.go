package database

import (
	"log"

	"github.com/xujiajun/nutsdb"
)

//HasKey 返回一个k-v型bucket是否含有某个key
func HasKey(db *nutsdb.DB, bucket, key string) bool {
	err := db.View(
		func(tx *nutsdb.Tx) error {
			_, err1 := tx.Get(bucket, []byte(key))
			return err1
		})
	if err == nil {
		return true
	}
	if err == nutsdb.ErrKeyNotFound {
		return false
	}
	log.Fatal(err)
	return false
}

//SetValue 设置k-v型bucket的键值
func SetValue(db *nutsdb.DB, bucket, key string, value []byte, ttl uint32) {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			return tx.Put(bucket, []byte(key), value, ttl)
		}); err != nil {
		log.Fatal(err)
	}
}

//GetValue 获取k-v型bucket的键值
func GetValue(db *nutsdb.DB, bucket, key string) []byte {
	var data []byte
	err := db.View(
		func(tx *nutsdb.Tx) error {
			d, err1 := tx.Get(bucket, []byte(key))
			if err1 != nil {
				data = d.Value
			}
			return err1
		})
	if err != nil {
		log.Fatal(err)
	}
	return data
}
