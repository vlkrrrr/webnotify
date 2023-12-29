package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
)

func PrintEvents() {
	c := http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/events", nil)

	resp, _ := c.Do(req)
	scn := bufio.NewScanner(resp.Body)
	scn.Split(scanEvents)
	for scn.Scan() {
		fmt.Println(scn.Text())
	}
	if err := scn.Err(); err != nil {
		fmt.Println("error reading events:", err)
	}

}

func scanEvents(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
