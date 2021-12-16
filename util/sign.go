package util

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"hash"
	"io/ioutil"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/crypto/pkcs12"
)

const (
	SignType_MD5    = "MD5"
	SignType_SHA1   = "SHA1"
	SignType_SHA256 = "SHA256"
)

func Md5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// 公钥加密  注意：用了base64 编码
func PublicEncrypt(data string, rootPEM string) (string, error) {
	block, _ := pem.Decode([]byte(rootPEM))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	pub := cert.PublicKey.(*rsa.PublicKey)
	v, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(v), err
}

// 私钥解析 注意：用了base64 编码
func PrivateEncrypt(encryptedKey string, ca string, certData string, password string) (k []byte, err error) {
	key, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return nil, err
	}
	var certD []byte
	if ca != "" {
		certD, err = ioutil.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("unable to find cert path=%s, error=%v", err)
		}
	}
	if certData != "" {
		certD, err = base64.StdEncoding.DecodeString(certData)
		if err != nil {
			return nil, fmt.Errorf("certData 商家秘钥文件转码错误", err)
		}
	}
	privateKey, _, err := pkcs12.Decode(certD, password)
	if err != nil {
		return nil, err
	}
	pri := privateKey.(*rsa.PrivateKey)
	v, err := rsa.DecryptPKCS1v15(rand.Reader, pri, key)
	if err != nil {
		return nil, err
	}
	return v, err
}

// Sha256Base64 进行SHA-256哈希值并进行Base64编码
func Sha256Base64(body []byte) (sum string) {
	h := sha256.New()
	h.Write(body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// VerifySign 验证支付
func VerifySign(params map[string]interface{}, sign string, rootPEM string, signType string) (ok bool, err error) {
	encodeSignParams := EncodeSignParams(params)
	var (
		h     hash.Hash
		hashs crypto.Hash
	)
	signBytes, _ := base64.StdEncoding.DecodeString(sign)
	block, _ := pem.Decode([]byte(rootPEM))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	publicKey := cert.PublicKey.(*rsa.PublicKey)
	switch signType {
	case "RSA":
		hashs = crypto.SHA1
	case "RSA2":
		hashs = crypto.SHA256
	default:
		hashs = crypto.SHA256
	}
	h = hashs.New()
	h.Write([]byte(encodeSignParams))
	err = rsa.VerifyPKCS1v15(publicKey, hashs, h.Sum(nil), signBytes)
	if err != nil {
		return ok, err
	}
	return true, err
}

// EncodeSignParams 编码符号参数
func EncodeSignParams(params map[string]interface{}) string {
	var buf strings.Builder
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "signatureString" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := params[k]
		if v == "" {
			continue
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(InterfaceToString(v))
		buf.WriteByte('&')
	}
	return buf.String()[:buf.Len()-1]
}

// Sign 开发平台签名支付签名.
//  params: 待签名的参数集合
//  privateKey: 密钥
func Sign(params map[string]interface{}, ca string, certData string, password string, signType string) (sign string, err error) {
	encodeSignParams := EncodeSignParams(params)
	var (
		// block          *pem.Block
		h              hash.Hash
		key            *rsa.PrivateKey
		hashs          crypto.Hash
		encryptedBytes []byte
	)
	var certD []byte
	if ca != "" {
		certD, err = ioutil.ReadFile(ca)
		if err != nil {
			return "", fmt.Errorf("unable to find cert path=%s, error=%v", err)
		}
	}
	if certData != "" {
		certD, err = base64.StdEncoding.DecodeString(certData)
		if err != nil {
			return "", fmt.Errorf("certData 商家秘钥文件转码错误", err)
		}
	}
	privateKey, _, err := pkcs12.Decode(certD, password)
	if err != nil {
		return "", err
	}
	key = privateKey.(*rsa.PrivateKey)
	switch signType {
	case "RSA":
		h = sha1.New()
		hashs = crypto.SHA1
	case "RSA2":
		h = sha256.New()
		hashs = crypto.SHA256
	default:
		h = sha256.New()
		hashs = crypto.SHA256
	}
	if _, err = h.Write([]byte(encodeSignParams)); err != nil {
		return
	}
	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, key, hashs, h.Sum(nil)); err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
}

// Sign 开发平台签名支付签名.
//  data: 待签名数据字符串
//  secret: 密钥
func HmacSha1(data string, secret string) (sign string) {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// ParseNotifyResult 解析异步通知
func InterfaceToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case float32:
		return strconv.FormatFloat(v.(float64), 'E', -1, 32)
	case float64:
		return strconv.FormatFloat(v.(float64), 'E', -1, 64)
	}
	return ""
}

// FormatPrivateKey 格式化 普通应用秘钥
func FormatPrivateKey(privateKey string) (pKey string) {
	var buffer strings.Builder
	buffer.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")
	rawLen := 64
	keyLen := len(privateKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen
	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(privateKey[start:])
		} else {
			buffer.WriteString(privateKey[start:end])
		}
		buffer.WriteByte('\n')
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END RSA PRIVATE KEY-----\n")
	pKey = buffer.String()
	return
}

// FormatParam 格式化请求参数
func FormatParam(params map[string]interface{}, appSecret string) (s string) {
	for key, value := range params {
		s = s + key + InterfaceToString(value)
	}
	return s + appSecret
}

// FormatURLParam 格式化请求URL参数
func FormatURLParam(params map[string]interface{}) (urlParam string) {
	v := url.Values{}
	for key, value := range params {
		v.Add(key, InterfaceToString(value))
	}
	return v.Encode()
}

// getSignData 获取数据字符串
func GetSignData(str string) (signData string) {
	indexStart := strings.Index(str, `":`)
	indexEnd := strings.Index(str, `,"sign"`)
	if indexEnd == -1 {
		indexEnd = strings.Index(str, `}}`) + 1
	}
	signData = str[indexStart+2 : indexEnd]
	return
}
