package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pin/tftp/v3"
)

func main() {
	s := tftp.NewServer(readHandler, nil)
	err := s.ListenAndServe(":69")
	if err != nil {
		log.Fatal(err)
	}
}

func readHandler(filename string, rf io.ReaderFrom) error {
	err := func() error {
		log.Printf("Received TFTP request for: %q\n", filename)

		filename = strings.TrimPrefix(filename, ":")

		resp, err := http.Get(filename)
		if err != nil {
			return fmt.Errorf("HTTP GET %q failed: %w\n", filename, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP GET %q failed with status: %q", filename, resp.Status)
		}

		rf.(tftp.OutgoingTransfer).SetSize(resp.ContentLength)
		n, err := rf.ReadFrom(resp.Body)
		if err != nil {
			return fmt.Errorf("Error sending data to TFTP client: %w\n", err)
		}

		log.Printf("Transferred %d bytes from %q\n", n, filename)
		return nil
	}()
	if err != nil {
		log.Printf("Error processing TFTP request for %q: %v\n", filename, err)
	}
	return err
}
