package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/alicebob/miniredis"
	yaml "gopkg.in/yaml.v3"
)

///////
type (
	_ = yaml.Encoder

	PersistedHash = map[string]string
	PersistedList = []string
	PersistedSet  = []string
	PersistedZSet = map[float64]string

	PersistedKey struct {
		Type  string
		Key   string
		Value interface{}
	}

	PersistedDB struct {
		// DBNum int
		Keys []PersistedKey
	}
)

func (pdb *PersistedDB) DumpYML(f io.Writer, indent int) error {
	enc := yaml.NewEncoder(f)
	enc.SetIndent(indent)
	defer enc.Close()

	if err := enc.Encode(pdb); err != nil {
		log.Println("Failed to marshal into YAML:", err)
		return err
	}
	return nil
}

///////
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
	fmt.Fprintln(f, "---")
	for _, k := range db.Keys() {
		t := db.Type(k)
		switch t {
		case "string":
			v, _ := db.Get(k)
			fmt.Fprintf(f, "%s: %s\n", k, v)
			continue
		default:
			fmt.Fprintf(f, "%s:\n", k)
		}
		switch t {
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

func (db *RedisDB) Load() error {
	data, err := ioutil.ReadFile(db.filePath)
	if err != nil {
		return err
	}

	dbData := map[string]interface{}{}
	if err := yaml.Unmarshal(data, dbData); err != nil {
		return err
	}

	spew.Dump(dbData)
	return nil
}

func main() {
	//// --
	var (
		srvAddr       string
		persistPeriod time.Duration
		dbsNums       IntSliceFlag
	)
	flag.StringVar(&srvAddr, "listen", "127.0.0.1:6379", "Listen address")
	flag.DurationVar(&persistPeriod, "period", 30*time.Second, "DB store period")
	flag.Var(&dbsNums, "db", "DB number to store (multiple DBs are allowed)")
	flag.Parse()

	//// --
	srv := miniredis.NewMiniRedis()
	defer srv.Close()

	//// --
	dbs := []*RedisDB{}
	for _, dbNum := range dbsNums {
		dbs = append(dbs, WrapRedisDB(srv, dbNum, fmt.Sprintf("redis_db%d.yaml", dbNum)))
	}

	//// --
	loadAll := func() {
		for _, db := range dbs {
			if err := db.Load(); err != nil {
				fmt.Printf("{WARN} no persisted db: %v\n", err)
			} else {
				fmt.Println(">> DB loaded:", db.filePath)
			}
		}
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
	//defer persistAll()

	//// --
	loadAll()

	if err := srv.StartAddr(srvAddr); err != nil {
		fmt.Println("Unable to start miniredis:", err)
		return
	}

	for true {
		time.Sleep(persistPeriod)
		persistAll()
	}
}
