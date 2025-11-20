package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"ordersystem/httptools"
	"ordersystem/model"
	"ordersystem/repository"
	"ordersystem/storage"
	"strings"

	"github.com/go-chi/render"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// GetMenu 			godoc
// @tags 			Menu
// @Description 	Returns the menu of all drinks
// @Produce  		json
// @Success 		200 {array} model.Drink
// @Failure     	500
// @Router 			/api/menu [get]
func GetMenu(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allDrinks, err := db.GetDrinks()
		if err != nil {
			slog.Error("Unable to load drinks", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load drinks")
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, allDrinks)
	}
}

// GetOrders		godoc
// @tags 			Order
// @Description 	Returns all orders
// @Produce  		json
// @Success 		200 {array} model.Order
// @Failure     	500
// @Router 			/api/order/all [get]
func GetOrders(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allOrders, err := db.GetOrders()
		if err != nil {
			slog.Error("Unable to load orders", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load order")
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, allOrders)
	}
}

// GetOrdersTotal		godoc
// @tags 				Order
// @Description 		Gets totalled orders
// @Produce  			json
// @Success 			200
// @Failure     		500
// @Router 				/api/order/totalled [get]
func GetOrdersTotal(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalledOrders, err := db.GetTotalledOrders()
		if err != nil {
			slog.Error("Unable to load order totals", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load order totals")
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, totalledOrders)
	}
}

// GetReceiptFile		godoc
// @tags 				Order
// @Description 		Get receipt for order
// @Produce 			text/markdown
// @Success 			200 {file} markdown file
// @Param 				orderId path int true "Order ID"
// @Failure     		404
// @Failure     		500
// @Router 				/api/receipt/{orderId} [get]
func GetReceiptFile(db *repository.DatabaseHandler, s3 *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uintId, err := httptools.ParseUintUrlParam("orderId", r)
		if err != nil {
			slog.Error("No order id set")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "No order id set")
			return
		}
		// get order from db
		order, err := db.GetOrder(uintId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, "This order does not exist")
				return
			}
			slog.Error("Unable to load order", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load order")
			return
		}

		// read from s3
		obj, err := s3.GetObject(context.Background(), storage.OrdersBucket, order.GetFilename(), minio.GetObjectOptions{})
		if err != nil {
			slog.Error("Unable to get file from S3", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to retrieve receipt")
			return
		}
		defer obj.Close()

		// serve file
		w.Header().Set("Content-Type", "text/markdown")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", order.GetFilename()))

		if _, err := io.Copy(w, obj); err != nil {
			slog.Error("Unable to write file to response", slog.String("error", err.Error()))
			// Headers are already written, so we can't really send a JSON error status here easily
		}
	}
}

// PostOrder 		godoc
// @tags 			Order
// @Description 	Adds an order to the db
// @Accept 			json
// @Param 			b body model.Order true "Order"
// @Produce  		json
// @Success 		200
// @Failure     	400
// @Failure     	500
// @Router 			/api/order [post]
func PostOrder(db *repository.DatabaseHandler, s3 *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order model.Order
		// read body
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Unable to read body", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Unable to read body")
			return
		}
		err = json.Unmarshal(payload, &order)
		if err != nil {
			slog.Error("Unable to decode body", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Unable to decode body")
			return
		}
		// store to db
		dbOrder, err := db.AddOrder(&order)
		if err != nil {
			slog.Error("Unable to add order to db", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to add order to db")
			return
		}

		// store to s3
		markdown := dbOrder.ToMarkdown()
		reader := strings.NewReader(markdown)

		_, err = s3.PutObject(context.Background(), storage.OrdersBucket, dbOrder.GetFilename(), reader, int64(len(markdown)), minio.PutObjectOptions{
			ContentType: "text/markdown",
		})
		if err != nil {
			slog.Error("Unable to upload receipt to S3", slog.String("error", err.Error()))
			// Note: We successfully saved to DB, so we might return OK but log the S3 failure,
			// or return 500. Returning 500 for consistency with the exercise instructions.
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Order saved, but failed to generate receipt")
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, "ok")
	}
}
