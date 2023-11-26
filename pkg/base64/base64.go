package base64

import "encoding/base64"

// DecodeString 将字符串两次base64解密
func DecodeString(s string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	key, err = base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return "", err
	}
	return string(key), nil
}

// EncodeString 将字符串两次base64加密
func EncodeString(s string) string {
	key := base64.StdEncoding.EncodeToString([]byte(s))
	key1 := base64.StdEncoding.EncodeToString([]byte(key))
	return key1
}
