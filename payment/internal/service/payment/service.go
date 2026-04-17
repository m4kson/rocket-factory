package payment

type service struct{}

func NewPaymentService() *service {
	return &service{}
}
