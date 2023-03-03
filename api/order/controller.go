package order

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shopping/domain/order"
	"shopping/utils/api_helper"
	"shopping/utils/pagination"
)

type Controller struct {
	orderService *order.Service
}

func NewOrderController(orderService *order.Service) *Controller {
	return &Controller{
		orderService: orderService,
	}
}

// CompleteOrder godoc
// @Summary 完成订单
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} api_helper.Response
// @Failure 400 {object} api_helper.ErrResponse
// @Router /order [post]
func (c *Controller) CompleteOrder(g *gin.Context) {
	userId := api_helper.GetUserId(g)
	err := c.orderService.CompleteOrder(userId)
	if err != nil {
		api_helper.HandlerError(g, err)
		return
	}
	g.JSON(
		http.StatusCreated, api_helper.Response{
			Message: "Order Created",
		})
}

// CancelOrder godoc
// @Summary 取消订单
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param CancelOrderRequest body CancelOrderRequest true "order information"
// @Success 200 {object} api_helper.Response
// @Failure 400 {object} api_helper.ErrResponse
// @Router /order [delete]
func (c *Controller) CancelOrder(g *gin.Context) {
	var req CancelOrderRequest
	if err := g.ShouldBind(&req); err != nil {
		api_helper.HandlerError(g, err)
		return
	}
	userId := api_helper.GetUserId(g)
	err := c.orderService.CancelOrder(userId, req.OrderId)
	if err != nil {
		api_helper.HandlerError(g, err)
		return
	}
	g.JSON(
		http.StatusCreated, api_helper.Response{
			Message: "Order Canceled",
		})
}

// GetOrders godoc
// @Summary 获得订单列表
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} pagination.Pages
// @Router /order [get]
func (c *Controller) GetOrders(g *gin.Context) {
	page := pagination.NewFromGinRequest(g, -1)
	userId := api_helper.GetUserId(g)
	page = c.orderService.GetAll(page, userId)
	g.JSON(http.StatusOK, page)
}
