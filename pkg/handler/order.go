package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getOrder(c *gin.Context) {
	order_uid := c.Param("order_uid")

	item, err := h.services.Cache.GetOrder(order_uid)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(
		http.StatusOK, "index.html", gin.H{
			"Order_uid":          item.Order_uid,
			"Track_number":       item.Track_number,
			"Entry":              item.Entry,
			"Delivery":           item.Delivery,
			"Payment":            item.Payment,
			"Items":              item.Items,
			"Locale":             item.Locale,
			"Internal_signature": item.Internal_signature,
			"Customer_id":        item.Customer_id,
			"Delivery_service":   item.Delivery_service,
			"ShardKey":           item.ShardKey,
			"SM_id":              item.SM_id,
			"Date_created":       item.Date_created.Format("2021-11-26T06:22:19Z"),
			"OOF_shard":          item.OOF_shard,
		},
	)
}
