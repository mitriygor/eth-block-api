package eth_block_helper

import (
	"fmt"
	"regexp"
	"strconv"
)

func IsInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func IsHex(s string) bool {
	pattern := `^0x[0-9a-fA-F]+$`
	match, _ := regexp.MatchString(pattern, s)
	return match && len(s) < 32
}

func StringToInt(str string) int {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return -1
	}

	return intValue
}

func IntToHex(num int) string {
	n := int64(num)
	return fmt.Sprintf("0x%x", n)
}
