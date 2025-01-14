package dongle

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	tripleDesKey = "0123456789abcdefghijklmn"
	tripleDesIV  = "12345678"
)

var tripleDesTests = []struct {
	mode     cipherMode
	padding  cipherPadding
	input    string
	toHex    string
	toBase64 string
}{
	{CBC, PKCS7, "", "", ""},
	{CBC, No, "hello world, go!", "7e9194cc827a325d49111eaa503110fe", "fpGUzIJ6Ml1JER6qUDEQ/g=="},
	{CBC, Zero, "hello world", "7e9194cc827a325ddaee992a89c5cd8d", "fpGUzIJ6Ml3a7pkqicXNjQ=="},
	{CBC, PKCS5, "hello world", "7e9194cc827a325db9c765859716cc97", "fpGUzIJ6Ml25x2WFlxbMlw=="},
	{CBC, PKCS7, "hello world", "7e9194cc827a325db9c765859716cc97", "fpGUzIJ6Ml25x2WFlxbMlw=="},
	{CBC, AnsiX923, "hello world", "7e9194cc827a325d2793bb48a7971825", "fpGUzIJ6Ml0nk7tIp5cYJQ=="},
	{CBC, ISO97971, "hello world", "7e9194cc827a325d89d4f50218d6e511", "fpGUzIJ6Ml2J1PUCGNblEQ=="},

	{CFB, PKCS7, "", "", ""},
	{CFB, No, "hello world, go!", "f52dc896ceff0ecc9393c19d4e1e7591", "9S3Ils7/DsyTk8GdTh51kQ=="},
	{CFB, Zero, "hello world", "f52dc896ceff0ecc9393c1b16e791ab0", "9S3Ils7/DsyTk8GxbnkasA=="},
	{CFB, PKCS5, "hello world", "f52dc896ceff0ecc9393c1b46b7c1fb5", "9S3Ils7/DsyTk8G0a3wftQ=="},
	{CFB, PKCS7, "hello world", "f52dc896ceff0ecc9393c1b46b7c1fb5", "9S3Ils7/DsyTk8G0a3wftQ=="},
	{CFB, AnsiX923, "hello world", "f52dc896ceff0ecc9393c1b16e791ab5", "9S3Ils7/DsyTk8GxbnkatQ=="},
	{CFB, ISO97971, "hello world", "f52dc896ceff0ecc9393c1316e791ab0", "9S3Ils7/DsyTk8ExbnkasA=="},

	{OFB, PKCS7, "", "", ""},
	{OFB, No, "hello world, go!", "f52dc896ceff0eccf869cb59735c3766", "9S3Ils7/Dsz4actZc1w3Zg=="},
	{OFB, Zero, "hello world", "f52dc896ceff0eccf869cb75533b5847", "9S3Ils7/Dsz4act1UztYRw=="},
	{OFB, PKCS5, "hello world", "f52dc896ceff0eccf869cb70563e5d42", "9S3Ils7/Dsz4actwVj5dQg=="},
	{OFB, PKCS7, "hello world", "f52dc896ceff0eccf869cb70563e5d42", "9S3Ils7/Dsz4actwVj5dQg=="},
	{OFB, AnsiX923, "hello world", "f52dc896ceff0eccf869cb75533b5842", "9S3Ils7/Dsz4act1UztYQg=="},
	{OFB, ISO97971, "hello world", "f52dc896ceff0eccf869cbf5533b5847", "9S3Ils7/Dsz4acv1UztYRw=="},

	{CTR, PKCS7, "", "", ""},
	{CTR, No, "hello world, go!", "f52dc896ceff0ecc366b2281038f6f7f", "9S3Ils7/Dsw2ayKBA49vfw=="},
	{CTR, Zero, "hello world", "f52dc896ceff0ecc366b22ad23e8005e", "9S3Ils7/Dsw2ayKtI+gAXg=="},
	{CTR, PKCS5, "hello world", "f52dc896ceff0ecc366b22a826ed055b", "9S3Ils7/Dsw2ayKoJu0FWw=="},
	{CTR, PKCS7, "hello world", "f52dc896ceff0ecc366b22a826ed055b", "9S3Ils7/Dsw2ayKoJu0FWw=="},
	{CTR, AnsiX923, "hello world", "f52dc896ceff0ecc366b22ad23e8005b", "9S3Ils7/Dsw2ayKtI+gAWw=="},
	{CTR, ISO97971, "hello world", "f52dc896ceff0ecc366b222d23e8005e", "9S3Ils7/Dsw2ayItI+gAXg=="},

	{ECB, PKCS7, "", "", ""},
	{ECB, No, "hello world, go!", "b8097975c76319c623be7c7aa6e0f3fc", "uAl5dcdjGcYjvnx6puDz/A=="},
	{ECB, Zero, "hello world", "b8097975c76319c61971a986e579cdf9", "uAl5dcdjGcYZcamG5XnN+Q=="},
	{ECB, PKCS5, "hello world", "b8097975c76319c6172687e0d90fd4d1", "uAl5dcdjGcYXJofg2Q/U0Q=="},
	{ECB, PKCS7, "hello world", "b8097975c76319c6172687e0d90fd4d1", "uAl5dcdjGcYXJofg2Q/U0Q=="},
	{ECB, AnsiX923, "hello world", "b8097975c76319c6d98a83ce5ec18698", "uAl5dcdjGcbZioPOXsGGmA=="},
	{ECB, ISO97971, "hello world", "b8097975c76319c66b1c75b91028ca62", "uAl5dcdjGcZrHHW5ECjKYg=="},
}

