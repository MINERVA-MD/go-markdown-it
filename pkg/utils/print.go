package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
}
