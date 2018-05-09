package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/dkaps125/go-contract/contract"
)

type FileData struct {
	Chunks map[string][]byte
}

type MemStorage struct {
	Data map[string][]byte
}

var myAddr string
var c contract.Contract
var prevSubs map[uint32]struct{}
var memory MemStorage

func main() {
	fmt.Println("Running...")

	if len(os.Args) < 3 {
		panic("Expected contract address and account parameters")
	}

	prevSubs = make(map[uint32]struct{})
	memory.Data = make(map[string][]byte)

	addr := os.Args[1]
	myAddr = os.Args[2]
	c, _ = c.Init("contracts/Catalog.json", addr, "http://localhost:9545")

	var b [32]byte
	copy(b[:], "http://localhost:5437")

	r, _ := c.Transact("chunkServerJoin", myAddr, b)
	fmt.Printf("%v\n", r)

	go fileListener()

	n, _ := c.RegisterEventListener("SongPublished")
	c.Listen(n, "SongPublished", fn)
}

func addAll(req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var f FileData
	err := decoder.Decode(&f)

	if err != nil {
		panic(err)
	}

	defer req.Body.Close()

	for k, v := range f.Chunks {
		memory.Data[k] = v
		fmt.Printf("Storing chunk with hash %s\n", k)
	}
}

func getChunk(key string) []byte {
	return memory.Data[key]
}

func fileListener() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Write(getChunk(r.URL.Query().Get("hash")))
			break
		case "PUT":
			addAll(r)
			w.WriteHeader(http.StatusOK)
			break
		default:
			break
		}
	})

	http.ListenAndServe(":5437", nil)
}

func fn(data []interface{}) error {
	songIndex := data[0].(uint32)

	if _, ok := prevSubs[songIndex]; ok {
		return nil
	}

	ran := uint32(rand.Int31())

	c.Transact("chunkServerSubmitRandomness", myAddr, ran, songIndex)

	fmt.Printf("Submitted rand for %d\n", songIndex)

	prevSubs[songIndex] = struct{}{}

	return nil
}
