package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type AppointmentCache struct {
	Rdb *redis.Client
	Ctx context.Context
}

func NewAppointmentCache(rdb *redis.Client, ctx context.Context) *AppointmentCache {
	return &AppointmentCache{
		Rdb: rdb,
		Ctx: ctx,
	}
}

func (ac *AppointmentCache) Invalidate(appointmentID, doctorID, patientID string) {
	keys := []string{
		fmt.Sprintf("cache:appointment:%s", appointmentID),
		fmt.Sprintf("cache:appointments:doctor:%s*", doctorID),
		fmt.Sprintf("cache:appointments:patient:%s*", patientID),
		//	fmt.Sprintf("cache:doctorSchedule:%s:%s", doctorID, date),
	}

	for _, key := range keys {
		iter := ac.Rdb.Scan(ac.Ctx, 0, key, 0).Iterator()
		for iter.Next(ac.Ctx) {
			ac.Rdb.Del(ac.Ctx, iter.Val())
		}
	}
}
