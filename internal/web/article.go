// Package web
// @Description: 帖子模块
package web

import (
	"github.com/gin-gonic/gin"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	ijwt "kitbook/internal/web/jwt"
	"kitbook/pkg/logger"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.Logger
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/articles")
	group.POST("/edit", a.Edit)         // 编辑帖子
	group.POST("/publish", a.Publish)   // 发表帖子
	group.POST("/withdraw", a.Withdraw) // 撤回帖子(更改可见状态)

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
	claims := ijwt.UserClaims{}
	fileds := logger.Fields{}

	err = ctx.Bind(&req)
	if err != nil {
		fileds = fileds.Add(logger.String("请求解析错误"))

		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		goto ERR
	}

	// 作者Id通过jwt来解析
	claims = ctx.MustGet("user_token").(ijwt.UserClaims)

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
		a.l.INFO(logKey,
			fileds.Add(logger.String("帖子保存成功")).
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
			Msg: "保存失败",
		})
	}

ERR:
	a.l.ERROR(logKey, fileds.Add(logger.Error(err)).
		Add(logger.Field{"IP", ctx.ClientIP()}).
		Add(logger.Int[int64]("artId", req.Id)).
		Add(logger.Int[int64]("userId", claims.UserID))...)
	return
}

// @func: Publish
// @date: 2023-11-26 00:00:30
// @brief: 帖子模块-帖子发表
// @author: Kewin Li
// @receiver a
// @param context
func (a *ArticleHandler) Publish(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req Req
	var err error
	var artId int64
	logKey := logger.ArticleLogMsgKey[logger.LOG_ART_EDIT]
	claims := ijwt.UserClaims{}
	fileds := logger.Fields{}

	err = ctx.Bind(&req)
	if err != nil {
		fileds = fileds.Add(logger.String("请求解析失败"))
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		goto ERR
	}

	// 作者Id通过jwt来解析
	claims = ctx.MustGet("user_token").(ijwt.UserClaims)

	// 发表
	artId, err = a.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.UserID,
		},
	})

	switch err {
	case nil:
		a.l.INFO(logKey,
			fileds.Add(logger.String("帖子发表成功")).
				Add(logger.Field{"IP", ctx.ClientIP()}).
				Add(logger.Int[int64]("artId", artId)).
				Add(logger.Int[int64]("userId", claims.UserID))...)

		ctx.JSON(http.StatusOK, Result{
			Msg:  "发表成功",
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
			Msg:  "发表失败",
			Data: -1,
		})
	}

ERR:
	a.l.ERROR(logKey,
		fileds.Add(logger.Error(err)).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Int[int64]("artId", req.Id)).
			Add(logger.Int[int64]("userId", claims.UserID))...)
	return
}

// @func: Withdraw
// @date: 2023-11-28 12:33:42
// @brief: 帖子模块-帖子撤回
// @author: Kewin Li
// @receiver a
// @param context
func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req Req
	var err error

	logKey := logger.ArticleLogMsgKey[logger.LOG_ART_WITHDRAW]
	claims := ijwt.UserClaims{}
	fields := logger.Fields{}

	err = ctx.Bind(&req)
	if err != nil {
		fields = fields.Add(logger.String("请求解析失败"))
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		goto ERR
	}

	// 作者Id通过jwt来解析
	claims = ctx.MustGet("user_token").(ijwt.UserClaims)

	// 撤回
	err = a.svc.Withdraw(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.UserID,
		},
	})

	switch err {
	case nil:
		a.l.INFO(logKey,
			fields.Add(logger.String("帖子撤回成功")).
				Add(logger.Field{"IP", ctx.ClientIP()}).
				Add(logger.Int[int64]("artId", req.Id)).
				Add(logger.Int[int64]("userId", claims.UserID))...)

		ctx.JSON(http.StatusOK, Result{
			Msg:  "撤回成功",
			Data: req.Id,
		})
		return
	case service.ErrInvalidUpdate:
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法操作",
			Data: -1,
		})

	default:
		ctx.JSON(http.StatusOK, Result{
			Msg:  "撤回失败",
			Data: -1,
		})
	}

ERR:
	a.l.ERROR(logKey,
		fields.Add(logger.Error(err)).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Int[int64]("artId", req.Id)).
			Add(logger.Int[int64]("userId", claims.UserID))...)
	return
}
