package seats

type Seat struct {
	ID         int32  `json:"id"`
	RowNumber  string `json:"row_number"`
	SeatNumber int32  `json:"seat_number"`
	Status     string `json:"status"`
	Price      int32  `json:"price"`
}
