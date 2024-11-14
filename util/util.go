package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/xlzd/gotp"
	"golang.org/x/crypto/bcrypt"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	md := h.Sum(nil)
	return hex.EncodeToString(md)
}

func PasswordMake(p string) ([]byte, string, error) {
	s := gotp.RandomSecret(16)
	hash, err := bcrypt.GenerateFromPassword([]byte(p+s), bcrypt.MinCost)
	return hash, s, err
}

func PasswordMakeWithSalt(p string, s string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p+s), bcrypt.MinCost)
	return hash, err
}

func PasswordCheck(h string, p string, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p+s))
}

func OtpUrlMake(u string, s string) string {
	otp := gotp.NewDefaultTOTP(s).ProvisioningUri(u, "Happybutler")
	// 구글 서비스 종료 된 url
	// return "https://www.google.com/chart?chs=200x200&chld=M%%7C0&cht=qr&chl=" + url.QueryEscape(otp)
	return "https://api.qrserver.com/v1/create-qr-code/?size=250x250&data=" + url.QueryEscape(otp)
}

func GenerateRandomId(prefix string, appendMillisecond int) string {
	now := time.Now()
	if appendMillisecond > 0 {
		now = now.Add(time.Millisecond * time.Duration(appendMillisecond))
	}
	// configs.Log.Info(time.Now().UnixNano())
	i := now.UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond)
	// configs.Log.Info(i)
	r := rand.Intn(100)
	// configs.Log.Info(fmt.Sprintf("P%09d%02d", i, r))
	return fmt.Sprintf("%s%s%02d", prefix, strconv.Itoa(int(i))[5:], r)
}
