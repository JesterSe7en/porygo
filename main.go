/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com
*/

package main

import (
	"github.com/JesterSe7en/scrapego/cmd"
	"github.com/JesterSe7en/scrapego/internal/logger"
)

func main() {
	logger.InitLogger()
	cmd.Execute()
}
