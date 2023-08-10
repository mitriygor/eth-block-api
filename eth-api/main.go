package main

import (
	"eth-api/app"
	"fmt"
	"gorm.io/gorm"
)

// EthBlock model for GORM
type EthBlock struct {
	gorm.Model
	BlockNumber uint64 `gorm:"uniqueIndex;not null"`
	BlockHash   string `gorm:"uniqueIndex;not null"`
}

func main() {
	fmt.Println("************** Hi A!dfsdfsdfsdfPI **************")
	app := app.NewApp()

	// Start the server
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
