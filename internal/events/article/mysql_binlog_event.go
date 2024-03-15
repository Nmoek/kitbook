package article

import (
	"context"
	"github.com/IBM/sarama"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	"kitbook/internal/repository/dao"
	"kitbook/pkg/canalx"
	"kitbook/pkg/logger"
	"kitbook/pkg/saramax"
	"time"
)

const (
	TopicBinlogEvent = "kitbook_binlog"
)

type MySqlBinlogConsumer struct {
	repo   *repository.CacheArticleRepository
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
			saramax.NewHandler[canalx.Message[dao.PublishedArticle]](m.Consume, m.l))
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
	event canalx.Message[dao.PublishedArticle]) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if event.Table != "published_articles" {
		return nil
	}

	for _, data := range event.Data {
		var err error
		switch data.Status {
		case domain.ArticleStatusPublished:

			err = m.repo.Cache().SetPub(ctx, repository.ConvertsDomainArticleFromLive(&data))
			if err != nil {
				return err
			}
		case domain.ArticleStatusPrivate:
			err = m.repo.Cache().DelPub(ctx, data.Id)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *MySqlBinlogConsumer) StartV2() error {
	cg, err := sarama.NewConsumerGroupFromClient("pub_articles_cache", m.client)
	if err != nil {
		return err
	}

	go func() {

		err2 := cg.Consume(context.Background(), []string{TopicBinlogEvent},
			saramax.NewBatchHandler[canalx.Message[dao.PublishedArticle]](m.BatchConsume, m.l))
		if err2 != nil {
			//TODO: 日志埋点

		}

	}()

	return nil
}

// @func: Consume
// @date: 2023-12-19 03:06:25
// @brief: 帖子模块-实际消费业务处理-批量提交
// @author: Kewin Li
// @receiver i
// @param msg
// @param event
// @return error
func (m *MySqlBinlogConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	event []canalx.Message[dao.PublishedArticle]) error {

	return nil

}
