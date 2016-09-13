package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/timehop/go-mixpanel"
)

func usage() {
	fmt.Printf("usage: %v <token> <distinct_id> '<operation>' prop1=val1 [prop2=val2 ...]\n", os.Args[0])
}

func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	token := os.Args[1]
	distinctID := os.Args[2]

	op := &mixpanel.Operation{Name: os.Args[3], Values: map[string]interface{}{}}
	if op.Name == "" {
		op.Name = "$set"
	}

	if len(os.Args) >= 5 {
		for i := 4; i < len(os.Args); i++ {
			parts := strings.Split(os.Args[i], "=")
			if len(parts) == 2 {
				op.Values[parts[0]] = parts[1]
			}
		}
	}

	mp := mixpanel.NewMixpanel(token)
	err := mp.Engage(distinctID, map[string]interface{}{}, op)
	if err != nil {
		fmt.Println("Error occurred:", err)
	}
}
