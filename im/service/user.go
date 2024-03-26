package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"kitbook/im/domain"
	"net/http"
)

type UserService interface {
	Sync(ctx context.Context, user domain.User) error
}

type RESTUserService struct {
	base string
	// 默认是 openIM123
	secret string
	client *http.Client
}

func NewRESTUserService(base string, secret string) *RESTUserService {
	return &RESTUserService{
		base:   base,
		secret: secret,
		client: http.DefaultClient,
	}
}

func (r *RESTUserService) Sync(ctx context.Context, user domain.User) error {
	var operationID string
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		operationID = spanCtx.TraceID().String()
	} else {
		operationID = uuid.New().String()
	}

	reqBody := request{
		Secret: r.secret,
		Users:  []domain.User{user},
	}
	val, _ := json.Marshal(&reqBody)

	req, err := http.NewRequest(http.MethodPost, r.base+"/user/user_register", bytes.NewReader(val))
	if err != nil {
		return err
	}
	req.Header.Add("operationID", operationID)

	httpResp, err := r.client.Do(req)
	if err != nil {
		return err
	}

	var resp response
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return err
	}

	//var resp response
	//err := httpx.NewRequst(ctx, http.MethodPost, r.base+"/user/user_register").AddHeader("operationID", operationID).
	//	JSONBody(request{
	//		Secret: r.secret,
	//		Users:  []domain.User{user},
	//	}).Do().JSONScan(&resp)
	//if err != nil {
	//	return err
	//}

	if resp.ErrCode != 0 {
		return fmt.Errorf("同步用户数据失败, %v \n", err)
	}

	return nil
}

type request struct {
	Secret string        `json:"secret"`
	Users  []domain.User `json:"users"`
}

type response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
}
