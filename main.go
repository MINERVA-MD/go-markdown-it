package main

import (
	"fmt"
	"go-markdown-it/pkg"
	"go-markdown-it/pkg/utils"
	"regexp"
)

func main() {
	var str = "For more information, see Chapter 3.4.5.1"
	var re = regexp.MustCompile("(?i)see (chapter \\d+(\\.\\d)*)")
	var found = re.FindStringSubmatch(str)
	fmt.Println(found)

	var mdurl = pkg.MdUrl{}
	utils.PrettyPrint(mdurl.Parse("http://example.com?find=\\*", true))
}
