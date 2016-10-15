package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/azert-software/virus-proxy/config"
)

func main() {
	path := flag.String("file", "", "A file path to scan")
	flag.Parse()

	cfg := &config.Config{}
	config.ReadConfig(cfg)

	v := NewVirusTotal(cfg.APIkey)

	w := setupLog()

	log.Printf("-------- STARTING SCAN %s -------", time.Now())
	log.Printf("File name %s ", *path)

	code := 0
	defer func() {
		log.Printf("-------- SCAN FINISHED %s -------", time.Now())
		w.Close()
		os.Exit(code)
	}()

	log.Println("Starting scan of file " + *path)

	_, h := v.GetFileBytesWithHash(*path)

	code = getAndCheckReport(v, h)

	if code != 0 {
		return
	}

	log.Println("File not known submitting file")

	v.ScanFile(*path)

	log.Println("Waiting on report")

	code = getAndCheckReport(v, h)
}

// sets up the standard go logger to write to a file
// this is setup as a rolling log with a max size around 1mb
func setupLog() *os.File {
	if _, err := os.Stat("log.txt"); os.IsNotExist(err) {
		os.Create("log.txt")
	}

	f, _ := os.Stat("log.txt")
	w, _ := os.OpenFile("log.txt", os.O_APPEND, 0666)

	if (f.Size()) > 1000000 {
		b, _ := ioutil.ReadFile("log.txt")
		bs := len(b)

		bSlice := make([]byte, 1000000)
		start := bs - 1000000
		iter := 0

		for i := start; i < bs; i++ {
			bSlice[iter] = b[i]
			iter++
		}

		ioutil.WriteFile("log.txt", bSlice, 0666)
	}

	log.SetOutput(w)

	return w
}

func getAndCheckReport(v *VirusTotal, h string) int {
	k, vir := v.GetReport(h, false)

	if k && vir {
		log.Println("File has virus")
		return -1
	}

	if k && !vir {
		log.Println("File has no virus")
		return 1
	}

	// file not known
	return 0
}
