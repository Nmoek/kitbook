// Package web
// @Description: 帖子模块
package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	"kitbook/internal/web/jwt"
	"kitbook/pkg/logger"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleServer
	l   logger.Logger
}

func NewArticleHandler(svc service.ArticleServer, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/articles")
	group.POST("/edit", a.Edit) // 编辑帖子

}

// @func: Edit
// @date: 2023-11-22 22:47:05
// @brief: 帖子模块-编辑后保存，不发表
// @author: Kewin Li
// @receiver a
// @param context
// @接收文章内容输入，返回文章的ID
func (a *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req Req
	var err error
	var artId int64
	logKey := logger.ArticleLogMsgKey[logger.LOG_ART_EDIT]
	claims := jwt.UserClaims{}
	fileds := logger.Fields{}

	err = ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		fileds = fileds.Add(logger.Error(err))
		goto ERR
	}

	// 作者Id通过jwt来解析
	claims = ctx.MustGet("user_token").(jwt.UserClaims)

	// 保存
	artId, err = a.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.UserID,
		},
	})

	if err != nil {

		ctx.JSON(http.StatusOK, Result{
			Msg: "保存失败",
		})

		fileds = fileds.Add(logger.Error(err)).Add(
			logger.Int[int64]("userId", claims.UserID),
		).Add(
			logger.Field{"title", req.Title},
		)

		goto ERR
	}

	a.l.INFO(logKey, logger.Field{
		"success",
		fmt.Sprintf("[%s], [%d][%d] 帖子保存成功", ctx.ClientIP(), artId, claims.UserID),
	})
	ctx.JSON(http.StatusOK, Result{
		Msg:  "保存成功",
		Data: artId,
	})

	return

ERR:
	a.l.ERROR(logKey, fileds...)
	return
}
