package mail

type MockMailService struct {
	SendReservationConfirmationCalled bool
	SendReservationConfirmationError  error
}

func (dep *MockMailService) SendReservationConfirmation() error {
	dep.SendReservationConfirmationCalled = true
	return dep.SendReservationConfirmationError
}
