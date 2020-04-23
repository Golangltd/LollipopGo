/*
 *  Copyright (c) 2017, https://github.com/nebulaim
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
)

func TestNewRSACryptor(t *testing.T) {

	rsa := NewRSACryptor()

	testData := []byte("rsa 2048 key!!!!")
	fmt.Println(string(testData))

	encData := rsa.Encrypt(testData)
	fmt.Println(hex.EncodeToString(encData))

	decData := rsa.Decrypt(encData)
	fmt.Println("len = ", len(decData), ", data: ", string(decData))
}

//import (
//"fmt"
//"github.com/nebulaim/telegramd/mtproto"
//"crypto/x509"
//"crypto/rsa"
//"encoding/pem"
//"errors"
//"crypto/rand"
//"math/big"
//)
//
//type Message interface {
//	Decode([]byte) error
//}
//
//type A1 struct {
//}
//
//func (m *A1) Decode(b []byte) error  {
//	fmt.Println("A1.Decode()")
//	return nil
//}
//
//type A2 struct {
//}
//
//func (m *A2) Decode(b []byte) error  {
//	fmt.Println("A2.Decode()")
//	return nil
//}
//
//type NewMessageFunc func() Message
//
//var registers = map[int32]NewMessageFunc{
//	1 : func() (Message) { return new(A1) },
//	2 : func() (Message) { return &A2{} },
//}
//
//func NewMessage(i int32) *Message {
//	// return nil
//	m := registers[i]()
//	m.Decode([]byte{})
//	return &m
//}
//
//type newTLObjectFunc func() interface{}
//
//var registers2 = map[mtproto.TLConstructor]newTLObjectFunc {
//	mtproto.TLConstructor_CRC32_p_q_inner_data : func() (interface{}) {return new(mtproto.TLPQInnerData) },
//}
//
//func NewTLObjectByClassID(classId mtproto.TLConstructor) interface{} {
//	m, ok := registers2[classId]
//	if !ok {
//		return nil
//	}
//	return m()
//}
//
//var serverPublicKeys = []byte(`
//-----BEGIN RSA PUBLIC KEY-----
//MIIBCgKCAQEAtUXgOV7DZ1d9M08gYVOMU/fenTbbLU12b1LoL9sYfRycEpF4aqA9
//RW0rPfzY6yZkfTlQdoFaGxVpUiNMv5V3xY+aVoFQbT7rlsevE87tHK90yG1OYysl
//V7IJiCy/tLu/2DVhbZqg4fgPjs4XYrt7CABmsy/OtHJy6A9I1qPQ40MlSB21lmAQ
//I9gKHBc2BGZCvQ/NDj1elun9Qitf3ziH8g/Xsxv18CO8hAev56FUMIFzMtGOmhpJ
//DAkQ+zg22yLlxKuxjkWSEkYYzHgzrCiDfqcfSkG34veRdD9CGnLsIPvHtTFV/+b0
//5xTUyQxFzZ3kOl41KsTY9OsebYxYThHbTQIDAQAB
//-----END RSA PUBLIC KEY-----
//`)
//
//// 加密
//func RsaEncrypt(origData []byte) ([]byte, error) {
//	block, _ := pem.Decode(serverPublicKeys)
//	if block == nil {
//		return nil, errors.New("public key error")
//	}
//	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//	pub := pubInterface.(*rsa.PublicKey)
//	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
//}
//
///*
//// 解密
//func RsaDecrypt(ciphertext []byte) ([]byte, error) {
//	block, _ := pem.Decode(privateKey)
//	if block == nil {
//		return nil, errors.New("private key error!")
//	}
//	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
//}
//*/
//
//// var a = rand.Reader
//
//
//var pcks8PemPublicKey = []byte(`
//-----BEGIN PUBLIC KEY-----
//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvKLEOWTzt9Hn3/9Kdp/R
//dHcEhzmd8xXeLSpHIIzaXTLJDw8BhJy1jR/iqeG8Je5yrtVabqMSkA6ltIpgylH/
////FojMsX1BHu4EPYOXQgB0qOi6kr08iXZIH9/iOPQOWDsL+Lt8gDG0xBy+sPe/2Z
//HdzKMjX6O9B4sOsxjFrk5qDoWDrioJorAJ7eFAfPpOBf2w73ohXudSrJE0lbQ8pC
//WNpMY8cB9i8r+WBitcvouLDAvmtnTX7akhoDzmKgpJBYliAY4qA73v7u5UIepE8Q
//gV0jCOhxJCPubP8dg+/PlLLVKyxU5CdiQtZj2EMy4s9xlNKzX8XezE0MHEa6bQpn
//FwIDAQAB
//-----END PUBLIC KEY-----
//`)
//
//var pcks1PemPublicKey = []byte(`
//-----BEGIN RSA PUBLIC KEY-----
//MIIBCgKCAQEAvKLEOWTzt9Hn3/9Kdp/RdHcEhzmd8xXeLSpHIIzaXTLJDw8BhJy1
//jR/iqeG8Je5yrtVabqMSkA6ltIpgylH///FojMsX1BHu4EPYOXQgB0qOi6kr08iX
//ZIH9/iOPQOWDsL+Lt8gDG0xBy+sPe/2ZHdzKMjX6O9B4sOsxjFrk5qDoWDrioJor
//AJ7eFAfPpOBf2w73ohXudSrJE0lbQ8pCWNpMY8cB9i8r+WBitcvouLDAvmtnTX7a
//khoDzmKgpJBYliAY4qA73v7u5UIepE8QgV0jCOhxJCPubP8dg+/PlLLVKyxU5Cdi
//QtZj2EMy4s9xlNKzX8XezE0MHEa6bQpnFwIDAQAB
//-----END RSA PUBLIC KEY-----
//`)
//
//var pcks1PemPrivateKey = []byte(`
//-----BEGIN RSA PRIVATE KEY-----
//MIIEpAIBAAKCAQEAvKLEOWTzt9Hn3/9Kdp/RdHcEhzmd8xXeLSpHIIzaXTLJDw8B
//hJy1jR/iqeG8Je5yrtVabqMSkA6ltIpgylH///FojMsX1BHu4EPYOXQgB0qOi6kr
//08iXZIH9/iOPQOWDsL+Lt8gDG0xBy+sPe/2ZHdzKMjX6O9B4sOsxjFrk5qDoWDri
//oJorAJ7eFAfPpOBf2w73ohXudSrJE0lbQ8pCWNpMY8cB9i8r+WBitcvouLDAvmtn
//TX7akhoDzmKgpJBYliAY4qA73v7u5UIepE8QgV0jCOhxJCPubP8dg+/PlLLVKyxU
//5CdiQtZj2EMy4s9xlNKzX8XezE0MHEa6bQpnFwIDAQABAoIBACd+SGjfyursZoiO
//MW/ejAK/PFJ3bKtNI8P++v9Enh8vF8swUBgMmzIdv93jZfnnD1mtT46kU6mXd3fy
//FMunGVrjlwkLKET9MC8B5U46Es6T/H4fAA8KCzA+ywefOEnVA5pIsB7dIFFhyNDB
//uO8zrBAFfsu+Y1KMlggsZaZGDXB/WVyUJDbEOMZstVx4uNhpcEgKYp28YQMP/yvv
//dp4UgnTxXXXpDghzO5iqi5tUWY0p1lH2ii2OZBxEdqdDl7TirorhUDYIivyoe3B5
//H30RNBRok/6w7W0WPyY2lSIcjd3cLPte6vx0QfBXVo2A6N9LTKAtAw3iWBp0x9NZ
//N5p8OeECgYEA8QywXlM8nH5M7Sg2sMUYBOHA22O26ZPio7rJzcb8dlkV5gVHm+Kl
//aDP61Uy8KoYABQ5kFdem/IQAUPepLxmJmiqfbwOIjfajOD3uVAQunFnDCHBWm4Uk
//onbpdA5NlT/OUoSjIBemiBR/4CpDK1cEby/sg+EvQaGxqtedEe4xFmcCgYEAyFXe
//MyAAOLpzmnCs9NYTTvMPofW8y+kLDodfbskl7M8q6l20VMo/E+g1gQ+65Aah901Z
///LKGi6HpzmHi5q9O2OJtqyI6FVwjXa07M5ueDbHcVKJw4hC9W0oHpMg8hqumPAWF
//+MoN/Toy77p5LzoR30WUdhPvOAJPEL1p2a6r29ECgYEAiXfCEVkI5PqGZm2bmv4b
//75TLhpJ8WwMSqms48Vi828V8Xpy+NOFxkVargv9rBBk9Y6TMYUSGH9Yr1AEZhBnd
//RoVuPUJXmxaACPAQvetQpavvNR3T1od82AZWpvANQMONp7Oqz/+M4mhGcRHJEqti
//hQJgsOk4KQbMqvChy/r6FZsCgYEAwyaqgkD9FkXC0UJLqWFUg8bQhqPcGwLUC34h
//n8kAUbPpiU5omWQ+mATPAf8xvmkbo81NCJVb7W93U90U7ET/2NSRonCABkiwBtP2
//ZKqGB68oA6YNspo960ytL38DPui80aFLxXQGtpPYBKEw5al6uXWNTozSrkvJe3QY
//Rb4amdECgYBpGk7zPcK1TbJ++W5fkiory4qOdf0L1Zf0NbML4fY6dIww+dwMVUpq
//FbsgCLqimqOFaaECU+LQEFUHHM7zrk7NBf7GzBvQ+qJx8zhJ66sFVox+IirBUyR9
//Vh0+z5tIbFbKmYkO06NbeMlq87JexSlocPZtA3HMhEga5/0fHNHsNw==
//-----END RSA PRIVATE KEY-----
//`)
//
//var pcks8PemPrivateKey = []byte(`
//-----BEGIN PRIVATE KEY-----
//MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC8osQ5ZPO30eff
///0p2n9F0dwSHOZ3zFd4tKkcgjNpdMskPDwGEnLWNH+Kp4bwl7nKu1VpuoxKQDqW0
//imDKUf//8WiMyxfUEe7gQ9g5dCAHSo6LqSvTyJdkgf3+I49A5YOwv4u3yAMbTEHL
//6w97/Zkd3MoyNfo70Hiw6zGMWuTmoOhYOuKgmisAnt4UB8+k4F/bDveiFe51KskT
//SVtDykJY2kxjxwH2Lyv5YGK1y+i4sMC+a2dNftqSGgPOYqCkkFiWIBjioDve/u7l
//Qh6kTxCBXSMI6HEkI+5s/x2D78+UstUrLFTkJ2JC1mPYQzLiz3GU0rNfxd7MTQwc
//RrptCmcXAgMBAAECggEAJ35IaN/K6uxmiI4xb96MAr88Undsq00jw/76/0SeHy8X
//yzBQGAybMh2/3eNl+ecPWa1PjqRTqZd3d/IUy6cZWuOXCQsoRP0wLwHlTjoSzpP8
//fh8ADwoLMD7LB584SdUDmkiwHt0gUWHI0MG47zOsEAV+y75jUoyWCCxlpkYNcH9Z
//XJQkNsQ4xmy1XHi42GlwSApinbxhAw//K+92nhSCdPFddekOCHM7mKqLm1RZjSnW
//UfaKLY5kHER2p0OXtOKuiuFQNgiK/Kh7cHkffRE0FGiT/rDtbRY/JjaVIhyN3dws
//+17q/HRB8FdWjYDo30tMoC0DDeJYGnTH01k3mnw54QKBgQDxDLBeUzycfkztKDaw
//xRgE4cDbY7bpk+KjusnNxvx2WRXmBUeb4qVoM/rVTLwqhgAFDmQV16b8hABQ96kv
//GYmaKp9vA4iN9qM4Pe5UBC6cWcMIcFabhSSidul0Dk2VP85ShKMgF6aIFH/gKkMr
//VwRvL+yD4S9BobGq150R7jEWZwKBgQDIVd4zIAA4unOacKz01hNO8w+h9bzL6QsO
//h19uySXszyrqXbRUyj8T6DWBD7rkBqH3TVn8soaLoenOYeLmr07Y4m2rIjoVXCNd
//rTszm54NsdxUonDiEL1bSgekyDyGq6Y8BYX4yg39OjLvunkvOhHfRZR2E+84Ak8Q
//vWnZrqvb0QKBgQCJd8IRWQjk+oZmbZua/hvvlMuGknxbAxKqazjxWLzbxXxenL40
//4XGRVquC/2sEGT1jpMxhRIYf1ivUARmEGd1GhW49QlebFoAI8BC961Clq+81HdPW
//h3zYBlam8A1Aw42ns6rP/4ziaEZxEckSq2KFAmCw6TgpBsyq8KHL+voVmwKBgQDD
//JqqCQP0WRcLRQkupYVSDxtCGo9wbAtQLfiGfyQBRs+mJTmiZZD6YBM8B/zG+aRuj
//zU0IlVvtb3dT3RTsRP/Y1JGicIAGSLAG0/ZkqoYHrygDpg2ymj3rTK0vfwM+6LzR
//oUvFdAa2k9gEoTDlqXq5dY1OjNKuS8l7dBhFvhqZ0QKBgGkaTvM9wrVNsn75bl+S
//KivLio51/QvVl/Q1swvh9jp0jDD53AxVSmoVuyAIuqKao4VpoQJT4tAQVQcczvOu
//Ts0F/sbMG9D6onHzOEnrqwVWjH4iKsFTJH1WHT7Pm0hsVsqZiQ7To1t4yWrzsl7F
//KWhw9m0DccyESBrn/R8c0ew3
//-----END PRIVATE KEY-----
//`)
//
///*
//  GO语言没找到支持PKCS8格式的操作
//  http://blog.qiujinwu.com/2017/07/14/rsa/
//  https://medium.com/@oyrxx/rsa%E7%A7%98%E9%92%A5%E4%BB%8B%E7%BB%8D%E5%8F%8Aopenssl%E7%94%9F%E6%88%90%E5%91%BD%E4%BB%A4-d3fcc689513f
//  openssl genrsa -out server.key 2048
//  openssl rsa -in server.key -pubout > public_pcks8.pub
//  openssl rsa -in server.key -outform PEM -RSAPublicKey_out -out public_pcks1.key
// */
//
//// 加密
//func getPublicKey(pemPublicKey []byte) (pub *rsa.PublicKey, err error) {
//	block, _ := pem.Decode(pemPublicKey)
//	if block == nil {
//		fmt.Println("public key error")
//		return nil, errors.New("public key error")
//	}
//
//	// ParsePKCS1PrivateKey
//	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		fmt.Println("parse pkix key error: ", err)
//		// return
//		// return nil, err
//	}
//	pub = pubInterface.(*rsa.PublicKey)
//	fmt.Println("Modulus : ", pub.N.String())
//	fmt.Println(">>> ", pub.N)
//	fmt.Printf("Modulus(Hex) : %X\n", pub.N)
//	fmt.Println("Public Exponent : ", pub.E)
//
//	// return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
//	return
//}
//
//
//func getPrivateKey(pemPrivateKey []byte) (prv *rsa.PrivateKey, err error) {
//	// testPrivateKey
//	block, _ := pem.Decode([]byte(pemPrivateKey))
//	if block == nil {
//		fmt.Println("public key error")
//		return nil, errors.New("public key error")
//	}
//
//	prv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
//
//	if err != nil {
//		panic("Failed to parse private key: " + err.Error())
//	}
//	// pri := pubInterface.(*rsa.PrivateKey)
//	fmt.Println("Modulus : ", prv.D.String())
//	fmt.Println(">>> ", prv.Primes)
//	// fmt.Printf("Modulus(Hex) : %X\n", prv.N)
//	// fmt.Println("Public Exponent : ", prv.E)
//	// return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
//	return
//}
//
//
//func doRSAencrypt(em []byte) []byte {
//	publicKey, _ := getPublicKey(pcks8PemPublicKey)
//
//	// z := make([]byte, 255)
//	// copy(z, em)
//
//	c := new(big.Int)
//	c.Exp(new(big.Int).SetBytes(em), big.NewInt(int64(publicKey.E)), publicKey.N)
//
//	res := make([]byte, 256)
//	copy(res, c.Bytes())
//
//	return res
//}
//
//func doRSAdecrypt(em []byte) []byte {
//	privateKey, _ := getPrivateKey(pcks1PemPrivateKey)
//
//	//z := make([]byte, 255)
//	//copy(z, em)
//
//	c := new(big.Int)
//	c.Exp(new(big.Int).SetBytes(em), privateKey.D, privateKey.N)
//
//	// res := make([]byte, 256)
//	// copy(res, c.Bytes())
//
//	return c.Bytes()
//}
//
//func main()  {
//	// processPrivateKey()
//	// fmt.Println("-----------------------------------------------------------------------")
//	// processPublicKey()
//
//	/*
//		testData := []byte("rsa 2048 key!!!!")
//		fmt.Println(string(testData))
//
//		encData := rsa.Encrypt(testData)
//		fmt.Println(string(encData))
//
//		decData := rsa.Decrypt(encData)
//		fmt.Println("len = ", len(decData), ", data: ", string(decData))
//	 */
//
//
//	//rsa := mtproto.NewRSACryptor()
//	//
//	//var PQ = string([]byte{0x17, 0xED, 0x48, 0x94, 0x1A, 0x08, 0xF9, 0x81})
//	//var P = string([]byte{0x49, 0x4C, 0x55, 0x3B})
//	//var Q = string([]byte{0x53, 0x91, 0x10, 0x73})
//	//
//	//pqInnerData := &mtproto.TLPQInnerData{}
//	//pqInnerData.Nonce = mtproto.GenerateNonce(16)
//	//pqInnerData.ServerNonce = mtproto.GenerateNonce(16)
//	//pqInnerData.NewNonce = mtproto.GenerateNonce(32)
//	//pqInnerData.P = P
//	//pqInnerData.Q = Q
//	//pqInnerData.Pq = PQ
//	//fmt.Println(pqInnerData)
//	//
//	//b := pqInnerData.Encode()
//	//sha1_b := sha1.Sum(b)
//	//
//	//b = append(sha1_b[:],b...)
//	//fmt.Println(hex.EncodeToString(b))
//	//
//	//e_b := rsa.Encrypt(b)
//	//fmt.Println(hex.EncodeToString(e_b))
//	//
//	//d_b := rsa.Decrypt(e_b)
//	//fmt.Println(hex.EncodeToString(d_b))
//	//
//	//dPQInnerData := &mtproto.TLPQInnerData{}
//	//dbuf := mtproto.NewDecodeBuf(d_b[24:])
//	//// dbuf.Int()
//	//err := dPQInnerData.Decode(dbuf)
//	//if err != nil {
//	//	fmt.Println(err)
//	//	return
//	//}
//	//
//	//fmt.Println(dPQInnerData)
//
//	// pqInnerData.
//	// mtproto.TLServer_DHInnerData{}
//
//	/*
//		// a1 := NewMessage(1)
//		// (*a1).Decode([]byte{})
//
//		// a2 := NewMessage(2)
//		// (*a2).Decode([]byte{})
//
//		block, _ := pem.Decode([]byte(pemBytes))
//
//		// publickey, err := x509.ParsePKIXPublicKey(block.Bytes)
//		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//
//		if err != nil {
//			fmt.Println("parse pkix key error: ", err)
//			return
//			// return nil, err
//		}
//		//if err != nil {
//		//	fmt.Println(err)
//		//	os.Exit(1)
//		//}
//
//		// convert publickey interface to rsa.PublicKey type
//		rsaPubKey := pubInterface.(*rsa.PublicKey)
//
//		fmt.Printf("Public Key N value(modulus) :  : %d\n\n", rsaPubKey.N)
//
//		fmt.Printf("Public Key E value(exponent) :  : %d\n\n", rsaPubKey.E)
//	 */
//}
//
