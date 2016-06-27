package user

import (
	"fmt"
	"projects/resellbay/server/cfg"
	"projects/resellbay/server/util"
	"projects/sample-tutorials/server/pkg/crypt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
)

//Login logs user into the app
func (u *User) Login(input *LoginInput) *LoginOutput {
	item := u.GetByEmail(input.Useremail)
	if *item.Count > 0 && crypt.ComparePasswords(input.Password, *item.Items[0]["password"].S) {
		token := createToken(*item.Items[0]["userID"].S, input)
		return &LoginOutput{true, token, *item.Items[0]["userID"].S, "Login successful"}
	}
	return &LoginOutput{Valid: false, Msg: "Login failed: Wrong Email or Password"}
}

func signupValid(input *AddInput) *AddOutput {
	if govalidator.IsEmail(input.Useremail) {
		if len(input.Password) > 8 {
			return &AddOutput{true, "", ""}
		}
		return &AddOutput{false, "Password invalid", "Password"}
	}
	return &AddOutput{false, "Email Invalid", "Email"}
}

func createToken(userID string, input *LoginInput) string {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["userID"] = userID
	token.Claims["useremail"] = input.Useremail
	token.Claims["password"] = input.Password
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(cfg.SigningKey)
	util.CheckErr(err)
	return tokenString
}

func verifyToken(myToken string, userID string) bool {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		if token.Claims["userID"] == userID {
			return cfg.SigningKey, nil
		}

		return "fail", nil

	})
	util.CheckErr(err)

	if err == nil && token.Valid {
		return true
	}
	return false

}
