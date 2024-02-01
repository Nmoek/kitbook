package ioc

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"kitbook/payment/repository"
	"kitbook/payment/service/wechat"
	"kitbook/pkg/logger"
	"os"
)

func InitWechatClient(cfg WechatConfig) *core.Client {
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.KeyPath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	client, err := core.NewClient(ctx,
		option.WithWechatPayAutoAuthCipher(
			cfg.AppID, cfg.MchSerialNum,
			mchPrivateKey, cfg.MchKey))
	if err != nil {
		panic(err)
	}

	return client
}

func InitWechatNativeService(cli *core.Client,
	repo repository.PaymentRepository,
	l logger.Logger,
	cfg WechatConfig) *wechat.NativePaymentService {
	return wechat.NewNativePaymentService(cfg.AppID, cfg.MchId, repo, &native.NativeApiService{
		Client: cli,
	}, l)
}

func InitWechatNotifyHandler(cfg WechatConfig) *notify.Handler {
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(cfg.MchId)
	handler, err := notify.NewRSANotifyHandler(cfg.MchKey, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		panic(err)
	}

	return handler
}

func InitWechatConfig() WechatConfig {
	return WechatConfig{
		AppID: os.Getenv("WEPAY_APP_ID"),
		MchId: os.Getenv("WEPAY_MCH_ID"),

		MchKey: os.Getenv("WEPAY_MCH_KEY"),

		MchSerialNum: os.Getenv("WEPAY_MCH_SERIAL_NUM"),
		//CertPath:     "./config/cert/apiclient_cert.pem",
		//KeyPath:      "./config/cert/apiclient_key.pem",
	}
}

type WechatConfig struct {
	AppID        string
	MchId        string
	MchKey       string
	MchSerialNum string

	//证书
	KeyPath  string
	CertPath string
}
