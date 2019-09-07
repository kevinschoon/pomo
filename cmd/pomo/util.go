package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
)

func maybe(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseTags(kvs []string) (map[string]string, error) {
	tags := map[string]string{}
	for _, kv := range kvs {
		split := strings.Split(kv, "=")
		if len(split) == 2 {
			tags[split[0]] = split[1]
		} else {
			return nil, fmt.Errorf("bad tag: %s", kv)
		}
	}
	return tags, nil
}

func makeUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func promptConfirm(question string) error {
	reader := bufio.NewReader(os.Stdin)
	result, _ := reader.ReadString('\n')
	result = strings.Replace(result, "\n", "", -1)
	if result != question {
		return fmt.Errorf("cancelled")
	}
	return nil
}
