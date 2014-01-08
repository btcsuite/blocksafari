// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/conformal/btcjson"
	"github.com/davecgh/go-spew/spew"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	numMainPageBlocks = 20 // number of blocks to display on main page
)

func handleBlock(w http.ResponseWriter, r *http.Request) {
	blockhash := r.URL.Path[len("/block"):]
	if len(blockhash) < 2 || len(blockhash[1:]) != 64 {
                printErrorPage(w, "Invalid block hash")
                return
        }

	b, err := getBlock(blockhash[1:])
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	tx := make([]btcjson.TxRawResult, len(b.Tx))
	for i := range b.Tx {
		tx[i], err = getTx(b.Tx[i])
		if err != nil {
			break
		}
	}

	title := fmt.Sprintf("Block %v", b.Height)
	printHTMLHeader(w, title)
	printBlock(w, b, tx)
	printHTMLFooter(w)
}

func handleBlockNum(w http.ResponseWriter, r *http.Request) {
	blocknum := r.URL.Path[len("/b"):]
	if len(blocknum) < 2 {
		return
	}

	blocknum = blocknum[1:]
	b, err := strconv.Atoi(blocknum)
	if err != nil {
		fmt.Fprintf(w, "invalid block number: %v", blocknum)
		return
	}

	hash, err := getBlockHash(int64(b))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	uri := fmt.Sprintf("http://%v/block/%v", r.Host, hash)
	w.Header().Set("Location", uri)
	w.WriteHeader(307)
}

func handleCSS(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/css"):]
	if len(file) < 2 {
		return
	}

	http.ServeFile(w, r, "css/"+file[1:])
}

func handleJS(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/js"):]
	if len(file) < 2 {
		return
	}

	http.ServeFile(w, r, "js/"+file[1:])
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	nblocks, err := getBlockCount()
	if err != nil {
		printErrorPage(w, "Unable to get block count")
		return
	}

	blocks := make([]btcjson.BlockResult, numMainPageBlocks)
	for j := 0; j < numMainPageBlocks; j++ {
		cblock := int64(int(nblocks) - j)
		jstr, err := getBlockHash(cblock)
		if err != nil {
			printErrorPage(w, "Error retrieving block hash")
			return
		}
		blocks[j], err = getBlock(jstr)
		if err != nil {
			printErrorPage(w, "Error retrieving block")
			return
		}
	}

	printHTMLHeader(w, "Welcome")
	printMainBlock(w, blocks)
	printHTMLFooter(w)
}

func handleRawBlock(w http.ResponseWriter, r *http.Request) {
	block := r.URL.Path[len("/rawblock"):]
	if len(block) < 2 || len(block[1:]) != 64 {
		printErrorPage(w, "Invalid block hash")
		return
	}

	output, err := getRawBlock(block[1:])
	if err != nil {
		printErrorPage(w, "Block not found")
		return
	}

	printContentType(w, "text/plain")
	fmt.Fprintf(w, "%v", spew.Sdump(output))
}

func handleRawTx(w http.ResponseWriter, r *http.Request) {
	tx := r.URL.Path[len("/rawtx"):]
	if len(tx) < 2 || len(tx[1:]) != 64 {
		printErrorPage(w, "Invalid transaction id")
		return
	}

	output, err := getRawTx(tx[1:])
	if err != nil {
		printErrorPage(w, "Transaction not found")
		return
	}

	printContentType(w, "text/plain")
	fmt.Fprintf(w, "%v", spew.Sdump(output))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Path[len("/search"):]
	if len(q) < 2 {
		return
	}
	q = q[1:]
	if len(q) == 64 {
		uri := fmt.Sprintf("http://%v/block/%v", r.Host, q)
		w.Header().Set("Location", uri)
		w.WriteHeader(307)
	} else if _, err := strconv.ParseInt(q, 0, 64); err == nil {
		uri := fmt.Sprintf("http://%v/b/%v", r.Host, q)
		w.Header().Set("Location", uri)
		w.WriteHeader(307)
	} else {
		str := "Unknown search value: " + q
		printErrorPage(w, str)
	}
}

func handleTx(w http.ResponseWriter, r *http.Request) {
	tx := r.URL.Path[len("/tx"):]
	if len(tx) < 2 {
		printErrorPage(w, "Invalid TX hash")
		return
	}
	t, err := getTx(tx[1:])
	if err != nil {
		printErrorPage(w, "Unable to retrieve tx")
		return
	}

	title := fmt.Sprintf("Tx %v\n", t.Txid)
	printHTMLHeader(w, title)
	printTx(w, t)
	printHTMLFooter(w)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	page := strings.Split(path, "/")[0]
	switch page {
	case "b":
		handleBlockNum(w, r)
	case "block":
		handleBlock(w, r)
	case "css":
		handleCSS(w, r)
	case "js":
		handleJS(w, r)
	case "rawblock":
		handleRawBlock(w, r)
	case "rawtx":
		handleRawTx(w, r)
	case "search":
		handleSearch(w, r)
	case "tx":
		handleTx(w, r)
	case "":
		handleMain(w, r)
	default:
		/* XXX - serve 404's */
		fmt.Fprintf(w, "404 - Not found")
	}
}

func main() {
	var err error

	cfg, _, err = loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loadConfig failed: %v\n", err)
		os.Exit(-1)
	}

	listeners := make([]net.Listener, 0, len(cfg.Listeners))
	for _, addr := range cfg.Listeners {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Listen failed: %v\n", err)
			os.Exit(-1)
		}
		listeners = append(listeners, listener)
	}

	http.HandleFunc("/", handleRequest)

	httpServeMux := http.NewServeMux()
	httpServer := &http.Server{Handler: httpServeMux}
	httpServeMux.HandleFunc("/", handleRequest)

	var wg sync.WaitGroup
	for _, listener := range listeners {
		wg.Add(1)
		go func(listener net.Listener) {
			fmt.Fprintf(os.Stderr, "HTTP server listening on %s\n", listener.Addr())
			httpServer.Serve(listener)
			wg.Done()
		}(listener)
	}
	wg.Wait()
}
