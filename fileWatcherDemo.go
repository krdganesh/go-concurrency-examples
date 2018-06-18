package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const watchedPath = "./data"

// startfileWatcherDemo : Function watches for files inside folder named "data" which needs to be proccessed.
// Input Parameters
//         ch chan bool : boolen channel
func startfileWatcherDemo(ch chan bool) {
	ch <- true
	for {
		d, _ := os.Open(watchedPath)
		files, _ := d.Readdir(-1)
		for _, fi := range files {
			filePath := watchedPath + "/" + fi.Name()
			f, _ := os.Open(filePath)
			data, _ := ioutil.ReadAll(f)
			f.Close()
			if data != nil {
				os.Remove(filePath)
			}

			go func(data string) {
				fmt.Println("Record proccessed - ", data)
			}(string(data))
		}
	}
}
