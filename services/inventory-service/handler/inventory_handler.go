package handler

import (
	"net/http"

	"github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, svc inventorypb.InventoryServiceServer) {
	inv := r.Group("/inventory")
	{
		inv.GET("/:id", func(c *gin.Context) {
			resp, err := svc.GetInventory(c.Request.Context(), &inventorypb.GetInventoryRequest{
				ProductId: c.Param("id"),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})

		inv.POST("/:id", func(c *gin.Context) {
			var body struct {
				Change int32 `json:"change"`
			}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
				return
			}

			resp, err := svc.UpdateInventory(c.Request.Context(), &inventorypb.UpdateInventoryRequest{
				ProductId:      c.Param("id"),
				QuantityChange: body.Change,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})
	}
}
