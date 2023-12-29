package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func PrintEvents() {
	c := http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/events", nil)

	getEvents := func() error {
		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		scn := bufio.NewScanner(resp.Body)
		scn.Split(scanEvents)

		for scn.Scan() {
			fmt.Println(scn.Text())
		}
		if err := scn.Err(); err != nil {
			return err
		}
		return nil
	}

	retry(getEvents)
}

func retry(f func() error) {
	bo := 2
	ty := 3
	at := 0
	for at <= ty {
		err := f()
		if err != nil {
			fmt.Println("Error occured: ", err)
			fmt.Println("Retrying")
			time.Sleep(time.Duration(at*bo) * time.Second)
			at++
		} else {
			break
		}
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
