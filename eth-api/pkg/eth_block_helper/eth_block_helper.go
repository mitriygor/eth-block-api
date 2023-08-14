package eth_block_helper

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func IsInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func IsHex(s string) bool {
	_, err := strconv.ParseUint(s, 16, 64)
	return err == nil
}

func StringToInt(str string) int {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("StringToInt::error: %v\n", err)
		return -1
	}

	return intValue
}

func HexToInt(hexStr string) int {

	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}

	intValue, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		log.Printf("SetCurrentBlockNumber::error: %v\n", err)
		return -1
	}

	return int(intValue)
}

func IntToHex(num int) string {
	n := int64(num)
	return fmt.Sprintf("0x%x", n)
}
