package helper

// 常用工具

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var mykey = []byte("gin_gorm_framework_key")

// GetMd5
// 生成MD5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))

}

// GenerateToken
// 生成token
func GenerateToken(identity, name string, IsAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        IsAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(mykey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析token
func AnalyseToken(tokenstring string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenstring, userClaim, func(t *jwt.Token) (interface{}, error) {
		return mykey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims.Valid {
		// fmt.Println(userClaim)
		return userClaim, nil
	} else {
		return nil, fmt.Errorf("analyze token error:%v", err)
	}
}

// 发送验证码

func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "zhangkun <zhangkun@lhcis.com>"
	e.To = []string{toUserEmail}
	// e.Bcc = []string{"test_bcc@example.com"}
	// e.Cc = []string{"test_cc@example.com"}
	e.Subject = "验证码已发送请查收"
	// e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	// err := e.Send("smtp.exmail.qq.com:465", smtp.PlainAuth("",
	// 	"zhangkun@lhcis.com",
	// 	"5MJmGLienHk83a5t",
	// 	"smtp.exmail.qq.com"))
	// 返回EOF时候关闭SSL重试
	err := e.SendWithTLS("smtp.exmail.qq.com:465", smtp.PlainAuth("",
		"zhangkun@lhcis.com",
		"5MJmGLienHk83a5t",
		"smtp.exmail.qq.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.exmail.qq.com"})
	return err
}

// GetUUID
// 生成UUID
func GetUUID() string {
	return uuid.NewV4().String()
}

// 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10)) // 0-9随机数
	}
	return s
}

// 代码保存
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil

}
