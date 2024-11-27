package main

import (
	"math/rand"
	"order-processing/internal/models"
	"sync"
	"time"
)

func (app *application) watch() <-chan *models.Order {
	pendingStream := make(chan *models.Order, 500)

	fetch := func() {
		pending, err := app.store.Orders.GetCreatedOrders(app.ctx)
		if err != nil {
			app.logger.Errorf("Error fetching pending orders: %v", err)
			return
		}

		app.logger.Infof("Orders watch query finished: %d new orders", len(pending))

		for _, order := range pending {
			select {
			case <-app.ctx.Done():
				return
			case pendingStream <- order:
			}
		}
	}

	fetch()

	go func() {
		defer close(pendingStream)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-app.ctx.Done():
				return
			case <-ticker.C:
				fetch()
			}
		}
	}()

	return pendingStream
}

func (app *application) managePending(pending map[int]*models.Order, watchStream <-chan *models.Order) <-chan *models.Order {
	pendingStream := make(chan *models.Order)

	go func() {
		m := &sync.Mutex{}

		for {
			select {
			case <-app.ctx.Done():
				return
			case order := <-watchStream:
				m.Lock()
				if _, exists := pending[order.ID]; exists {
					m.Unlock()
					continue
				}

				err := app.store.Orders.ChangeOrderStatus(app.ctx, order.ID, "pending")
				if err != nil {
					m.Unlock()
					app.logger.Warnf("Error while setting order %d as pending: %v", order.ID, err)
					continue
				}

				pending[order.ID] = order
				m.Unlock()

				pendingStream <- order
			}
		}
	}()

	return pendingStream
}

func (app *application) orderSimulate() {
	go func() {
		ticker := time.NewTicker(time.Duration(rand.Intn(5)+2) * time.Second)
		defer ticker.Stop()

		simulate := func() {
			amount := rand.Intn(300) + 1

			app.logger.Infof("Adding %v new simulated orders", amount)
			for i := 0; i <= amount; i++ {
				order := &models.Order{
					Status: "created",
				}
				app.store.Orders.CreateOrder(app.ctx, order)
			}
		}

		simulate()

		for {
			select {
			case <-app.ctx.Done():
				return
			case <-ticker.C:
				simulate()
			}
		}
	}()
}

func (app *application) run() {
	app.orderSimulate()
	pendingOrders := make(map[int]*models.Order)

	watchStream := app.watch()
	pendingStream := app.managePending(pendingOrders, watchStream)

	go func() {
		for {
			select {
			case <-app.ctx.Done():
				return
			case order := <-pendingStream:
				app.logger.Info("Pending order added:", order)
			}
		}
	}()
}
