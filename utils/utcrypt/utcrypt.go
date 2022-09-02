package utcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/koinworks/asgard-heimdal/libs/serror"
	"github.com/koinworks/asgard-heimdal/utils/utstring"
)

var key = "kVBkcfp3ivy3GqzSn62Dn2sW27X408VN"
var size = 16

type dfCrypt struct {
	IV    string `json:"iv"`
	Value string `json:"value"`
	Mac   string `json:"mac"`
}

func Decrypt(v string) (string, serror.SError) {
	base, _ := base64.StdEncoding.DecodeString(v)

	var err error

	if err != nil {
		return "", serror.New("Decrypt failed, 0x001")
	}

	jd := dfCrypt{}
	err = json.Unmarshal(base, &jd)
	if err != nil {
		return "", serror.NewFromError(err)
	}

	hs := hmac.New(md5.New, []byte(jd.IV))

	_, err = hs.Write([]byte(jd.Value))
	if err != nil {
		return "", serror.NewFromErrorc(err, "Failed to write")
	}

	hsh := hs.Sum(nil)

	if hex.EncodeToString(hsh) != jd.Mac {
		return "", serror.New("Decrypt failed, 0x002")
	}

	iv, err := base64.StdEncoding.DecodeString(jd.IV)
	if err != nil {
		return "", serror.New("Decrypt failed, 0x003")
	}
	ct, err := base64.StdEncoding.DecodeString(jd.Value)
	if err != nil {
		return "", serror.New("Decrypt failed, 0x004")
	}

	k := []byte(key)
	blc, err := aes.NewCipher(k)
	if err != nil {
		return "", serror.New("Decrypt failed, 0x005")
	}

	if len(ct) < aes.BlockSize {
		return "", serror.New("Decrypt failed, 0x006")
	} else if len(ct)%aes.BlockSize != 0 {
		return "", serror.New("Decrypt failed, 0x007")
	}

	m := cipher.NewCBCDecrypter(blc, iv)
	m.CryptBlocks(ct, ct)
	ct = pkcs5UnPadding(ct)
	return string(ct), nil
}

func Encrypt(v string) (string, serror.SError) {
	vb := []byte(v)
	vb = pkcs5Padding(vb, size)

	if len(vb)%aes.BlockSize != 0 {
		return "", serror.New("Encrypt failed, 0x101")
	}

	k := []byte(key)
	b, err := aes.NewCipher(k)
	if err != nil {
		return "", serror.New("Encrypt failed, 0x102")
	}

	material := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	iv := []byte(utstring.RandString(size, material))

	m := cipher.NewCBCEncrypter(b, iv)
	m.CryptBlocks(vb, vb)

	j := dfCrypt{
		IV:    base64.StdEncoding.EncodeToString(iv),
		Value: base64.StdEncoding.EncodeToString(vb),
	}

	hs := hmac.New(md5.New, []byte(j.IV))
	_, err = hs.Write([]byte(j.Value))
	if err != nil {
		return "", serror.NewFromErrorc(err, "Failed to write")
	}
	hsh := hs.Sum(nil)

	j.Mac = hex.EncodeToString(hsh)

	js, err := json.Marshal(j)
	if err != nil {
		return "", serror.New("Encrypt failed, 0x103")
	}

	return base64.StdEncoding.EncodeToString(js), nil
}

func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
