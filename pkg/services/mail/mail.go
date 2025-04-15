package mail

import (
	"github.com/iTchTheRightSpot/reservation/config"
	"github.com/iTchTheRightSpot/utility/utils"
)

type IMailService interface {
	SendReservationConfirmation() error
}

type mailService struct {
	logger utils.ILogger
	env    *config.SecretVariables
}

func NewMailService(l utils.ILogger, e *config.SecretVariables) IMailService {
	return &mailService{logger: l, env: e}
}

func (dep *mailService) SendReservationConfirmation() error {
	return nil
}