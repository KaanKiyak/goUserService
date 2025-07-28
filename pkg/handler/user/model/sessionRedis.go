package model

import (
	"fmt"
	"github.com/mileusna/useragent"
	_ "github.com/mileusna/useragent"
	"time"
	"user-service/pkg/config"
)

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
	timesTime := time.Now().Unix()

	hashKey := fmt.Sprintf("user:%d session:%d", userID, timesTime)

	sessionInfo := map[string]string{
		"IP":        ip,
		"UserAgent": userAgent,
		"Browser":   ua.Name + " " + ua.Version,
		"OS":        ua.OS,
		"Device":    ua.Device,
		"LoginTime": time.Now().Format(time.RFC3339),
	}
	if err := config.Rdb.HSet(config.Ctx, hashKey, sessionInfo).Err(); err != nil {
		return err
	}
	config.Rdb.Expire(config.Ctx, hashKey, 30*24*time.Hour)
	redisKey := fmt.Sprintf("user:%d:session", userID)
	err := config.Rdb.LPush(config.Ctx, redisKey, hashKey).Err()
	if err != nil {
		return err
	}
	config.Rdb.LTrim(config.Ctx, redisKey, 0, 9)
	config.Rdb.Expire(config.Ctx, redisKey, 30*24*time.Hour)

	return nil
}
