package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"github.com/golang-jwt/jwt"
	uuid2 "github.com/google/uuid"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func NewUuid() string {
	uuid := uuid2.New()
	return uuid.String()
}

// NowUnix 当前时间戳
func NowUnix() int {
	var sysTimeLocation, _ = time.LoadLocation("Asia/Shanghai")
	return int(time.Now().In(sysTimeLocation).Unix())
}

// FormatFromUnixTime 将时间戳转为 yyyy-mm-dd H:i:s 格式
func FormatFromUnixTime(t int64) string {
	cstSh, _ := time.LoadLocation("Asia/Shanghai") //上海
	if t > 0 {
		return time.Unix(t, 0).In(cstSh).Format(constant.SysTimeFormat)
	} else {
		return time.Now().In(cstSh).Format(constant.SysTimeFormat)
	}
}

// FormatFromUnixTimeShort 将时间戳转为 yyyy-mm-dd 格式
func FormatFromUnixTimeShort(t int64) string {
	cstSh, _ := time.LoadLocation("Asia/Shanghai") //上海
	if t > 0 {
		return time.Unix(t, 0).In(cstSh).Format(constant.SysTimeFormatShort)
	} else {
		return time.Now().In(cstSh).Format(constant.SysTimeFormatShort)
	}
}

// ParseTime 将字符串转成时间
func ParseTime(str string) (time.Time, error) {
	var sysTimeLocation, _ = time.LoadLocation("Asia/Shanghai")
	return time.ParseInLocation(constant.SysTimeFormat, str, sysTimeLocation)
}

// Random 得到一个随机数
func Random(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if max < 1 {
		return r.Int()
	} else {
		return r.Intn(max)
	}
}

// encrypt 对一个字符串进行加密
func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// decrypt 对一个字符串进行解密
func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext is too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 在预定义字符前添加 \
// ' " \
// http://www.ruanyifeng.com/blog/2007/10/ascii_unicode_and_utf-8.html
func AddSlashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// StripsSlashes 删除
func StripsSlashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	for i := 0; i < len(strRune); i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

// Ip4toInt 将字符串的 IP 转化为数字
func Ip4toInt(ip string) int64 {
	bits := strings.Split(ip, ".")
	var sum int64 = 0
	if len(bits) == 4 {
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])
		sum += int64(b0) << 24
		sum += int64(b1) << 16
		sum += int64(b2) << 8
		sum += int64(b3)
	}
	return sum
}

// NextDayDuration 得到当前时间到下一天零点的延时
func NextDayDuration() time.Duration {
	year, month, day := time.Now().Add(time.Hour * 24).Date()
	var sysTimeLocation, _ = time.LoadLocation("Asia/Shanghai")
	next := time.Date(year, month, day, 0, 0, 0, 0, sysTimeLocation)
	return next.Sub(time.Now())
}

func NextDayTime() time.Time {
	year, month, day := time.Now().Add(time.Hour * 24).Date()
	var sysTimeLocation, _ = time.LoadLocation("Asia/Shanghai")
	next := time.Date(year, month, day, 0, 0, 0, 0, sysTimeLocation)
	return next
}

// isLittleEndian 判断当前系统中的字节序类型是否是小端字节序
func isLittleEndian() bool {
	var i int = 0x1
	bs := (*[int(unsafe.Sizeof(0))]byte)(unsafe.Pointer(&i))
	return bs[0] == 0
}

// GetInt64 从接口类型安全获取到int64，d 是默认值
func GetInt64(i interface{}, d int64) int64 {
	if i == nil {
		return d
	}
	switch i.(type) {
	case string:
		num, err := strconv.Atoi(i.(string))
		if err != nil {
			return d
		} else {
			return int64(num)
		}
	case []byte:
		bits := i.([]byte)
		if len(bits) == 8 {
			if isLittleEndian() {
				return int64(binary.LittleEndian.Uint64(bits))
			} else {
				return int64(binary.BigEndian.Uint64(bits))
			}
		} else if len(bits) <= 4 {
			num, err := strconv.Atoi(string(bits))
			if err != nil {
				return d
			} else {
				return int64(num)
			}
		}
	case uint:
		return int64(i.(uint))
	case uint8:
		return int64(i.(uint8))
	case uint16:
		return int64(i.(uint16))
	case uint32:
		return int64(i.(uint32))
	case uint64:
		return int64(i.(uint64))
	case int:
		return int64(i.(int))
	case int8:
		return int64(i.(int8))
	case int16:
		return int64(i.(int16))
	case int32:
		return int64(i.(int32))
	case int64:
		return i.(int64)
	case float32:
		return int64(i.(float32))
	case float64:
		return int64(i.(float64))
	}
	return d
}

// GetString 从接口安全获取到字符串类型
func GetString(str interface{}, d string) string {
	if str == nil {
		return d
	}
	switch str.(type) {
	case string:
		return str.(string)
	case []byte:
		return string(str.([]byte))
	}
	return fmt.Sprintf("%s", str)
}

// GetInt64FromMap 从map中得到指定的key
func GetInt64FromMap(dm map[string]interface{}, key string, d int64) int64 {
	data, ok := dm[key]
	if !ok {
		return d
	}
	return GetInt64(data, d)
}

func GetStringFromMap(dm map[string]interface{}, key string, d string) string {
	data, ok := dm[key]
	if !ok {
		return d
	}
	return GetString(data, d)
}

func GetTodayIntDay() int {
	y, m, d := time.Now().Date()
	strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	day, _ := strconv.Atoi(strDay)
	return day
}

// JWTClaims 自定义格式内容
type JWTClaims struct {
	UserID         uint   `json:"user_id"`
	UserName       string `json:"user_name"`
	StandardClaims jwt.StandardClaims
}

func (j JWTClaims) Valid() error {
	return nil
}

// GenerateJwtToken 生成token
func GenerateJwtToken(secret string, issuer string, userId uint, userName string) (string, error) {
	hmacSampleSecret := []byte(secret) //密钥，不能泄露
	token := jwt.New(jwt.SigningMethodHS256)
	nowTime := time.Now().Unix()
	token.Claims = JWTClaims{
		UserID:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: nowTime,                                             // 签名生效时间
			ExpiresAt: time.Now().Add(constant.TokenExpireDuration).Unix(), // 签名过期时间
			Issuer:    issuer,                                              // 签名颁发者
		},
	}
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

// ParseJwtToken 解析token
func ParseJwtToken(tokenString string, secret string) (*JWTClaims, error) {
	var hmacSampleSecret = []byte(secret)
	//前面例子生成的token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	claims := token.Claims.(*JWTClaims)
	return claims, nil
}
