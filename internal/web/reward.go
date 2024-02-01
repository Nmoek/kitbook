package web

import (
	"github.com/gin-gonic/gin"
	rewardv1 "kitbook/api/proto/gen/reward/v1"
	ijwt "kitbook/internal/web/jwt"
	"kitbook/pkg/logger"
	"net/http"
)

// RewardHandler
// @Description: web接口-打赏模块
type RewardHandler struct {
	rewardSvc rewardv1.RewardServiceClient

	l logger.Logger
}

func (h *RewardHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("reward")

	group.Any("/detail", h.GetReward)
}

// @func: GetWard
// @date: 2024-02-05 23:54:31
// @brief: web接口-轮询查询打赏结果
// @author: Kewin Li
// @receiver h
// @param ctx
func (h *RewardHandler) GetReward(ctx *gin.Context) {
	type GetRewardReq struct {
		Rid int64 `json:"rid"`
		Amt int64 `json:"amt"`
	}

	var req GetRewardReq
	var err error

	err = ctx.Bind(&req)
	if err != nil {

		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		return
	}

	// 作者Id通过jwt来解析
	claims := ctx.MustGet("user_token").(ijwt.UserClaims)

	_, err = h.rewardSvc.GetReward(ctx, &rewardv1.GetRewardRequest{
		Rid: req.Rid,
		Uid: claims.UserID,
	})
	if err != nil {

	}

	ctx.JSON(http.StatusOK, Result{})
}
