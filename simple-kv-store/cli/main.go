package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jacktantram/backend-go/simple-kv-store/kvstore"
)

func KVCLI(in io.Reader, w io.Writer) {
	store := kvstore.NewStore()
	scanner := bufio.NewScanner(in)
	for {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			continue
		}
		splitStr := strings.Split(text, " ")
		command := strings.ToUpper(splitStr[0])
		switch command {
		case "WRITE":
			if len(splitStr) <= 2 {
				_, _ = fmt.Fprintln(w, "missing arguments need <key>, <val>")
				continue
			}
			if err := store.Set(splitStr[1], splitStr[2]); err != nil {
				_, _ = fmt.Fprintf(w, "unable to create: %v\n", err)
			}
		case "READ":
			if len(splitStr) <= 1 {
				_, _ = fmt.Fprintln(w, "missing arguments need <key>")
				continue
			}
			val, ok := store.Get(splitStr[1])
			if !ok {
				_, _ = fmt.Fprintf(w, "Key not found: %s\n", splitStr[1])
				continue
			}
			_, _ = fmt.Fprintln(w, val)
		case "DELETE":
			if len(splitStr) <= 1 {
				_, _ = fmt.Fprintln(w, "missing arguments need <key>")
				continue
			}
			store.Delete(splitStr[1])
		case "START":
			store.Begin()
		case "COMMIT":
			if err := store.Commit(); err != nil {
				_, _ = fmt.Fprintf(w, "unable to commit transaction: %v\n", err)
			}
		case "ABORT":
			if err := store.Abort(); err != nil {
				_, _ = fmt.Fprintf(w, "unable to abort transaction: %v\n", err)
			}
		case "QUIT":
			_, _ = fmt.Fprintln(w, "Exiting...")
			return
		}
	}
}

func main() {
	KVCLI(os.Stdin, os.Stdout)
	os.Exit(0)
}
