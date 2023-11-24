// Package web
// @Description: 帖子模块
package web

import (
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
// @brief: 帖子模块-编辑后保存，不发表(无论新建、修改)
// @author: Kewin Li
// @receiver a
// @param context
// @接收文章内容输入，返回文章的ID
func (a *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
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
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.UserID,
		},
	})

	switch err {
	case nil:
		a.l.INFO(logKey, fileds.Add(logger.Field{"success", "帖子保存成功"}).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Int[int64]("artId", req.Id)).
			Add(logger.Int[int64]("userId", claims.UserID))...)

		ctx.JSON(http.StatusOK, Result{
			Msg:  "保存成功",
			Data: artId,
		})

		return
	case service.ErrInvalidUpdate:
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法操作",
			Data: artId,
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}

	fileds = fileds.Add(logger.Error(err)).
		Add(logger.Field{"IP", ctx.ClientIP()}).
		Add(logger.Int[int64]("artId", req.Id)).
		Add(logger.Int[int64]("userId", claims.UserID))

ERR:
	a.l.ERROR(logKey, fileds...)
	return
}
