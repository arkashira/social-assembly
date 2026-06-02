package main

import (
	"fmt"
	"log"

	"github.com/axentx/axentx-social-assembly/pkg/moderation"
)

func main() {
	if err := moderation.LoadRules(); err != nil {
		log.Fatal(err)
	}

	rule, ok := moderation.GetRule("test_rule")
	if !ok {
		log.Fatal("rule not found")
	}

	fmt.Println(rule)
}