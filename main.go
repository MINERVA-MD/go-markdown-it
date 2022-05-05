package main

import (
	"fmt"
	"go-markdown-it/pkg"
	"regexp"
)

func main() {
	var str = "For more information, see Chapter 3.4.5.1"
	var re = regexp.MustCompile("(?i)see (chapter \\d+(\\.\\d)*)")
	var found = re.FindStringSubmatch(str)
	fmt.Println(found)

	fmt.Println(pkg.FromCodePoint(0x20))

}
