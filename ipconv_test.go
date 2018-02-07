package ipgeo

import (
	"bytes"
	"fmt"
	"testing"
)

type TestData struct {
	ip      string
	ipmask  string
	u32     uint32
	u32mask uint32
	mask    uint32
	match   bool
}

var testData = []TestData{
	{"223.255.255.9", "223.255.255.0/24", 3758096137, 3758096128, 4294967040, true},
	{"223.255.215.9", "223.255.255.0/24", 3758085897, 3758096128, 4294967040, false},
	{"46.183.182.9", "46.183.182.0/24", 0, 0, 0, true},
	{"73.183.190.129", "73.183.190.128/25", 0, 0, 0, true},
}

func TestIPNet(t *testing.T) {
	for k, v := range testData {
		b := IPNet(v.ipmask)
		if len(b) != 8 {
			fmt.Println("len is :", len(b))
			t.Errorf("%d\n", k)
		}
		//fmt.Println("ipnet bytes:", b)
	}
}
func TestIP2IPNet(t *testing.T) {

	for k, v := range testData {
		match := false
		var b []byte
		for i := 0; i < 32; i++ {
			b = IP2IPNet(v.ip, uint32(i))
			fmt.Println(k, "ip2ipnet bytes:", b)
			if v.match == bytes.Equal(b, IPNet(v.ipmask)) {
				fmt.Println("b:", k, b, IPNet(v.ipmask))
				match = true
				break
			}
		}
		if match == false {
			fmt.Printf("b:%b\n", b)
			t.Errorf("b:%d", k)
		}
	}

}

/*
func TestMatch(t *testing.T) {
	for k, v := range testData {
		//fmt.Println(v.ip)
		//fmt.Println(v.ipmask)
		//fmt.Printf("ip:%b\n",IpGetUint32(v.ip))
		//fmt.Printf("ip:%d\n",IpGetUint32(v.ip))
		//fmt.Printf("ipmask:%b\n",IpMaskGetUint32(v.ipmask))
		//fmt.Printf("ipmask:%d\n",IpMaskGetUint32(v.ipmask))
		if IpMaskGetUint32(v.ipmask) != v.u32mask {
			t.Errorf("%d\n", k)
		}
		if IpGetUint32(v.ip) != v.u32 {
			t.Errorf("%d\n", k)
		}
		if Match(IpGetUint32(v.ip), IpMaskGetUint32(v.ipmask), v.mask) != v.match {
			t.Errorf("%d\n", k)
		}
	}

	return
}

func TestSearch(t *testing.T) {

	return
	for k, v := range testData {
		if v.match != bytes.Equal(Search(v.ip, []byte(v.ipmask)), []byte(v.ipmask)) {
			t.Errorf("%d\n", k)
		}
	}

	return
}
func TestContains(t *testing.T) {

	for k, v := range testData {
		if v.match != Contains(v.ipmask, v.ip) {
			t.Errorf("%d\n", k)
		}
	}

}*/
