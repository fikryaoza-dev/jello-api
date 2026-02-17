package mapper

import (
	"jello-api/internal/domain"
	"jello-api/internal/model"
)

type BookingMapper struct{}

func (BookingMapper) ToDomain(m model.Booking) domain.Booking {
	return domain.Booking{
		ID:              m.ID,
		Customer:        m.Customer,
		TableID:         m.TableID,
		DurationMinutes: m.DurationMinutes,
		Status:          domain.BookingStatus(m.Status),
	}
}

func (BookingMapper) ToModel(d domain.Booking) model.Booking {
	return model.Booking{
		ID:              d.ID,
		Customer:        d.Customer,
		TableID:         d.TableID,
		DurationMinutes: d.DurationMinutes,
		Status:          string(d.Status),
	}
}
