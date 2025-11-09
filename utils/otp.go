package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid OTP length")
	}

	digits := "0123456789"
	otp := make([]byte, length)

	_, err := rand.Read(otp)
	if err != nil {
		return "", fmt.Errorf("failed to generate otp: %w", err)
	}

	for i := 0; i < length; i++ {
		otp[i] = digits[int(otp[i])%len(digits)]
	}

	return string(otp), nil
}
