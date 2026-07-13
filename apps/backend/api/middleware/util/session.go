package middleware_util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
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

func ParseJWT(public_id uuid.UUID, token_string string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (any, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("invalid JWT claims")
		}

		user_public_id, ok := claims["user_public_id"].(string)
		if !ok {
			return nil, errors.New("invalid session: missing user_public_id claim")
		}

		ctx := context.Background()
		jwt_data, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key(user_public_id+":jwt_data").Build(),
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

		return []byte(jwtKey.Key), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("session expired")
	}

	return token.Claims.(jwt.MapClaims), nil
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
