package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Node struct {
	data int
	next *Node
}

func newNode(val int) *Node  {
	var n Node
	n.data = val
	n.next = nil
	return &n
}

func addNode(head *Node, val int) *Node {
	tptr := head
	head = newNode(val)
	head.next = tptr
	return head
}

func printList(head *Node) {
	for head != nil {
		println(head.data)
		head = head.next
	}
}

func fetch(url string, ch chan <- string){
	resp,err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error in http get %d",resp)
		return
	}
	nbytes,err := io.Copy(ioutil.Discard, resp.Body)

	defer resp.Body.Close()

	if err != nil {
		ch <- fmt.Sprintf("Error in copying body %s: %v", url, err)
		return
	}

	ch <- fmt.Sprintf("%s - Body Size in bytes = %d",url,nbytes)
	return
}

func main() {
	ch := make(chan string)
	for _,url := range os.Args[1:] {
		go fetch(url,ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
}

func myvariadic(multi ...int){
	for index,val := range multi{
		fmt.Println(index,val)
	}
}