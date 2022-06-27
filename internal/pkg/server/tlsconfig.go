// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// server
package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/marmotedu/errors"
	"github.com/notone/pigpig/pkg/log"
	"github.com/notone/pigpig/pkg/util/fileutil"
)

type CertsCache struct {
	CacheMap map[string]*tls.Certificate
	RootCa   CertKey
	mu       sync.RWMutex
}

var certsCache CertsCache

func GetCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	log.Debugf("ready to load %s certificate", clientHello.ServerName)
	hostName := clientHello.ServerName
	rootCA := certsCache.RootCa
	if rootCA.KeyFile == "" || rootCA.CertFile == "" {
		return nil, errors.New("failed to load rootCA key or certificate")
	}
	certsCache.mu.Lock()
	defer certsCache.mu.Unlock()

	if certsCache.CacheMap == nil {
		certsCache.CacheMap = make(map[string]*tls.Certificate)
	}
	if certs, exist := certsCache.CacheMap[hostName]; exist {
		log.Debugf("loaded %s certificate successful", clientHello.ServerName)
		return certs, nil
	}
	certPath := strings.Replace(rootCA.CertFile, "PigPig", hostName, 1)
	privateKeyPath := strings.Replace(rootCA.KeyFile, "PigPig", hostName, 1)

	loadCert, err := LoadCertificate(clientHello.ServerName, certPath, privateKeyPath)
	if err == nil && loadCert != nil {
		certsCache.CacheMap[hostName] = loadCert
		return loadCert, nil
	}
	defaultCert, err := tls.LoadX509KeyPair(rootCA.CertFile, rootCA.KeyFile)

	generateCert, err := GenerateCertsForHostname(hostName, &defaultCert, certPath, privateKeyPath)
	if err == nil {
		certsCache.CacheMap[hostName] = generateCert
		return generateCert, err
	}

	log.Errorf("failed load certificate and key all ---> domain: %s", hostName)
	return nil, err
}

// GetTlsConfig initialize tlsconfig must provided a rootCA by the calls
func GetTlsConfig(rootCA CertKey) *tls.Config {
	certsCache.RootCa = rootCA
	return &tls.Config{
		GetCertificate:     GetCertificate,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
	}
}

func GenerateCertsForHostname(host string, rootCA *tls.Certificate, certPath, privateKeyPath string) (*tls.Certificate, error) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)

	// 定义：引用IETF的安全领域的公钥基础实施（PKIX）工作组的标准实例化内容
	subject := GetSubject(host)

	rootCert, _ := x509.ParseCertificate(rootCA.Certificate[0])
	// 设置 SSL证书的属性用途
	certificate509 := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(100 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	// 生成指定位数密匙
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	// 生成 SSL公匙
	derBytes, _ := x509.CreateCertificate(rand.Reader, &certificate509, rootCert, pk.Public(), rootCA.PrivateKey)
	certBuffer := &bytes.Buffer{}
	keyBuffer := &bytes.Buffer{}
	_ = pem.Encode(certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	_ = pem.Encode(keyBuffer, &pem.Block{Type: "RAS PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})

	go func() {
		certOut, _ := os.Create(certPath)
		defer certOut.Close()
		_, _ = certOut.Write(certBuffer.Bytes())
		// 生成 SSL私匙
		keyOut, _ := os.Create(privateKeyPath)
		defer keyOut.Close()
		_, _ = keyOut.Write(keyBuffer.Bytes())
		log.Debugf("certificate and key are finish")
	}()

	certs, err := tls.X509KeyPair(certBuffer.Bytes(), keyBuffer.Bytes())
	if err != nil {
		log.Errorf("found an error when generate certificate and key for %s ----> error: %s", host, err.Error())
		return nil, err
	}

	return &certs, nil
}

func LoadCertificate(host, certPath, privkeyPath string) (*tls.Certificate, error) {
	// TODO handle race condition (ask Matt)
	// the transaction is idempotent, however, so it shouldn't matter
	exist, _ := fileutil.FileExists(certPath)
	if !exist {
		log.Warn("failed to load certificate or key file")
		return nil, errors.New("failed to load certificate or key file")
	}
	exist, _ = fileutil.FileExists(privkeyPath)
	if !exist {
		log.Errorf("failed to load certificate or key file")
		return nil, errors.New("failed to load certificate or key file")
	}

	cert, err := tls.LoadX509KeyPair(certPath, privkeyPath)
	if err != nil {
		log.Debugf("load %s certificate successfully", host)
		return &cert, nil
	}
	return nil, err

}

func GetSubject(host string) pkix.Name {
	subject := pkix.Name{
		Country:            []string{"CN"},
		Organization:       []string{"PP"},
		OrganizationalUnit: []string{"PigPig"},

		Locality:   []string{"ChengDu"},
		Province:   []string{"SiChuan"},
		CommonName: host,
	}
	return subject
}
