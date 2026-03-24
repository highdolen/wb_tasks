package database

import (
	"context"
	"errors"

	"L0_optimize/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder — создание заказа с транзакцией
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Вставка Delivery
	err = tx.QueryRow(ctx, `
		INSERT INTO delivery (name, phone, zip, city, address, region, email)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id
	`, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	).Scan(&order.Delivery.ID)
	if err != nil {
		return err
	}

	// Вставка Payment
	err = tx.QueryRow(ctx, `
		INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id
	`, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
	).Scan(&order.Payment.ID)
	if err != nil {
		return err
	}

	// Вставка Order
	_, err = tx.Exec(ctx, `
		INSERT INTO orders (
			order_uid, track_number, entry, delivery_id, payment_id, locale,
			internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Delivery.ID, order.Payment.ID,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	// Вставка Items
	for _, item := range order.Items {
		_, err := tx.Exec(ctx, `
			INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		`, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status, order.OrderUID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetOrderByUID — получение заказа по UID
func (r *OrderRepository) GetOrderByUID(ctx context.Context, uid string) (*models.Order, error) {
	var order models.Order

	err := r.db.QueryRow(ctx, `
		SELECT
			o.order_uid,
			o.track_number,
			o.entry,
			o.locale,
			o.internal_signature,
			o.customer_id,
			o.delivery_service,
			o.shardkey,
			o.sm_id,
			o.date_created,
			o.oof_shard,
			d.id,
			d.name,
			d.phone,
			d.zip,
			d.city,
			d.address,
			d.region,
			d.email,
			p.id,
			p.transaction,
			p.request_id,
			p.currency,
			p.provider,
			p.amount,
			p.payment_dt,
			p.bank,
			p.delivery_cost,
			p.goods_total,
			p.custom_fee
		FROM orders o
		JOIN delivery d ON o.delivery_id = d.id
		JOIN payment p ON o.payment_id = p.id
		WHERE o.order_uid = $1
	`, uid).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
		&order.Delivery.ID,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
		&order.Payment.ID,
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order.Items = make([]models.Item, 0, 4)

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &order, nil
}

// GetAllOrders — получить все заказы
func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := r.db.Query(ctx, `
		SELECT order_uid FROM orders
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}

		order, err := r.GetOrderByUID(ctx, uid)
		if err != nil {
			return nil, err
		}
		if order != nil {
			orders = append(orders, *order)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
