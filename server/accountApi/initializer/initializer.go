package initializer

import (
	"context"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/luxun9527/zlog"
	"github/lunxun9527/bestpractice/pkg/i18n"
	"github/lunxun9527/bestpractice/pkg/xvalidator"
	"github/lunxun9527/bestpractice/server/accountApi/config"
	"github/lunxun9527/bestpractice/server/accountApi/global"
	"github/lunxun9527/bestpractice/server/accountApi/rpcClient"
)

func Init(confPath string) {
	//初始化配置
	global.Config = config.InitConfig(confPath)
	//初始化日志
	zlog.InitDefaultLogger(&global.Config.Logger)
	//初始化grpc客户端
	client, err := global.Config.RpcClient.EtcdConf.NewEtcdClient()
	if err != nil {
		zlog.Panicf("etcd client init failed, err:%v", err)
	}
	if err := rpcClient.InitEtcdRpcClients(context.Background(), client, global.Config.RpcClient.TargetConfList); err != nil {
		zlog.Panicf("rpc client init failed, err:%v", err)
	}
	//加载多语言文件
	translator, err := i18n.NewTranslatorFormFile(global.Config.Lang.Path)
	if err != nil {
		zlog.Panicf("i18n init failed, err:%v", err)
	}
	i18n.SetDefaultTranslator(translator)

	//设置gin参数校验
	v, _ := xvalidator.NewValidateTranslator(binding.Validator.Engine().(*validator.Validate))
	xvalidator.SetDefaultValidateTranslator(v)
	xvalidator.RegisterValidations()
}
