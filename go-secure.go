package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// read file as byte
func ReadFileAsByte(filePath string) ([]byte, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WriteFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var output, key string
	var mode int

	const (
		ENCRYPT = iota
		DECRYPT
	)

	flag.IntVar(&mode, "m", ENCRYPT, "mode: 0: encrypt(DEFAULT), 1: decrypt")
	flag.StringVar(&output, "o", "output.txt", "output file path")
	flag.StringVar(&key, "k", "", "key to decrypt/encrypt in base64 if empty random string will be generated")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	filePath := flag.Arg(0)

	input, err := ReadFileAsByte(filePath)
	if err != nil {
		panic(err)
	}

	var secretb64 string
	if key == "" {
		secretb64, err = GenerateRandomSecretToBase64(32)
		if err != nil {
			panic(err)
		}
		log.Println("Please take a note of the secret:", secretb64)
	} else {
		secretb64 = key
	}

	secret, err := base64.StdEncoding.DecodeString(secretb64)
	if err != nil {
		panic(err)
	}

	switch mode {
	case ENCRYPT:
		encrypted, err := Encrypt(input, secret)
		if err != nil {
			panic(err)
		}

		err = WriteFile(output, encrypted)
		if err != nil {
			panic(err)
		}
	case DECRYPT:
		plainText, err := Decrypt(input, []byte(secret))
		if err != nil {
			panic(err)
		}

		err = WriteFile(output, plainText)
		if err != nil {
			panic(err)
		}
	}

}

func GenerateRandomSecretToBase64(len int) (result string, err error) {
	secret := make([]byte, len)
	_, err = rand.Read(secret)
	if err != nil {
		return "", err
	}
	result = base64.StdEncoding.EncodeToString(secret)
	return result, nil
}

func Encrypt(plainText, secret []byte) (cipherText []byte, errs error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		errs = err
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		errs = err
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		errs = err
		return
	}

	cipherText = gcm.Seal(
		nonce,
		nonce,
		plainText,
		nil)
	return
}

func Decrypt(cipherText, secret []byte) (plainText []byte, errs error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		errs = err
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		errs = err
		return
	}
	nonceSize := gcm.NonceSize()

	if len(cipherText) < nonceSize {
		errs = errors.New("Not enough size on ciphertext")
		return
	}
	nonce, cipherTextOnly := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err = gcm.Open(nil, nonce, cipherTextOnly, nil)
	if err != nil {
		errs = err
		return
	}
	return
}
