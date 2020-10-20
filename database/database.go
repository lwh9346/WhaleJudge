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
	return err == nil
}

//RemoveKey 删除一个key
func RemoveKey(db *nutsdb.DB, bucket, key string) {
	if !HasKey(db, bucket, key) {
		return
	}
	db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(bucket, []byte(key))
	})
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
			if err1 == nil {
				data = d.Value
			}
			return err1
		})
	if err != nil {
		log.Fatal(err)
	}
	return data
}

//SAdd 向集合型bucket中添加元素
func SAdd(db *nutsdb.DB, bucket, key string, value []byte) {
	err := db.Update(func(tx *nutsdb.Tx) error {
		return tx.SAdd(bucket, []byte(key), value)
	})
	if err != nil {
		log.Fatal(err)
	}
}

//SAllMembers 取出集合型bucket中某个key的所有元素，返回元素与成功与否
func SAllMembers(db *nutsdb.DB, bucket, key string) ([][]byte, bool) {
	var data [][]byte
	err := db.View(func(tx *nutsdb.Tx) error {
		members, err1 := tx.SMembers(bucket, []byte(key))
		data = members
		return err1
	})
	return data, err == nil
}

//SRemove 移除集合型bucket中某个key下的指定元素
func SRemove(db *nutsdb.DB, bucket, key string, value []byte) {
	if !SHasKey(db, bucket, key) {
		return
	}
	db.Update(func(tx *nutsdb.Tx) error {
		return tx.SRem(bucket, []byte(key), value)
	})
}

//SHasKey 返回集合型bucket中是否含有某个key
func SHasKey(db *nutsdb.DB, bucket, key string) bool {
	var has bool
	db.View(func(tx *nutsdb.Tx) error {
		ok, err := tx.SHasKey(bucket, []byte(key))
		has = ok
		return err
	})
	return has
}

//SIsMember 返回集合型bucket中某个key下是否有指定元素
func SIsMember(db *nutsdb.DB, bucket, key string, value []byte) bool {
	var is bool
	db.View(func(tx *nutsdb.Tx) error {
		is, _ = tx.SAreMembers(bucket, []byte(key), value)
		return nil
	})
	return is
}
