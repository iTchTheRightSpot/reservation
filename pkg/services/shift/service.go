package shift

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/shift"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IShiftService interface {
	Create(ctx context.Context, dto *shift.ShiftDto) error
}

type shiftService struct {
	logger   utils.ILogger
	adapters *stores.Adapters
}

func NewShiftService(l utils.ILogger, a *stores.Adapters) IShiftService {
	return &shiftService{logger: l, adapters: a}
}

func (dep *shiftService) Create(ctx context.Context, dto *shift.ShiftDto) error {
	_, err := dto.CheckForOverlappingSegments(dep.logger.Timezone())
	if err != nil {
		return err
	}

	_, err = dep.adapters.StaffStore.StaffByUUID(ctx, dto.StaffUUID)
	if err != nil {
		return err
	}

	// asynchronously validate no duplicate shift already saved

	return nil
}
