package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type JwtJson struct {
	ProjectId           string `json:"project_id"`
	KeyId               string `json:"key_id"`
	PrivateKey          string `json:"private_key"`
	SubAccount          string `json:"sub_account"`
	AuthUri             string `json:"auth_uri"`
	TokenUri            string `json:"token_uri"`
	AuthProviderCertUri string `json:"auth_provider_cert_uri"`
	ClientCertUri       string `json:"client_cert_uri"`
}

func GenerateJwtToken(jwtJson JwtJson) string {
	// 解码私钥
	privateKeyData := []byte(jwtJson.PrivateKey)
	block, _ := pem.Decode(privateKeyData)
	if block == nil || !strings.HasPrefix(block.Type, "PRIVATE KEY") {
		log.Fatalf("无效的私钥格式")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		log.Fatalf("私钥不是 RSA 私钥")
	}

	// 创建 JWT 载荷
	now := time.Now().UTC().Unix()
	payload := jwt.MapClaims{
		"aud": jwtJson.AuthUri,
		"iss": jwtJson.SubAccount,
		"exp": now + 300,
		"iat": now,
	}

	// 创建 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, payload)
	token.Header["kid"] = jwtJson.KeyId
	// 使用私钥进行签名
	tokenString, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		log.Fatalf("签名 JWT 失败: %v", err)
	}
	return tokenString
}

func readJwtJsonFile() JwtJson {
	// Open the JSON file
	file, err := os.Open("/app/data/private.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the JSON file
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal JSON data into the struct
	var config JwtJson
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return config
}
