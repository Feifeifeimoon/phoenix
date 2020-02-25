package common

/*
	jwt JSON Web Token
*/
import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	//
	Secret = "salt"
	//token过期时间 (s)
	ExpireTime = 60 * 30
	//token过期
	TokenExpired = errors.New("Token is expired")
	TokenInvalid = errors.New("Token is invalid")
)

// token里面添加用户信息，验证token后可能会用到用户信息
type JWTClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

// 生成token
func GenerateToken(id string) (string, error) {
	claims := JWTClaims{UserID: id}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

//验证token
func VerifyToken(strToken string) (string, error) {
	function := func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	}
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, function)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", TokenExpired
			} else {
				return "", TokenInvalid
			}
		}
		return "", TokenInvalid
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}
	return "", TokenInvalid
}
