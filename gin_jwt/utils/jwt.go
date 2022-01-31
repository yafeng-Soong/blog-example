package utils

import (
	"errors"
	"gin_jwt/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	User model.UserInfo
	jwt.StandardClaims
}

const (
	TokenExpireDuration = time.Hour * 2
	M                   = time.Minute * 5
)

var MySecret = []byte("yoursecret") // 生成签名的密钥

func GenerateToken(userInfo model.UserInfo) (string, error) {
	expirationTime := time.Now().Add(M) // 两个小时有效期
	claims := &MyClaims{
		User: userInfo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "yourname",
		},
	}
	// 生成Token，指定签名算法和claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名
	if tokenString, err := token.SignedString(MySecret); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}

}

func RenewToken(claims *MyClaims) (string, error) {
	// 若token过期不超过10分钟则给它续签
	if withinLimit(claims.ExpiresAt, 600) {
		return GenerateToken(claims.User)
	}
	return "", errors.New("登录已过期")
}

func ParseToken(tokenString string) (*MyClaims, error) {
	claims := &MyClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	// 若token只是过期claims是有数据的，若token无法解析claims无数据
	return claims, err
}

func ParseToken2(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token无法解析")
}

// 计算过期时间是否超过l
func withinLimit(s int64, l int64) bool {
	e := time.Now().Unix()
	// println(e - s)
	return e-s < l
}
