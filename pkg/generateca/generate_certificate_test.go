// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// server
package generateca

import (
	"testing"
)

func Test_crt(t *testing.T) {
	baseinfo := CertInformation{
		Country:            []string{"CN"},
		Organization:       []string{"PP"},
		IsCA:               true,
		OrganizationalUnit: []string{"PigPig"},
		EmailAddress:       []string{"aiphalv0@163.com"},

		Locality:   []string{"ChengDu"},
		Province:   []string{"SiChuan"},
		CommonName: "PigPig",
		CrtName:    "/Users/lvsongke/workspace/go/golang/src/github.com/notone0010/pigpig/configs/cert/PigPig.crt",
		KeyName:    "/Users/lvsongke/workspace/go/golang/src/github.com/notone0010/pigpig/configs/cert/PigPig.key",
	}

	err := CreateCRT(nil, nil, baseinfo)
	if err != nil {
		t.Log("Create crt error,Error info:", err)
		return
	}
	// crtinfo := baseinfo
	// crtinfo.IsCA = false
	// crtinfo.CrtName = "test_server.crt"
	// crtinfo.KeyName = "test_server.key"
	// crtinfo.Names = []pkix.AttributeTypeAndValue{{asn1.ObjectIdentifier{2, 1, 3}, "MAC_ADDR"}} // 添加扩展字段用来做自定义使用
	//
	// crt, pri, err := Parse(baseinfo.CrtName, baseinfo.KeyName)
	// if err != nil {
	// 	t.Log("Parse crt error,Error info:", err)
	// 	return
	// }
	// err = CreateCRT(crt, pri, crtinfo)
	// if err != nil {
	// 	t.Log("Create crt error,Error info:", err)
	// }
	// os.Remove(baseinfo.CrtName)
	// os.Remove(baseinfo.KeyName)
	// os.Remove(crtinfo.CrtName)
	// os.Remove(crtinfo.KeyName)
}
