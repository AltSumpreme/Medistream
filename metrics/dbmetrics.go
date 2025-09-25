package metrics

import (
	"time"

	"gorm.io/gorm"
)

func DbMetrics(db *gorm.DB, operation string, fn func(*gorm.DB) error) error {
	start := time.Now()
	err := fn(db)
	DBLatency.WithLabelValues(operation).Observe(time.Since(start).Seconds())
	return err
}
