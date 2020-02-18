package main

import (
	"fmt"
	"phpengine"
)

func main()  {
	pe := phpengine.NewPhpEngine("/home/wangchao/php74/bin/php")
	s := pe.NewSession("/home/wangchao/go/src/amdop/localpkg/phpengine/test/test.php")
	fmt.Println(s.Execute(), s.GetLastOutput())
}
