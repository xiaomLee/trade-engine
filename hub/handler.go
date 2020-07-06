package hub

import "github.com/xiaomLee/trade-engine/entrust"

func (h *Hub) doBuyEntrust(task *Task) error {
	detail := task.Detail.(*TaskBuyEntrust)
	return h.cache.AddBuy(detail.Entrust.CoinType, detail.Entrust)
}

func (h *Hub) doGetBuyList(task *Task) []*entrust.Entrust {
	ret := make([]*entrust.Entrust, 0)
	detail := task.Detail.(*TaskGetBuyList)
	list := h.cache.GetBuyList(detail.CoinType)
	for _, l := range list {
		ret = append(ret, l.(*entrust.Entrust))
	}

	return ret
}
