package main

func main() {
	ch := make(chan bool)
	go startfileWatcherDemo(ch)
	<-ch
	for {
		ch1 := make(chan bool)
		go startAsyncWebServiceDemo(ch1)
		<-ch1
		//Uncomment below line to add delay of 10 secs. and adjust it as per your need
		// time.Sleep(10000 * time.Millisecond)
	}
}
