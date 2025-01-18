package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"time"
)

func main() {
	secret := config.SecretVariables{}
	env := secret.Config()

	db, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	l, err := utils.NewLogger("UTC")
	if err != nil {
		return
	}

	a := stores.NewAdapters(l, db, stores.NewTransactionProvider(l, db))
	ctx := context.Background()

	// 1. Create profiles, roles, permissions & staffs
	stafs := staffs(ctx, a)

	// 2. Create services.
	s := services(ctx, a)

	// 3. Assign to services.
	assignServices(ctx, a, &s, &stafs)

	// 4. Create schedules for said staffs.
	schedules(ctx, l.Date(), a, &stafs)
}

func services(ctx context.Context, a *stores.Adapters) [10]service.ServiceTypeEntity {
	arr := [10]service.ServiceTypeEntity{}
	for i := 0; i < 10; i++ {
		s := service.ServiceTypeEntity{
			Name:        fmt.Sprintf("Service-%v-%s", i, uuid.New().String()),
			Price:       float64(20 * (i + 1)),
			IsVisible:   true,
			Duration:    3600,
			CleanUpTime: 30 * 60,
		}
		if err := a.ServiceStore.Save(ctx, &s); err != nil {
			log.Fatal(err.Error())
		}
		arr[i] = s
	}
	return arr
}

func staffs(ctx context.Context, a *stores.Adapters) [6]staff.Staff {
	arr := [6]staff.Staff{}

	for i := 0; i < 6; i++ {
		p := profile.Profile{
			Firstname: fmt.Sprintf("Firstname-%v", i),
			Lastname:  fmt.Sprintf("Lastname-%v", i),
			Email:     fmt.Sprintf("firstname-%v-%s@email.com", i, uuid.New().String()),
		}

		if err := a.ProfileStore.Save(ctx, &p); err != nil {
			log.Fatal(err.Error())
		}

		r := models.Role{
			Role:      models.STAFF,
			ProfileId: p.ProfileId,
		}

		if err := a.RoleStore.Save(ctx, &r); err != nil {
			log.Fatal(err.Error())
		}

		if err := a.PermissionStore.Save(ctx, &models.Permission{
			Permission: models.READ,
			RoleId:     r.RoleId,
		}); err != nil {
			log.Fatal(err.Error())
		}

		if i%2 == 0 {
			if err := a.PermissionStore.Save(ctx, &models.Permission{
				Permission: models.WRITE,
				RoleId:     r.RoleId,
			}); err != nil {
				log.Fatal(err.Error())
			}
		}

		lorem := "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Amet corporis, deserunt eaque earum molestias neque nesciunt nulla placeat sed suscipit."
		st := staff.Staff{
			UUID:      uuid.New(),
			Bio:       &lorem,
			ProfileId: &p.ProfileId,
		}

		if err := a.StaffStore.Save(ctx, &st); err != nil {
			log.Fatal(err.Error())
		}

		arr[i] = st
	}

	return arr
}

func assignServices(ctx context.Context, a *stores.Adapters, s *[10]service.ServiceTypeEntity, st *[6]staff.Staff) {
	windowSize := 3
	idx := 0

	for _, staf := range st {
		for i := 0; i < windowSize; i++ {
			currentIndex := (idx + i) % len(s)
			if err := a.StaffServiceStore.Save(ctx, &staff.StaffServiceEntity{
				StaffId:   staf.StaffId,
				ServiceId: s[currentIndex].ServiceId,
			}); err != nil {
				log.Fatal(err.Error())
			}
		}
		idx = (idx + windowSize) % len(s)
	}
}

func schedules(ctx context.Context, date time.Time, a *stores.Adapters, st *[6]staff.Staff) {
	for _, staf := range st {
		for j := 0; j < len(st); j++ {
			date = date.Add(time.Duration(48*(j+1)) * time.Hour)
			d := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location())
			if err := a.ScheduleStore.Save(ctx, &schedule.Schedule{
				StaffId:   staf.StaffId,
				Start:     d,
				End:       d.Add(8 * time.Hour),
				IsVisible: true,
			}); err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}
