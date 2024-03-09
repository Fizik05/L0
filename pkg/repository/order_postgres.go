package repository

import (
	"fmt"

	"github.com/Fizik05/L0"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) SaveOrder(order L0.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	insertOrderQuery := fmt.Sprintf(
		"INSERT INTO %s (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		ordersTable,
	)
	_, err = tx.Exec(insertOrderQuery, order.Order_uid, order.Track_number, order.Entry, order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.ShardKey, order.SM_id, order.Date_created, order.OOF_shard)
	if err != nil {
		logrus.Errorf("Error while insert into ordersTable: %s", err.Error())
		tx.Rollback()
		return err
	}

	insertDeliveryQuery := fmt.Sprintf(
		"INSERT INTO %s (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		deliveriesTable,
	)
	_, err = tx.Exec(insertDeliveryQuery, order.Order_uid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		logrus.Errorf("Error while insert into deliveryTable: %s", err.Error())
		tx.Rollback()
		return err
	}

	insertPaymnetQuery := fmt.Sprintf(
		"INSERT INTO %s (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		paymentsTable,
	)
	_, err = tx.Exec(insertPaymnetQuery, order.Payment.Transaction, order.Payment.Request_id, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.Payment_dt, order.Payment.Bank, order.Payment.Delivery_cost, order.Payment.Goods_total, order.Payment.Custom_fee)
	if err != nil {
		logrus.Errorf("Error while insert into paymentTable: %s", err.Error())
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		insertItemQuery := fmt.Sprintf(
			"INSERT INTO %s (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			itemsTable,
		)
		_, err = tx.Exec(insertItemQuery, order.Order_uid, item.CHRT_id, item.Track_number, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.Total_price, item.NM_id, item.Brand, item.Status)
		if err != nil {
			logrus.Errorf("Error while insert into itemTable: %s", err.Error())
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderPostgres) RecoverCache() ([]L0.Order, error) {
	var delivery L0.Delivery
	var payment L0.Payment
	var items []L0.Item
	var orders []L0.Order

	query := fmt.Sprintf(
		"SELECT * FROM %s",
		ordersTable,
	)
	err := r.db.Select(&orders, query)
	if err != nil {
		logrus.Errorf("Problem with connection to Order table")
		return nil, err
	}

	for ind, order := range orders {
		UID := order.Order_uid
		deliveryQuery := fmt.Sprintf(
			"SELECT name, phone, zip, city, address, region, email FROM %s WHERE order_uid = $1",
			deliveriesTable,
		)
		err = r.db.Get(&delivery, deliveryQuery, UID)
		if err != nil {
			logrus.Errorf("Problem with connection to Delivery table")
			return nil, err
		}
		orders[ind].Delivery = delivery

		paymentQuery := fmt.Sprintf(
			"SELECT * FROM %s WHERE transaction = $1",
			paymentsTable,
		)
		err = r.db.Get(&payment, paymentQuery, UID)
		if err != nil {
			logrus.Errorf("Problem with connection to Payment table")
			return nil, err
		}
		orders[ind].Payment = payment

		itemsQuery := fmt.Sprintf(
			"SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM %s WHERE order_uid = $1",
			itemsTable,
		)
		err = r.db.Select(&items, itemsQuery, UID)
		if err != nil {
			logrus.Errorf("Problem with connection to Item table")
			return nil, err
		}

		orders[ind].Items = make([]L0.Item, len(items))
		copy(orders[ind].Items, items)
	}

	return orders, nil
}
