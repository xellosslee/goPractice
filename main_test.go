package main_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/xlzd/gotp"
	"golang.org/x/crypto/bcrypt"
	"gs.lee.was/util"
)

func TestChannel(t *testing.T) {
	// // 채널 테스트
	// test := make(chan int)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		log.Print(i)
	// 	}
	// 	test <- 1
	// }()
	// var test2 int
	// test2 = <-test
	// log.Print("complete ", test2)
}

func TestEncrypt(t *testing.T) {
	log.Print("암호화 테스트")

	salt := util.RandString(10)
	log.Print(salt)

	hash, err := bcrypt.GenerateFromPassword([]byte("test"+salt), bcrypt.MinCost)
	if err != nil {
		log.Panic(err)
	}
	log.Print(string(hash))
	err = bcrypt.CompareHashAndPassword(hash, []byte("test"+salt))
	if err != nil {
		log.Print("not match password")
	}

	// sha256
	data := "sh12345!!"
	h := sha256.New()
	h.Write([]byte(data))
	md := h.Sum(nil)
	mdStr := hex.EncodeToString(md)

	fmt.Println(mdStr)
}

func TestTOTP(t *testing.T) {
	totp := gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO")
	totp.Now()          // current otp '123456'
	totp.At(1524486261) // otp of timestamp 1524486261 '123456'

	totp.Verify("492039", 1524486261) // true
	totp.Verify("492039", 1520000000) // false

	// generate a provisioning uri
	url := totp.ProvisioningUri("gslee", "cdms")
	t.Log(url)
	// otpauth://totp/issuerName:demoAccountName?secret=4S62BZNFXXSZLCRO&issuer=issuerName
}

func TestSendMail(t *testing.T) {
	util.SendEmail("gs.lee@cndfactory.com", "test", "join.html", struct{}{}, []string{})
}
