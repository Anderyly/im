/*
 * *
 *  @Author anderyly
 *  @email admin@aaayun.cc
 *  @link http://blog.aaayun.cc/
 *  @copyright Copyright (c) 2022
 *  *
 */

package ay

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func LastTime(t int) (msg string) {
	s := (int(time.Now().Unix()) - t) / 60

	switch {
	case s < 60:
		msg = strconv.Itoa(s) + "分钟前"

	case s >= 60 && s < (60*24):
		msg = strconv.Itoa(s/60) + "小时前"
	case s >= (60*24) && s < (60*24*3):
		msg = strconv.Itoa(s/24/60) + "天前"

	default:
		msg = time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
	}
	return
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

//Base64Decode  base64 解密
func Base64Decode(s string) string {
	var b []byte
	var err error
	x := len(s) * 3 % 4
	switch {
	case x == 2:
		s += "=="
	case x == 1:
		s += "="
	}
	if b, err = base64.StdEncoding.DecodeString(s); err != nil {
		return string(b)
	}
	return string(b)
}

func AuthCode(str, operation, key string, expiry int64) string {
	// 动态密匙长度，相同的明文会生成不同密文就是依靠动态密匙
	// 加入随机密钥，可以令密文无任何规律，即便是原文和密钥完全相同，加密结果也会每次不同，增大破解难度。
	// 取值越大，密文变动规律越大，密文变化 = 16 的 cKeyLength 次方
	// 当此值为 0 时，则不产生随机密钥
	cKeyLength := 1
	if len(str) < cKeyLength {
		return ""
	}

	// 密匙
	if key == "" {
		key = "#@!^5ebcQJx2Lz6GmcsqNiNHW^!@#"
	}
	key = MD5(key)

	// 密匙a会参与加解密
	keyA := MD5(key[:16])
	// 密匙b会用来做数据完整性验证
	keyB := MD5(key[16:])
	// 密匙c用于变化生成的密文
	keyC := ""
	if operation == "DECODE" {
		keyC = str[:cKeyLength]
	} else {
		sTime := MD5(time.Now().String())
		sLen := 32 - cKeyLength
		keyC = sTime[sLen:]
	}
	// 参与运算的密匙
	cryptKey := fmt.Sprintf("%s%s", keyA, MD5(keyA+keyC))
	keyLength := len(cryptKey)
	// 明文，前10位用来保存时间戳，解密时验证数据有效性，10到26位用来保存$keyB(密匙b)，解密时会通过这个密匙验证数据完整性
	// 如果是解码的话，会从第$ckey_length位开始，因为密文前$ckey_length位保存 动态密匙，以保证解密正确
	if operation == "DECODE" {
		str = strings.Replace(str, "-", "+", -1)
		str = strings.Replace(str, "_", "/", -1)
		str = strings.Replace(str, "*", "=", -1)
		strByte, err := base64.StdEncoding.DecodeString(str[cKeyLength:])
		if err != nil {
			log.Fatal(err)
		}
		str = string(strByte)
	} else {
		if expiry != 0 {
			expiry = expiry + time.Now().Unix()
		}
		tmpMd5 := MD5(str + keyB)
		str = fmt.Sprintf("%010d%s%s", expiry, tmpMd5[:16], str)
	}
	string_length := len(str)
	resdata := make([]byte, 0, string_length)
	var rndkey, box [256]int
	// 产生密匙簿
	j := 0
	a := 0
	i := 0
	tmp := 0
	for i = 0; i < 256; i++ {
		rndkey[i] = int(cryptKey[i%keyLength])
		box[i] = i
	}
	// 用固定的算法，打乱密匙簿，增加随机性，好像很复杂，实际上并不会增加密文的强度
	for i = 0; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		tmp = box[i]
		box[i] = box[j]
		box[j] = tmp
	}
	// 核心加解密部分
	a = 0
	j = 0
	tmp = 0
	for i = 0; i < string_length; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp = box[a]
		box[a] = box[j]
		box[j] = tmp
		// 从密匙簿得出密匙进行异或，再转成字符
		resdata = append(resdata, byte(int(str[i])^box[(box[a]+box[j])%256]))
	}
	result := string(resdata)
	if operation == "DECODE" {
		// substr($result, 0, 10) == 0 验证数据有效性
		// substr($result, 0, 10) - time() > 0 验证数据有效性
		// substr($result, 10, 16) == substr(md5(substr($result, 26).$keyB), 0, 16) 验证数据完整性
		// 验证数据有效性，请看未加密明文的格式
		frontTen, _ := strconv.ParseInt(result[:10], 10, 0)
		if (frontTen == 0 || frontTen-time.Now().Unix() > 0) && result[10:26] == MD5(result[26:] + keyB)[:16] {
			return result[26:]
		} else {
			return ""
		}
	} else {
		// 把动态密匙保存在密文里，这也是为什么同样的明文，生产不同密文后能解密的原因
		// 因为加密后的密文可能是一些特殊字符，复制过程可能会丢失，所以用base64编码
		result = keyC + base64.StdEncoding.EncodeToString([]byte(result))
		result = strings.Replace(result, "+", "-", -1)
		result = strings.Replace(result, "/", "_", -1)
		result = strings.Replace(result, "=", "*", -1)
		return result
	}
}
