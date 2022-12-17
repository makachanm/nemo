package build

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type TemplateTools struct {
	skconfig SkinConfig
}

func NewTemplateTools(config SkinConfig) TemplateTools {
	return TemplateTools{skconfig: config}
}

func (t *TemplateTools) GetTimeStamp(stamp TimeStamp) string {
	ttime := time.Date(stamp.Year, time.Month(stamp.Month), stamp.Day, stamp.Hour, stamp.Min, 0, 0, time.UTC)
	return ttime.Format(t.skconfig.DateType)
}

func (t *TemplateTools) GetTodayStamp() string {
	ttime := time.Now()
	return ttime.Format(t.skconfig.DateType)
}

func (t *TemplateTools) GetTagnameHash(name string) string {
	hashbyte := md5.Sum([]byte(name))
	nhash := hex.EncodeToString(hashbyte[:])

	return nhash[:8]
}
