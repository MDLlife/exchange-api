package exchange

import (
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

var store = &tracker{
	orders:    make(map[int]idx),
	opened:    make(map[int]Order),
	completed: make(map[int]Order),
}
var accepttime = time.Now()

func TestTrackerPush(t *testing.T) {
	var tests = []struct {
		order Order
		err   error
	}{
		{
			order: Order{
				OrderID:   1,
				Status:    Submitted,
				Market:    "LTC/BTC",
				Price:     decimal.NewFromFloat(123.456),
				Amount:    decimal.NewFromFloat(1.0),
				Type:      Buy,
				Submitted: accepttime,
				Accepted:  accepttime,
			},
			err: nil,
		},
		{
			order: Order{
				OrderID:   2,
				Status:    Submitted,
				Market:    "SKY/DOGE",
				Price:     decimal.NewFromFloat(654.321),
				Amount:    decimal.NewFromFloat(5.0),
				Type:      Sell,
				Submitted: accepttime,
				Accepted:  accepttime,
			},
			err: nil,
		},
		{
			order: Order{
				OrderID:   1,
				Status:    Opened,
				Price:     decimal.NewFromFloat(123.456),
				Amount:    decimal.NewFromFloat(10.0),
				Submitted: accepttime,
				Market:    "BTC/LTC",
				Type:      Sell,
				Accepted:  accepttime,
			},
			err: ErrExist,
		},
		{
			order: Order{
				OrderID:   5,
				Price:     decimal.NewFromFloat(123.456),
				Amount:    decimal.NewFromFloat(10.0),
				Submitted: accepttime,
				Market:    "BTC/LTC",
				Type:      Sell,
				Accepted:  accepttime,
			},
			err: ErrInvalidStatus,
		}, {
			order: Order{
				OrderID:   8,
				Status:    Submitted,
				Market:    "LTC/BTC",
				Price:     decimal.NewFromFloat(123.456),
				Amount:    decimal.NewFromFloat(20.0),
				Type:      Buy,
				Submitted: accepttime,
				Accepted:  accepttime,
			},
			err: nil,
		},
	}
	for _, v := range tests {
		err := store.Push(v.order)
		if err != v.err {
			t.Errorf("want error %v, expected %v", err, v.err)
		}
	}
}
func TestTrackerUpdate(t *testing.T) {
	var testtime = time.Now().Add(time.Minute * 1)
	var tests = []struct {
		upd     Order
		orderid int
		result  Order
		err     error
	}{
		{
			orderid: 1,
			err:     nil,
			result: Order{
				OrderID:         1,
				Status:          Completed,
				Market:          "LTC/BTC",
				Price:           decimal.NewFromFloat(123.456),
				Amount:          decimal.NewFromFloat(1.0),
				Type:            Buy,
				Submitted:       accepttime,
				Accepted:        accepttime,
				Completed:       testtime,
				Fee:             decimal.NewFromFloat(0.01),
				CompletedAmount: decimal.NewFromFloat(1.0),
			},
			upd: Order{
				OrderID:         1,
				Market:          "LTC/BTC",
				Price:           decimal.NewFromFloat(123.456),
				Amount:          decimal.NewFromFloat(1.0),
				Type:            Buy,
				Submitted:       accepttime,
				Status:          Completed,
				Accepted:        accepttime,
				Completed:       testtime,
				Fee:             decimal.NewFromFloat(0.01),
				CompletedAmount: decimal.NewFromFloat(1.0),
			},
		},
		{
			orderid: 2,
			err:     nil,
			result: Order{
				OrderID:         2,
				Market:          "SKY/DOGE",
				Price:           decimal.NewFromFloat(654.321),
				Amount:          decimal.NewFromFloat(5.0),
				Type:            Sell,
				Submitted:       accepttime,
				Status:          Cancelled,
				CompletedAmount: decimal.NewFromFloat(4.0),
				Fee:             decimal.NewFromFloat(0.04),
				Accepted:        accepttime,
				Completed:       testtime,
			},
			upd: Order{
				OrderID:         2,
				Market:          "SKY/DOGE",
				Price:           decimal.NewFromFloat(654.321),
				Amount:          decimal.NewFromFloat(5.0),
				Type:            Sell,
				Submitted:       accepttime,
				Status:          Cancelled,
				CompletedAmount: decimal.NewFromFloat(4.0),
				Fee:             decimal.NewFromFloat(0.04),
				Accepted:        accepttime,
				Completed:       testtime,
			},
		},
		{
			orderid: 8,
			result: Order{
				OrderID:   8,
				Status:    Completed,
				Market:    "LTC/BTC",
				Price:     decimal.NewFromFloat(123.456),
				Amount:    decimal.NewFromFloat(20.0),
				Type:      Buy,
				Submitted: accepttime,
				Completed: testtime,
				Accepted:  accepttime,
			},
			err: nil,
			upd: Order{
				OrderID:   123,
				Status:    Completed,
				Market:    "LTC/BTC",
				Price:     decimal.NewFromFloat(123.456),
				Amount:    decimal.NewFromFloat(20.0),
				Type:      Buy,
				Submitted: accepttime,
				Completed: testtime,
				Accepted:  accepttime,
			},
		},
	}
	for _, v := range tests {
		err := store.UpdateOrder(v.upd)
		if err != v.err {
			t.Errorf("want error %v, expected %v", v.err, err)
		}
		order, err := store.GetOrderInfo(v.orderid)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(order, v.result) {
			t.Errorf("want %#v, expected %#v", v.result, order)
		}
	}
}
