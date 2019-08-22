package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alicebob/miniredis"
)

type RedisDB struct {
	*miniredis.RedisDB
	filePath string
}

func WrapRedisDB(rd *miniredis.Miniredis, dbNum int, dumpPath string) *RedisDB {
	return &RedisDB{
		RedisDB:  rd.DB(dbNum),
		filePath: dumpPath,
	}
}

func (db *RedisDB) Persist() error {
	if db.filePath == "" {
		return errors.New("persist: no file path has been specified")
	}
	f, err := os.Create(db.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	indent := "  "
	for _, k := range db.Keys() {
		fmt.Fprintf(f, "- %s\n", k)
		t := db.Type(k)
		switch t {
		case "string":
			v, _ := db.Get(k)
			fmt.Fprintf(f, "%s%s\n", indent, v)
		case "hash":
			fields, _ := db.HKeys(k)
			for _, hk := range fields {
				fmt.Fprintf(f, "%s%s: %s\n", indent, hk, db.HGet(k, hk))
			}
		case "list":
			items, _ := db.List(k)
			for _, lk := range items {
				fmt.Fprintf(f, "%s%s\n", indent, lk)
			}
		case "set":
			membs, _ := db.Members(k)
			for _, mk := range membs {
				fmt.Fprintf(f, "%s%s\n", indent, mk)
			}
		case "zset":
			zmembs, _ := db.ZMembers(k)
			for _, zm := range zmembs {
				score, _ := db.ZScore(k, zm)
				fmt.Fprintf(f, "%s%f: %s\n", indent, score, zm)
			}
		default:
			fmt.Fprintf(os.Stderr, "FIXME: %s(a %s, fixme!)\n", indent, t)
		}
	}
	return nil
}

func main() {
	srv := miniredis.NewMiniRedis()
	defer srv.Close()

	dbs := []*RedisDB{
		WrapRedisDB(srv, 0, "redis_db0.txt"),
		WrapRedisDB(srv, 10, "redis_db10.txt"),
	}
	persistAll := func() {
		for _, db := range dbs {
			if err := db.Persist(); err != nil {
				fmt.Printf("{ERR} failed to persist db: %v\n", err)
			} else {
				fmt.Println(">> DB persisted:", db.filePath)
			}
		}
	}
	defer persistAll()

	if err := srv.StartAddr("127.0.0.1:6379"); err != nil {
		fmt.Println("Unable to start miniredis:", err)
		return
	}

	//db0, db10 := srv.DB(0), srv.DB(10)
	for true {
		time.Sleep(2 * time.Minute)
		persistAll()
	}
}
