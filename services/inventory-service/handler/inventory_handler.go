// package handler

// import (
// 	"net/http"

// 	"github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
// 	"github.com/gin-gonic/gin"
// )

// func RegisterRoutes(r *gin.Engine, svc inventorypb.InventoryServiceServer) {
// 	inv := r.Group("/inventory")
// 	{
// 		inv.GET("/:id", func(c *gin.Context) {
// 			resp, err := svc.GetInventory(c.Request.Context(), &inventorypb.GetInventoryRequest{
// 				ProductId: c.Param("id"),
// 			})
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			c.JSON(http.StatusOK, resp)
// 		})

// 		inv.POST("/:id", func(c *gin.Context) {
// 			var body struct {
// 				Change int32 `json:"change"`
// 			}
// 			if err := c.ShouldBindJSON(&body); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
// 				return
// 			}

// 			resp, err := svc.UpdateInventory(c.Request.Context(), &inventorypb.UpdateInventoryRequest{
// 				ProductId:      c.Param("id"),
// 				QuantityChange: body.Change,
// 			})
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			c.JSON(http.StatusOK, resp)
// 		})
// 	}
// }

package handler

import (
	"net/http"

	"github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	"google.golang.org/grpc/metadata"

	// "github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/middleware" // 🔐 Import JWT middleware
	"github.com/BhavaniNBL/ecommerce-backend/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, svc inventorypb.InventoryServiceServer) {
	inv := r.Group("/inventory")
	inv.Use(middleware.JWTMiddleware()) // 🔐 Protect all inventory routes with JWT

	inv.GET("/:id", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)
		resp, err := svc.GetInventory(ctx, &inventorypb.GetInventoryRequest{
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

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}

		// ✅ 2. Inject into gRPC metadata
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

		resp, err := svc.UpdateInventory(ctx, &inventorypb.UpdateInventoryRequest{
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
