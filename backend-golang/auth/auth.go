package auth

import (
	"errors"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	jwtReq "github.com/dgrijalva/jwt-go/request"
)

type jwtClaims struct {
	UserId   int64
	Admin    bool
	PassPart string
	jwt.StandardClaims
}

var (
	ErrJwtClaimsAssertFailed = errors.New("couldn't assert claim type")
	jwtKeyFunc               = func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return common.Config.JwtSecret, nil
	}
)

func GetUserIdAndPassPartAndAdminFromRequest(r *http.Request) (userID int64, passPart string, isAdmin bool, err error) {
	token, err := jwtReq.ParseFromRequestWithClaims(r, jwtReq.AuthorizationHeaderExtractor, &jwtClaims{}, jwtKeyFunc)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		err = ErrJwtClaimsAssertFailed
		return
	}

	userID = claims.UserId
	passPart = claims.PassPart
	isAdmin = claims.Admin

	return
}
