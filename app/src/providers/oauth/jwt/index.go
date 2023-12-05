package jwtAuth

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/types"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT secret, same must be used in api-gateway to authenticate
var jwtSecret = []byte("yOushaLlnotPass")

type JWTClaims struct {
	*jwt.RegisteredClaims
	Data *types.SessionUser `json:"data"`
}

//
// TODO - add logic to verify token
//

/*
To add custom `data` struct make sure to pass `&data`
where `GenerateToken` is called.

The mapping only works if the props in the struct are public
i.e Capitalized

	ex: type user struct{
		Username string `json:"username"`
		id string // Not accessible and will not be added to jwt token
	}

referred: https://medium.com/@nooraldinahmed/very-basic-jwt-authentication-with-golang-3516b21c2740
*/
func GenerateToken(data *types.SessionUser) (string, error) {
	logger.Logger.Log("jwt.GenerateToken: Generating JWT token for data: ", data)
	// Get the token instance with the Signing method
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	exp := time.Now().Add((time.Hour + 24) * 365)
	// Add your claims
	token.Claims = &JWTClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		data,
	}
	// Sign the token with your secret key
	return token.SignedString(jwtSecret)
}

func DecodeToken(token string) (*types.SessionUser, error) {
	claims, err := getClaimsFromToken(token)
	if err != nil {
		return nil, err
	}
	return claims.Data, nil
}

func getClaimsFromToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, custom_errors.InvalidJWTToken
	}
	claims := token.Claims.(*JWTClaims)
	return claims, nil
	// token, err := jwt.Parse(tokenString, func(tk *jwt.Token) (interface{}, error) {
	// 	if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		logger.Logger.Error("unexpected signing method:", tk.Header["alg"])
	// 		return nil, fmt.Errorf("unexpected signing method: %v", tk.Header["alg"])
	// 	}
	// 	return jwtSecret, nil
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(token.Claims)
	// if claims, ok := token.Claims.(JWTClaims); ok && token.Valid {
	// 	fmt.Println("Token is valid:", claims)
	// 	return claims, nil
	// } else {
	// 	fmt.Println("Token is invalid", ok, claims)
	// }
	// return nil, err
}
