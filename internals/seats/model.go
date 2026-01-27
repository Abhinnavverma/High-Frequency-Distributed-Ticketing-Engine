package seats

type Seat struct {
	ID         int32  `json:"id"`
	RowNumber  string `json:"row_number"`
	SeatNumber int32  `json:"seat_number"`
	Status     string `json:"status"`
	Price      int32  `json:"price"`
}

type SeatCreationRequest struct {
	RowNumber  string `json:"row_number"`
	SeatNumber int32  `json:"seat_number"`
	Price      int32  `json:"price"`
}

type SeatCreationResponse struct {
	Message string `json:"message"`
}
