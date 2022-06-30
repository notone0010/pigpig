// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	rd "math/rand"
	"os"
	"time"
)

// CertInformation certificate information.
type CertInformation struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Province           []string
	Locality           []string
	CommonName         string
	CrtName, KeyName   string
	IsCA               bool
	Names              []pkix.AttributeTypeAndValue
}

func newCertificate(info CertInformation) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			CommonName:         info.CommonName,
			Locality:           info.Locality,
			ExtraNames:         info.Names,
		},
		NotBefore:             time.Now(),                                                                 // 证书的开始时间
		NotAfter:              time.Now().AddDate(20, 0, 0),                                               // 证书的结束时间
		BasicConstraintsValid: true,                                                                       // 基本的有效性约束
		IsCA:                  info.IsCA,                                                                  // 是否是根证书
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, // 证书用途
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		EmailAddresses:        info.EmailAddress,
	}
}

// CreateCRT create certificate and key pair.
func CreateCRT(RootCa *x509.Certificate, RootKey *rsa.PrivateKey, info CertInformation) error {
	Crt := newCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var buf []byte
	if RootCa == nil || RootKey == nil {
		// 创建自签名证书
		buf, err = x509.CreateCertificate(rand.Reader, Crt, Crt, &Key.PublicKey, Key)
	} else {
		// 使用根证书签名
		buf, err = x509.CreateCertificate(rand.Reader, Crt, RootCa, &Key.PublicKey, RootKey)
	}
	if err != nil {
		return err
	}

	// return nil
	err = write(info.CrtName, "CERTIFICATE", buf)
	if err != nil {
		return err
	}

	buf = x509.MarshalPKCS1PrivateKey(Key)

	return write(info.KeyName, "PRIVATE KEY", buf)
}

func write(filename, Type string, p []byte) error {
	File, err := os.Create(filename)
	// defer File.Close()
	if err != nil {
		return err
	}
	b := &pem.Block{Bytes: p, Type: Type}

	defer File.Close()

	return pem.Encode(File, b)
}

// Parse parse ssl certificate and key.
func Parse(crtPath, keyPath string) (rootcertificate *x509.Certificate, rootPrivateKey *rsa.PrivateKey, err error) {
	rootcertificate, err = ParseCrt(crtPath)
	if err != nil {
		return
	}
	rootPrivateKey, err = ParseKey(keyPath)

	return
}

// ParseCrt parse ssl certificate.
func ParseCrt(path string) (*x509.Certificate, error) {
	var p *pem.Block

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, _ = pem.Decode(buf)

	return x509.ParseCertificate(p.Bytes)
}

// ParseKey parse ssl key.
func ParseKey(path string) (*rsa.PrivateKey, error) {
	var p *pem.Block
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, _ = pem.Decode(buf)

	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func main() {
	baseinfo := CertInformation{
		Country:            []string{"CN"},
		Organization:       []string{"PP"},
		IsCA:               true,
		OrganizationalUnit: []string{"PigPig"},
		EmailAddress:       []string{"aiphalv0@163.com"},

		Locality:   []string{"ChengDu"},
		Province:   []string{"SiChuan"},
		CommonName: "PigPig",
		CrtName:    "./configs/cert/pigpig.crt",
		KeyName:    "./configs/cert/pigpig.key",
	}

	err := CreateCRT(nil, nil, baseinfo)
	if err != nil {
		log.Printf("Create crt error,Error info: %s", err.Error())
		return
	}
	log.Printf("Create crt success file -> certificate: %s, key: %s", baseinfo.CrtName, baseinfo.KeyName)
}
