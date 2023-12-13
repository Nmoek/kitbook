// Package web
// @Description: 帖子模块
package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	ijwt "kitbook/internal/web/jwt"
	"kitbook/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc            service.ArticleService
	interactiveSvc service.InteractiveService
	l              logger.Logger
	biz            string
}

func NewArticleHandler(svc service.ArticleService,
	interactiveSvc service.InteractiveService,
	l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:            svc,
		interactiveSvc: interactiveSvc,
		l:              l,
		biz:            "article", // 业务标识
	}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/articles")
	group.POST("/edit", a.Edit)         // 编辑帖子
	group.POST("/publish", a.Publish)   // 发表帖子
	group.POST("/withdraw", a.Withdraw) // 撤回帖子(更改可见状态)

	// 创作者接口
	group.GET("/detail/:id", a.Detail) // 帖子内容详情
	// /list?offset=?&limit=?  带参查询
	//group.GET("/list", a.List) // 创作列表
	// 查询参数放在Body中
	group.POST("/list", a.List)

	// 分第二个层次
	pub := group.Group("/pub")

	// 读者接口
	pub.GET("/:id", a.PubDetail) // 内嵌阅读数接口

	// 点赞接口
	// 传入参数, true=点赞, false=取消点赞
	pub.POST("/like", a.Like) // 内嵌阅读数接口

	// 收藏接口
	pub.POST("/collect", a.Collect)

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

// @func: Detail
// @date: 2023-12-05 02:09:13
// @brief: 查询创作列表-内容详情
// @author: Kewin Li
// @receiver a
// @param ctx
func (a *ArticleHandler) Detail(ctx *gin.Context) {

	var id int64
	var err error
	var claims ijwt.UserClaims
	var art domain.Article
	logKey := logger.ArticleLogMsgKey[logger.LOG_ART_DETAIL]
	fields := logger.Fields{}

	idStr := ctx.Param("id")
	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		a.l.WARN(logKey, fields.Add(logger.String("请求参数解析错误")).
			Add(logger.Error(err)).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Field{"idStr", idStr})...)

		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		return
	}

	claims = ctx.MustGet("user_token").(ijwt.UserClaims)

	art, err = a.svc.GetById(ctx, id)

	// 防攻击
	if art.Author.Id != claims.UserID {

		fields = fields.Add(logger.String("用户与帖子ID不匹配"))

		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		goto ERR
	}

	switch err {
	case nil:

		a.l.INFO(logKey, fields.Add(logger.String("查询列表详情成功")).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Field{"artId", id}).
			Add(logger.Int[int64]("userId", claims.UserID))...)

		ctx.JSON(http.StatusOK, Result{
			Msg:  "查询详情成功",
			Data: ConvertArticleVo(&art, false),
		})

		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}

ERR:
	a.l.ERROR(logKey,
		fields.Add(logger.Error(err)).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Field{"artId", id}).
			Add(logger.Int[int64]("userId", claims.UserID)).
			Add(logger.Field{"authorId", art.Author.Id})...)
	return

}

// @func: List
// @date: 2023-12-04 00:10:52
// @brief: 帖子模块-查询创作列表
// @author: Kewin Li
// @receiver a
// @param context
func (a *ArticleHandler) List(ctx *gin.Context) {
	var reqPage Page
	var err error
	var arts []domain.Article
	var claims ijwt.UserClaims
	logKey := logger.ArticleLogMsgKey[logger.LOG_ART_LIST]
	fields := logger.Fields{}

	err = ctx.Bind(&reqPage)
	if err != nil {
		fields = fields.Add(logger.String("请求解析失败"))
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		goto ERR
	}

	claims = ctx.MustGet("user_token").(ijwt.UserClaims)

	arts, err = a.svc.GetByAuthor(ctx, claims.UserID, reqPage.Offset, reqPage.Limit)

	switch err {
	case nil:

		a.l.INFO(logKey, fields.Add(logger.String("创作列表查询成功")).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Int[int64]("userID", claims.UserID))...)

		ctx.JSON(http.StatusOK, Result{
			Msg:  "查询列表成功",
			Data: ConvertArticleVos(arts, true),
		})

		return

	default:
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

	}

ERR:
	a.l.ERROR(logKey,
		fields.Add(logger.Error(err)).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Int[int64]("userId", claims.UserID))...)
	return
}