func Test3Des_Encrypt_String(t *testing.T) {
	for index, test := range tripleDesTests {
		raw := Decode.FromString(test.toHex).ByHex().ToString()
		e := Encrypt.FromString(test.input).By3Des(getCipher(test.mode, test.padding, tripleDesKey, tripleDesIV))

		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_raw_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, raw, e.ToRawString())
		})
		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_hex_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, test.toHex, e.ToHexString())
		})
		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_base64_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, test.toBase64, e.ToBase64String())
		})
	}
}

func Test3Des_Decrypt_String(t *testing.T) {
	for index, test := range tripleDesTests {
		raw := Decode.FromString(test.toHex).ByHex().ToString()
		e := Decrypt.FromRawString(raw).By3Des(getCipher(test.mode, test.padding, tripleDesKey, tripleDesIV))

		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_raw_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, test.input, e.ToString())
			assert.Equal(t, test.input, fmt.Sprintf("%s", e))
		})
	}

	for index, test := range tripleDesTests {
		e := Decrypt.FromHexString(test.toHex).By3Des(getCipher(test.mode, test.padding, tripleDesKey, tripleDesIV))

		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_hex_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, test.input, e.ToString())
			assert.Equal(t, test.input, fmt.Sprintf("%s", e))
		})
	}

	for index, test := range tripleDesTests {
		e := Decrypt.FromBase64String(test.toBase64).By3Des(getCipher(test.mode, test.padding, tripleDesKey, tripleDesIV))

		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_base64_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, test.input, e.ToString())
			assert.Equal(t, test.input, fmt.Sprintf("%s", e))
		})
	}
}

func Test3Des_Encrypt_Bytes(t *testing.T) {
	for index, test := range tripleDesTests {
		raw := Decode.FromBytes([]byte(test.toHex)).ByHex().ToBytes()
		e := Encrypt.FromBytes([]byte(test.input)).By3Des(getCipher(test.mode, test.padding, []byte(tripleDesKey), []byte(tripleDesIV)))

		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_raw_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, raw, e.ToRawBytes())
		})
		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_hex_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, []byte(test.toHex), e.ToHexBytes())
		})
		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_base64_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, []byte(test.toBase64), e.ToBase64Bytes())
		})
	}
}

