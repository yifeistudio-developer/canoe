package util

import "testing"

func TestGetLocalIpAddress(t *testing.T) {

	ip, err := GetLocalIpAddress()
	println(ip, err)
}