// @func: PubDetail
// @date: 2023-12-06 02:22:27
// @brief: 帖子模块-读者查询
// @author: Kewin Li
// @receiver a
// @param context
func (a *ArticleHandler) PubDetail(ctx *gin.Context) {
	logKey := logger.ArticleLogMsgKey[logger.LOG_ART_PUBDETAIL]
	fields := logger.Fields{}
	var claims ijwt.UserClaims
	var artId int64
	var err error
	var art domain.Article

	idStr := ctx.Param("id")
	claims = ctx.MustGet("user_token").(ijwt.UserClaims)

	artId, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fields = fields.Add(logger.String("请求参数解析失败"))
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		goto ERR
	}

	art, err = a.svc.GetPubById(ctx, artId)

	switch err {
	case nil:
		a.l.INFO(logKey,
			fields.Add(logger.String("读者查询成功")).
				Add(logger.Field{"IP", ctx.ClientIP()}).
				Add(logger.Int[int64]("artId", artId))...)

		ctx.JSON(http.StatusOK, Result{
			Msg: "查询成功",
			Data: ArticleVo{
				Id:         art.Id,
				Title:      art.Title,
				Content:    art.Content,
				AuthorId:   art.Author.Id,
				AuthorName: art.Author.Name,
				Status:     art.Status.ToUint8(),
				Ctime:      art.Ctime.Format(time.DateTime),
				Utime:      art.Utime.Format(time.DateTime),
			}})

		go func() {
			newCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// 阅读数+1
			err2 := a.interactiveSvc.IncreaseReadCnt(newCtx, a.biz, artId)
			if err2 != nil {
				a.l.ERROR(logKey,
					fields.Add(logger.Error(err2)).
						Add(logger.String("阅读数更新失败")).
						Add(logger.Field{"IP", ctx.ClientIP()}).
						Add(logger.Field{"artId", artId}).
						Add(logger.Int[int64]("userId", claims.UserID))...)
			}

		}()

		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Msg: "获取文章失败",
		})
	}

ERR:
	a.l.ERROR(logKey,
		fields.Add(logger.Error(err)).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Field{"artId", artId}).
			Add(logger.Int[int64]("userId", claims.UserID))...)
	return
}

// @func: Like
// @date: 2023-12-13 21:46:31
// @brief: 点赞/取消点赞
// @author: Kewin Li
// @receiver a
// @param c
func (a *ArticleHandler) Like(ctx *gin.Context) {
	type LikeReq struct {
		Id     int64 `json:"id"`
		IsLike bool  `json:"isLike"`
	}
	var req LikeReq
	var err error
	var claims ijwt.UserClaims

	logKey := logger.ArticleLogMsgKey[logger.LOG_ART_LIKE]
	fields := logger.Fields{}

	err = ctx.Bind(&req)
	if err != nil {
		fields = fields.Add(logger.String("请求参数解析错误"))
		goto ERR
	}

	claims = ctx.MustGet("user_token").(ijwt.UserClaims)

	if req.IsLike {

		err = a.interactiveSvc.Like(ctx, a.biz, req.Id, claims.UserID)
	} else {
		err = a.interactiveSvc.CancelLike(ctx, a.biz, req.Id, claims.UserID)

	}

	switch err {
	case nil:

		a.l.INFO(logKey, fields.Add(logger.String("点赞/取消点赞成功")).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Field{"artId", req.Id}).
			Add(logger.Field{"isLike", req.IsLike}).
			Add(logger.Int[int64]("userId", claims.UserID))...)
		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}

ERR:
	a.l.ERROR(logKey,
		fields.Add(logger.Error(err)).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Field{"artId", req.Id}).
			Add(logger.Field{"isLike", req.IsLike}).
			Add(logger.Int[int64]("userId", claims.UserID))...)
	return
}

// @func: Collect
// @date: 2023-12-14 01:37:55
// @brief: 帖子收藏
// @author: Kewin Li
// @receiver a
// @param c
func (a *ArticleHandler) Collect(ctx *gin.Context) {
	type CollectReq struct {
		// 帖子ID
		Id int64 `json:"id"`
		// 收藏夹ID
		CollectId int64 `json:"collectId"`
	}

	var req CollectReq
	var err error
	var claims ijwt.UserClaims
	logKey := logger.ArticleLogMsgKey[logger.LOG_ART_COLLECT]
	fields := logger.Fields{}

	err = ctx.Bind(&req)
	if err != nil {
		fields = fields.Add(logger.String("请求参数解析错误"))
		goto ERR
	}

	claims = ctx.MustGet("user_token").(ijwt.UserClaims)

	err = a.interactiveSvc.Collect(ctx, a.biz, req.Id, req.CollectId, claims.UserID)

	switch err {
	case nil:
		a.l.INFO(logKey, fields.Add(logger.String("收藏成功")).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Field{"artId", req.Id}).
			Add(logger.Int[int64]("userId", claims.UserID))...)

		ctx.JSON(http.StatusOK, Result{
			Msg: "收藏成功",
		})
		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Msg: "收藏失败",
		})
	}

ERR:
	a.l.ERROR(logKey,
		fields.Add(logger.Error(err)).
			Add(logger.Field{"IP", ctx.ClientIP()}).
			Add(logger.Field{"artId", req.Id}).
			Add(logger.Int[int64]("userId", claims.UserID))...)
	return
}
