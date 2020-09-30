package benchmark

import (
	"fmt"
	"github.com/gawwo/fake115-go/utils"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var wgp = utils.NewWaitGroupPool(50)

func PrintLocalDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 15 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}
	fmt.Println("connect done, use", conn.LocalAddr().String())

	return conn, err
}

var client = &http.Client{
	Transport: &http.Transport{
		Dial: PrintLocalDial,
	},
}

func doGet(url string, id int) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%d: %s -- %v\n", id, string(buf), err)
	if err := resp.Body.Close(); err != nil {
		fmt.Println(err)
	}
}

func doMultiRequest(url string) bool {
	defer wgp.Done()
	for i := 0; i < 500000; i++ {
		_, err := utils.Get(url, nil, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
	return true
}

// 并未发现有内存泄漏的问题
func DoMultiRequest() {
	url := "http://localhost:8441/"
	for i := 0; i < 50; i++ {
		wgp.Add()
		go doMultiRequest(url)
	}

	wgp.Wait()
}
