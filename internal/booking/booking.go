package booking

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNotFound  = errors.New("booking: not found")
	ErrConflict  = errors.New("booking: slot conflict")
	ErrInvalid   = errors.New("booking: invalid argument")
	ErrCancelled = errors.New("booking: already cancelled")
)

type Room struct {
	ID   string
	Name string
}

type Reservation struct {
	ID        string
	RoomID    string
	Start     int // inclusive
	End       int // exclusive
	Cancelled bool
}

type Service struct {
	rooms        map[string]Room
	reservations map[string]Reservation
	byRoom       map[string][]string
	nextID       int
}

func NewService() *Service {
	return &Service{
		rooms:        make(map[string]Room),
		reservations: make(map[string]Reservation),
		byRoom:       make(map[string][]string),
	}
}

func (s *Service) AddRoom(id, name string) error {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" {
		return fmt.Errorf("add room: id and name must be non-empty: %w", ErrInvalid)
	}
	if _, exists := s.rooms[id]; exists {
		return fmt.Errorf("add room %q: duplicate id: %w", id, ErrInvalid)
	}
	s.rooms[id] = Room{
		ID:   id,
		Name: name,
	}
	return nil

}

func (s *Service) Book(roomID string, start, end int) (Reservation, error) {
	if _, ok := s.rooms[roomID]; !ok {
		return Reservation{}, fmt.Errorf("book %s: room: %w", roomID, ErrNotFound)

	}
	if start < 0 || start >= end || end > 24 {
		return Reservation{}, fmt.Errorf("book %s: slot [%d,%d): %w", roomID, start, end, ErrInvalid)
	}

	for _, id := range s.byRoom[roomID] {
		r := s.reservations[id]
		if r.Cancelled {
			continue
		}
		if start < r.End && r.Start < end {
			return Reservation{}, fmt.Errorf("book %s: [%d,%d) overlaps %s: %w", roomID, start, end, r.ID, ErrConflict)
		}
	}
	s.nextID++
	res := Reservation{
		ID:     strconv.Itoa(s.nextID),
		RoomID: roomID,
		Start:  start,
		End:    end,
	}
	s.reservations[res.ID] = res
	s.byRoom[roomID] = append(s.byRoom[roomID], res.ID)
	return res, nil
}

func (s *Service) Cancel(reservationID string) error {
	res, ok := s.reservations[reservationID]
	if !ok {
		return fmt.Errorf("cancel %s: %w", reservationID, ErrNotFound)
	}
	if res.Cancelled {
		return fmt.Errorf("cancel %s: %w", reservationID, ErrCancelled)
	}
	res.Cancelled = true
	s.reservations[reservationID] = res
	return nil
}

func (s *Service) Get(reservationID string) (Reservation, error) {
	res, ok := s.reservations[reservationID]
	if !ok {
		return Reservation{}, fmt.Errorf("get %s: %w", reservationID, ErrNotFound)
	}
	return res, nil
}
func (s *Service) ListRoom(roomID string) ([]Reservation, error) {
	if _, ok := s.rooms[roomID]; !ok {
		return nil, fmt.Errorf("list %s: room: %w", roomID, ErrNotFound)
	}
	ids := s.byRoom[roomID]
	out := make([]Reservation, 0, len(ids))

	for _, id := range ids {
		out = append(out, s.reservations[id])
	}
	return out, nil

} // includes cancelled
