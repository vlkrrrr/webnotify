package internal

import (
	"bufio"
	"fmt"
	"net/http"
)

func PrintEvents() {
	c := http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/events", nil)

	resp, _ := c.Do(req)
	r := bufio.NewReader(resp.Body)
	res, _ := r.ReadString('\n')

	fmt.Println("event", string(res))
}
