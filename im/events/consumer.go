package events

import (
	"context"
	"github.com/IBM/sarama"
	"kitbook/im/domain"
	"kitbook/im/service"
	"kitbook/pkg/canalx"
	"kitbook/pkg/logger"
	"kitbook/pkg/saramax"
	"strconv"
	"time"
)

const (
	TopicBinlogEvent = "kitbook_binlog"
)

type MySqlBinlogConsumer struct {
	svc    service.UserService
	client sarama.Client
	l      logger.Logger
}

// @func: Start
// @date: 2023-12-17 20:25:40
// @brief: 启动消费
// @author: Kewin Li
// @receiver i
// @return error
func (m *MySqlBinlogConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("pub_articles_cache", m.client)
	if err != nil {
		return err
	}

	go func() {

		err2 := cg.Consume(context.Background(), []string{TopicBinlogEvent},
			saramax.NewHandler[canalx.Message[User]](m.Consume, m.l))
		if err2 != nil {
			//TODO: 日志埋点

		}

	}()

	return nil
}

// @func: Consume
// @date: 2023-12-17 20:31:03
// @brief: 帖子模块-实际消费业务处理-阅读数+1
// @author: Kewin Li
// @receiver i
// @param msg
// @param events
// @return error
func (m *MySqlBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	event canalx.Message[User]) error {
	if event.Table != "users" || event.Type != "INSERT" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	for _, data := range event.Data {
		err := m.svc.Sync(ctx, domain.User{
			UserID:   strconv.FormatInt(data.Id, 10),
			Nickname: data.Nickname,
		})
		if err != nil {
			// TODO: 记录日志
			continue
		}
	}

	return nil
}

type User struct {
	Id    int64  `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`

	Openid   string `json:"openid,omitempty"`
	Unionid  string `json:"unionid,omitempty"`
	Password string `json:"password,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Birthday int64  `json:"birthday,omitempty"`
	AboutMe  string `json:"about_me,omitempty"`
	Ctime    int64  `json:"ctime,omitempty"`
	Utime    int64  `json:"utime,omitempty"`
}
