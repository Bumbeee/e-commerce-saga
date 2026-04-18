package order

func CanTransitionTo(from, to Status) bool {
	switch from {
	case StatusCreated:
		return to == StatusPaid || to == StatusCancelled
	case StatusPaid:
		return to == StatusShipped || to == StatusCancelled
	case StatusShipped:
		return to == StatusDelivered || to == StatusCancelled
	case StatusDelivered, StatusCancelled:
		return false
	default:
		return false
	}
}
