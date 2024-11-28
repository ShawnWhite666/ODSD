package utils

import "testing"

func TestIp4toInt(t *testing.T) {
	str := "192.168.32.33"
	t.Log(Ip4toInt(str))
}
