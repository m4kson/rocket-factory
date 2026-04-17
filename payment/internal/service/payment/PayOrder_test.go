package payment

//func (s *ServiceSuite) TestPayOrderSuccess() {
//	_, err := s.service.PayOrder(s.ctx)
//
//	s.NoError(err)
//}

func (s *ServiceSuite) TestPayOrder() {
	s.Run("PayOrderSuccess", func() {
		_, err := s.service.PayOrder(s.ctx)

		s.Require().NoError(err)
	})
}
