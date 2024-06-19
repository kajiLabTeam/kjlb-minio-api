package helper

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

func AuthBasic(authHeader string) error {

	// "Basic " 接頭辞を削除して、Base64でエンコードされた文字列を取得
	authValue := strings.TrimPrefix(authHeader, "Basic ")

	// Base64デコード
	decoded, err := base64.StdEncoding.DecodeString(authValue)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return err
	}

	// デコードされた文字列を取得
	credentials := string(decoded)

	// ユーザー名とパスワードの分割
	split := strings.SplitN(credentials, ":", 2)
	user := split[0]
	pass := split[1]

	if user != os.Getenv("MINIO_ROOT_USER") {
		return errors.New("invalid username")
	}

	if pass != os.Getenv("MINIO_ROOT_PASSWORD") {
		return errors.New("invalid password")
	}
	return nil
}
