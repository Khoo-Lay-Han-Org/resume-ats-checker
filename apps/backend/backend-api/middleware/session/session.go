package middleware_session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	valkey "github.com/valkey-io/valkey-go"
	"resuming/tool"
)

func ExtractSessionCookie(cookie string) (uuid.UUID, string, error) {
	var session_data map[string]string
	err := json.Unmarshal([]byte(cookie), &session_data)
	if err != nil {
		return uuid.Nil, "", errors.New("failed to process session")
	}

	public_id, err := uuid.Parse(session_data["public_id"])
	if err != nil {
		return uuid.Nil, "", errors.New("failed to process session")
	}

	token_string := session_data["token"]

	return public_id, token_string, nil
}

func ParseJWT(public_id uuid.UUID, token_string string) (map[string]any, error) {
	ctx := context.Background()
	jwt_data, err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Get().Key(public_id.String()+":jwt_data").Build(),
	).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, errors.New("session key not found")
		}
		return nil, errors.New("failed to read session key")
	}

	var jwtKey struct {
		Key string `json:"Key"`
	}
	if err := json.Unmarshal([]byte(jwt_data), &jwtKey); err != nil {
		return nil, errors.New("corrupted session key data")
	}

	object, err := jose.ParseEncrypted(token_string, []jose.KeyAlgorithm{jose.DIRECT}, []jose.ContentEncryption{jose.A128GCM})
	if err != nil {
		return nil, errors.New("invalid token format")
	}

	decoded, err := object.Decrypt([]byte(jwtKey.Key))
	if err != nil {
		return nil, errors.New("failed to decrypt token")
	}

	var claims map[string]any
	json.Unmarshal(decoded, &claims)

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func CheckSession(cookie string) (string, error) {
	user_uuid, token_string, e := ExtractSessionCookie(cookie)
	if e != nil {
		return "", e
	}

	parsed_token, e := ParseJWT(user_uuid, token_string)
	if e != nil {
		return "", e
	}

	jwt_user_id, ok := parsed_token["user_public_id"].(string)
	if !ok {
		return "", errors.New("invalid session: missing user_public_id claim")
	}

	ctx := context.Background()
	_, err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Get().
			Key(jwt_user_id+":session_data").
			Build(),
	).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "", errors.New("session expired or not found in store")
		}
		return "", fmt.Errorf("failed to read session from Valkey: %w", err)
	}

	return jwt_user_id, nil
}
