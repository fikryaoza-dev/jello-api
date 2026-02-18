package usecase

import "jello-api/internal/repository"

type OrderUsecase struct {
	Repo repository.IBookingRepository
}

func NewBookingUsecase(repo repository.IBookingRepository) *BookingUsecase {
	return &BookingUsecase{
		Repo: repo,
	}
}

// func generateOrderNumber() string {
//     now := time.Now()
//     // Format: ORD-YYYYMMDD-Last4CharsOfUUID
//     datePart := now.Format("20060102")
//     uniquePart := strings.ToUpper(uuid.New().String()[:4])
//     return fmt.Sprintf("ORD-%s-%s", datePart, uniquePart)
// }