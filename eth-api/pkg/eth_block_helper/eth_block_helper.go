package eth_block_helper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func IsInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func IsHex(s string) bool {
	pattern := `^0x[0-9a-fA-F]+$`
	match, _ := regexp.MatchString(pattern, s)
	log.Printf("eth-api::IsHex::match: %v\n", match)
	return match && len(s) < 32
}

func StringToInt(str string) int {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("eth-api::ERROR::StringToInt::error: %v\n", err)
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
		log.Printf("eth-api::ERROR::SetCurrentBlockNumber::error: %v\n", err)
		return -1
	}

	return int(intValue)
}

func IntToHex(num int) string {
	n := int64(num)
	return fmt.Sprintf("0x%x", n)
}
