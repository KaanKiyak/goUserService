package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	_ "github.com/mileusna/useragent"
	"time"
	"user-service/pkg/config"
)

// bu kod hashı ayrı sessionlıyacağımız zaman kullanılır
type SessionData struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Browser   string `json:"browser"`
	OS        string `json:"os"`
	Device    string `json:"device"`
	LoginTime string `json:"login_time"`
}

func TrackLogin(userID int, ip string, userAgent string) error {
	ua := useragent.Parse(userAgent)

	sessionInfo := map[string]string{
		"IP":        ip,
		"UserAgent": userAgent,
		"Browser":   ua.Name + " " + ua.Version,
		"OS":        ua.OS,
		"Device":    ua.Device,
		"LoginTime": time.Now().Format(time.RFC3339),
	}
	data, err := json.Marshal(sessionInfo)
	if err != nil {
		return err
	}
	redisKey := fmt.Sprintf("user:%d:session", userID)

	sessionID := uuid.New().String()

	if err := config.Rdb.HSet(context.Background(), redisKey, sessionID, data).Err(); err != nil {
		return err
	}

	return config.Rdb.Expire(config.Ctx, redisKey, 30*24*time.Hour).Err()
}
