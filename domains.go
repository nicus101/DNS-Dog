package main

import (
	"bufio"
	"fmt"
	"os"
)

func loadDomainList() []string {

	data, err := os.Open("domains.conf")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	var DomainList []string
	for scanner.Scan() {
		line := scanner.Text()
		DomainList = append(DomainList, line)
	}

	return DomainList
}
