package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
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
	logLevel := slog.LevelDebug
	// Logger colors
	w := os.Stderr
	// Set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.DateTime,
		}),
	))

	fmt.Println(skemrAscii)

	cmd.Execute()
}
