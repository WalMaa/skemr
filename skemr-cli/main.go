package main

import (
	"fmt"
	"github.com/walmaa/skemr-cli/cmd"
)

const skemrAscii = `   _____ _                       
  / ____| |                      
 | (___ | | _____ _ __ ___  _ __ 
  \___ \| |/ / _ \ '_ ` + "`" + ` _ \| '__|
  ____) |   <  __/ | | | | | |   
 |_____/|_|\_\___|_| |_| |_|_|   
                                 
                                 `

func main() {

	fmt.Println(skemrAscii)

	cmd.Execute()
}
