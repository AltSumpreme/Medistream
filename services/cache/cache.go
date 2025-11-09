package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Rdb *redis.Client
	Ctx context.Context
}

func NewCache(rdb *redis.Client, ctx context.Context) *Cache {
	return &Cache{
		Rdb: rdb,
		Ctx: ctx,
	}
}

func (c *Cache) AppointmentInvalidate(appointmentID, doctorID, patientID, date string) {
	keys := []string{
		fmt.Sprintf("cache:appointment:%s", appointmentID),
		fmt.Sprintf("cache:appointments:doctor:%s*", doctorID),
		fmt.Sprintf("cache:appointments:patient:%s*", patientID),
		fmt.Sprintf("cache:doctorSchedule:%s:%s", doctorID, date),
	}

	for _, key := range keys {
		iter := c.Rdb.Scan(c.Ctx, 0, key, 0).Iterator()
		for iter.Next(c.Ctx) {
			c.Rdb.Del(c.Ctx, iter.Val())
		}
	}
}

func (c *Cache) MedicalRecordInvalidate(medicalRecordID, patientID string) {
	keys := []string{
		fmt.Sprintf("cache:medicalRecord:%s", medicalRecordID),
		fmt.Sprintf("cache:medicalRecords:patient:%s*", patientID),
	}

	for _, key := range keys {
		iter := c.Rdb.Scan(c.Ctx, 0, key, 0).Iterator()
		for iter.Next(c.Ctx) {
			c.Rdb.Del(c.Ctx, iter.Val())
		}
	}
}

func (c *Cache) PrescriptionInvalidate(patientID string) {
	keys := []string{
		fmt.Sprintf("cache:prescriptions:patient:%s*", patientID),
	}

	for _, key := range keys {
		iter := c.Rdb.Scan(c.Ctx, 0, key, 0).Iterator()
		for iter.Next(c.Ctx) {
			c.Rdb.Del(c.Ctx, iter.Val())
		}
	}
}

func (c *Cache) ReportInvalidate(patientID string) {
	keys := []string{
		fmt.Sprintf("cache:reports:patient:%s*", patientID),
	}

	for _, key := range keys {
		iter := c.Rdb.Scan(c.Ctx, 0, key, 0).Iterator()
		for iter.Next(c.Ctx) {
			c.Rdb.Del(c.Ctx, iter.Val())
		}
	}
}

func (c *Cache) VitalsInvalidate(patientID string) {
	keys := []string{
		fmt.Sprintf("cache:vitals:patient:%s*", patientID),
	}

	for _, key := range keys {
		iter := c.Rdb.Scan(c.Ctx, 0, key, 0).Iterator()
		for iter.Next(c.Ctx) {
			c.Rdb.Del(c.Ctx, iter.Val())
		}
	}
}

func SaveOTP(email, otp string, ttl time.Duration) error {
	return config.Rdb.Set(config.Ctx, fmt.Sprintf("otp:%s", email), otp, ttl).Err()
}

func VerifyOTP(email, otp string) error {
	key := fmt.Sprintf("otp:%s", email)

	storedOTP, err := config.Rdb.Get(config.Ctx, key).Result()
	if err != nil {
		return fmt.Errorf("OTP not found or expired")
	}

	if storedOTP != otp {
		return fmt.Errorf("invalid otp")
	}

	config.Rdb.Del(config.Ctx, key)

	return nil
}
