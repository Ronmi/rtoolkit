package conenv

import "crypto/aes"

// AESExtension creates an extension to decrypt the value which is encrypted with
// shared key using AES algorithm.
//
// The extension returns empty string, which means "unsetting" the value, if
// decrypting failed.
//
// key MUST be 16, 24 or 32 bytes string. See DESEncrypt() for implementation
// detail.
func AESExtension(key string) (ext Extension) {
	return genSharedKeyExtension(key, AESDecrypt)
}

// AESEncrypt encrypts the val using AES/CTR algorithm.
//
//  * As needed by aes.NewCipher(), key MUST be 16, 24 or 32 bytes, which selects
//    AES-128, AES-192 or AES-256.
//  * Returned code is encoded with base64.StdEncoding.
//  * val is padded with space before encrypting.
func AESEncrypt(key, val string) (code string, err error) {
	// encryption
	b, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	return blockEncrypt(val, b)
}

// AESDecrypt decrypts the val using AES/CTR algorithm.
//
//  * As needed by aes.NewCipher(), key MUST be 16, 24 or 32 bytes, which selects
//    AES-128, AES-192 or AES-256.
//  * Returned code is encoded with base64.StdEncoding.
//  * val is padded with space before encrypting.
func AESDecrypt(key, code string) (val string, err error) {
	// decryption
	b, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	return blockDecrypt(code, b)
}
