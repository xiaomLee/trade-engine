package entrust

import "github.com/xiaomLee/trade-engine/entrust/queue"

type Cache interface {
	AddBuy(coinType string, item queue.Item) error
	AddSell(coinType string, item queue.Item) error
	RemoveBuy(coinType string, index int) error
	RemoveSell(coinType string, index int) error
	GetBuyList(coinType string) []queue.Item
	GetSellList(coinType string) []queue.Item
	AddFollower(string, string) error
	Close() error
}
