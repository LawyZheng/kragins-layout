package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// CBCEncrypt: 对字符串用 secretKey 进行 CBC 加密
// 参数
//
//	text: 需要加密的字符串
//	secretKey: 用于加密的key
//
// 返回
//
//	加密结果
//	错误
func CBCEncrypt(plaintext []byte, secretKey string) (string, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	//key, _ := hex.DecodeString(secretKey)
	key := []byte(secretKey)

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2.
	// Here we'll padding it
	//if len(plaintext)%aes.BlockSize != 0 {
	plaintext = pkcs7Padding(plaintext, aes.BlockSize)
	//}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	//ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	//iv := ciphertext[:aes.BlockSize]
	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	return "", errors.Wrap(err)
	//}

	// 昌平版本，放弃随机iv， 用key的前16个block作为iv
	iv := key[:block.BlockSize()]

	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func CBCDecrypt(crypt string, secretKey string) (string, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	//key, _ := hex.DecodeString(secretKey)
	key := []byte(secretKey)
	//ciphertext, _ := hex.DecodeString(crypt)
	ciphertext, err := base64.StdEncoding.DecodeString(crypt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	//if len(ciphertext) < aes.BlockSize {
	//	return "", errors.Errorf("ciphertext too short")
	//}
	//iv := ciphertext[:aes.BlockSize]
	//ciphertext = ciphertext[aes.BlockSize:]

	// 昌平版本iv从key中取
	iv := key[:aes.BlockSize]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	origin := make([]byte, len(ciphertext))
	mode.CryptBlocks(origin, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	return string(pkcs7UnPadding(origin)), nil
}

// pkcs7Padding: 补码,CBC 模式需要 被加密的字符串为 block的整数倍 https://tools.ietf.org/html/rfc5246#section-6.2.3.2.
// 参数
//
//	cipherText: 需要被加密的字符串
//	blockSize: 块大小
//
// 返回
//
//	补码后的字符串
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	// 查看缺失的数量
	padding := blockSize - len(ciphertext)%blockSize
	// 用该数量作为补码重复补充至末尾
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// pkcs7UnPadding: 去除补码
// 参数
//
//	cipherText: 需要去除补码的字符串
//
// 返回
//
//	原字符串
func pkcs7UnPadding(origData []byte) []byte {
	// 获取最末尾的数据
	padding := origData[len(origData)-1]
	// 根据获取到的补码长度来截取
	return origData[:(len(origData) - int(padding))]
}