func Test3Des_Decrypt_Bytes(t *testing.T) {
	for index, test := range tripleDesTests {
		raw := Decode.FromBytes([]byte(test.toHex)).ByHex().ToBytes()
		e := Decrypt.FromRawBytes(raw).By3Des(getCipher(test.mode, test.padding, []byte(tripleDesKey), []byte(tripleDesIV)))

		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_raw_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, []byte(test.input), e.ToBytes())
		})
	}

	for index, test := range tripleDesTests {
		e := Decrypt.FromHexBytes([]byte(test.toHex)).By3Des(getCipher(test.mode, test.padding, []byte(tripleDesKey), []byte(tripleDesIV)))

		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_hex_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, []byte(test.input), e.ToBytes())
		})
	}

	for index, test := range tripleDesTests {
		e := Decrypt.FromBase64Bytes([]byte(test.toBase64)).By3Des(getCipher(test.mode, test.padding, []byte(tripleDesKey), []byte(tripleDesIV)))

		t.Run(fmt.Sprintf(string(test.mode)+"_"+string(test.padding)+"_base64_test_%d", index), func(t *testing.T) {
			assert.Nil(t, e.Error)
			assert.Equal(t, []byte(test.input), e.ToBytes())
		})
	}
}

func Test3Des_Key_Error(t *testing.T) {
	e := Encrypt.FromString("hello world").By3Des(getCipher(CBC, PKCS7, "xxxx", tripleDesIV))
	assert.Equal(t, invalid3DesKeyError(), e.Error)

	d := Decrypt.FromHexString("68656c6c6f20776f726c64").By3Des(getCipher(CBC, PKCS7, "xxxx", tripleDesIV))
	assert.Equal(t, invalid3DesKeyError(), d.Error)
}

func Test3Des_IV_Error(t *testing.T) {
	e := Encrypt.FromString("hello world").By3Des(getCipher(OFB, PKCS7, tripleDesKey, "xxxx"))
	assert.Equal(t, invalid3DesIVError(), e.Error)

	d := Decrypt.FromHexString("68656c6c6f20776f726c64").By3Des(getCipher(CBC, PKCS7, tripleDesKey, "xxxx"))
	assert.Equal(t, invalid3DesIVError(), d.Error)
}

func Test3Des_Mode_Error(t *testing.T) {
	e := Encrypt.FromString("hello world").By3Des(getCipher("xxxx", PKCS7, tripleDesKey, tripleDesIV))
	assert.Equal(t, invalidModeError("xxxx"), e.Error)

	d := Decrypt.FromHexString("7e9194cc827a325db9c765859716cc97").By3Des(getCipher("xxxx", PKCS7, tripleDesKey, tripleDesIV))
	assert.Equal(t, invalidModeError("xxxx"), d.Error)
}

func Test3Des_Padding_Error(t *testing.T) {
	e := Encrypt.FromString("hello world").By3Des(getCipher(CFB, "xxxx", tripleDesKey, tripleDesIV))
	assert.Equal(t, invalidPaddingError("xxxx"), e.Error)

	d := Decrypt.FromHexString("7e9194cc827a325db9c765859716cc97").By3Des(getCipher(CBC, "xxxx", tripleDesKey, tripleDesIV))
	assert.Equal(t, invalidPaddingError("xxxx"), d.Error)
}

func Test3Des_Src_Error(t *testing.T) {
	e := Encrypt.FromString("hello world").By3Des(getCipher(CFB, No, tripleDesKey, tripleDesIV))
	assert.Equal(t, invalid3DesSrcError(), e.Error)

	d := Decrypt.FromHexString("68656c6c6f20776f726c64").By3Des(getCipher(CBC, PKCS7, tripleDesKey, tripleDesIV))
	assert.Equal(t, invalid3DesSrcError(), d.Error)
}

func Test3Des_Decoding_Error(t *testing.T) {
	d1 := Decrypt.FromHexBytes([]byte("xxxx")).By3Des(getCipher(CTR, Zero, []byte(tripleDesKey), []byte(tripleDesIV)))
	assert.Equal(t, invalidDecodingError("hex"), d1.Error)

	d2 := Decrypt.FromBase64Bytes([]byte("xxxxxx")).By3Des(getCipher(CFB, PKCS7, []byte(tripleDesKey), []byte(tripleDesIV)))
	assert.Equal(t, invalidDecodingError("base64"), d2.Error)
}
