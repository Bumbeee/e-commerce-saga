package order

// CanTransitionTo проверяет, разрешён ли переход из текущего статуса в целевой
// Матрица переходов:
// created   → paid, cancelled
// paid      → shipped, cancelled
// shipped   → delivered, cancelled
// delivered → (terminal)
// cancelled → (terminal)
func (s Status) CanTransitionTo(next Status) bool {
	allowed := map[Status][]Status{
		StatusCreated:   {StatusPaid, StatusCancelled},
		StatusPaid:      {StatusShipped, StatusCancelled},
		StatusShipped:   {StatusDelivered, StatusCancelled},
		StatusDelivered: {},
		StatusCancelled: {},
	}

	for _, allowedStatus := range allowed[s] {
		if allowedStatus == next {
			return true
		}
	}
	return false
}
