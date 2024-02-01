package web

import (
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"kitbook/payment/service/wechat"
	"kitbook/pkg/logger"
	"net/http"
)

type WeChatNativeHandler struct {
	handler *notify.Handler

	svc *wechat.NativePaymentService
	l   logger.Logger
}

func (w *WeChatNativeHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/pay/callback", w.HandleNative)
}

func NewWeChatNativeHandler(handler *notify.Handler, svc *wechat.NativePaymentService, l logger.Logger) *WeChatNativeHandler {
	return &WeChatNativeHandler{
		handler: handler,
		svc:     svc,
		l:       l,
	}
}

func (w *WeChatNativeHandler) HandleNative(ctx *gin.Context) {
	transaction := new(payments.Transaction)

	_, err := w.handler.ParseNotifyRequest(ctx, ctx.Request, transaction)
	if err != nil {
		ctx.String(http.StatusBadRequest, "参数解析失败")
		w.l.ERROR("notifyURL解析失败",
			logger.Error(err),
			logger.Field{"peer_id", ctx.ClientIP()},
			logger.Field{"notifyURL", ctx.Request.URL})
		//TODO: 接入系统观测、监控
	}

	err = w.svc.HandleCallback(ctx, transaction)
	if err != nil {
		w.l.ERROR("回调处理失败",
			logger.Error(err),
			logger.Field{"biz_trade_no", transaction.OutTradeNo})
		ctx.String(http.StatusInternalServerError, "系统异常")
	}

	ctx.String(http.StatusOK, "处理完毕")
}
