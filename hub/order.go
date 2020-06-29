package hub

const (
	BuyOrderSide  = 1
	SellOrderSide = 2
)

type Order struct {
	Side    int8
	Price   float64
	Num     int
	Uid     int64
	OrderId int64
}
