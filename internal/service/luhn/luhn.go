package luhn

import (
	"strconv"
	"strings"
)

func CheckLuhn(number string) bool {
	number = strings.TrimSpace(number)
	sum := 0
	n := len(number)
	parity := n % 2
	for i := 0; i < n; i++ {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
