## go私有库目录结构

私有仓库中一般会存什么内容，应该如何组织目录结构。

1、将中间件sdk的再次封装，主要是将中间件各种配置项加载后，初始化实例。能用框架就是框架没有再自己搞。其他框架有的，就抄过来，推荐目录名字 

**pkg/gormx或pkg/xgorm**

**pkg/kafkax pkg/xkafka** 

**pkg/xnacos pkg/nacosx 最好不要和官方库同名**

2、grpc pb协议。推荐是分仓库，不要多个项目的协议放在一个仓库中，避免相互影响。推荐 

**pb/projectname/** 

3、公共model，预定义错误，redis key，以及各种常量。推荐目录名字 

**common/model/项目名.go**

**common/rediskey/项目名称.go** 

**common/constrant/项目名.go** 

**common/errs/项目名.go**

4、utils：类型装换，hash方法，生成随机数。 一般只有一个函数。能用开源就用开源不要自己写。推荐目录名字 	**common/utils/timex.go**

 **common/utils/stringx.go**

推荐几个常用的https://github.com/gookit/goutil，https://github.com/duke-git/lancet/tree/main

如果是只有在自己用到的方法，不用放到common 直接在项目下定义projectName/tools 和utils做区分。



## 快速curd

**后端的核心基本curd, 我比较推荐的快速curd的方法。 表结构驱动，gorm+gen组合**

1、业务相关的单词提前定义，避免一个意思多个单词，积极维护相关的文档。

2、数据类型提前定义，如涉及到钱的应该要提前定义好精度，数据库大家统一使用decimal字段。

​     数据表如果存储大量的数据要定义每一列的数据类型。只要是值类型，int 能存下的不用bigint。

​     时间类型统一时间戳。使用int64类型，主要是考虑到如果要存毫秒等情况，否则无符号int32也是可以

​     常用的表、且数据量比较大，索引设计要慎重。

3、使用创建表的sql，通过grom 的gen生成对于的model，要提前定义好类型表结构到结构体的映射， 如果一个程序会用到多个数据库。要提前设计包的结构 一般的结构是 projectname/dao/model  如果用了多个数据库 则为 projectname/dao/dbname/model

**重点，提前统一这些会在后期的开发中，避免很多很多坑。**

## 微服务

**为啥微服务，我的理解**

1、业务复杂，如果一直在一个单体应用上改。后面会搞的难以维护。

2、降低风险，提高容错。微服务一个服务挂了，影响的范围会小一点。

3、方便横向扩容。

4、一致性要求不强的情况，使用消息组件，可以提高接口的响应速度。



**1、避免api和rpc接口大面积重复**

大部分微服务都是 请求-->nginx--(使用http)-->api服务---(使用grpc)-->rpc服务

我们开发的时候，我们大部分的业务一般都是写在rpc里面，但是一般rpc服务不直接面向前端，这种情况，我们一般要在api服务里面做一层转换。示例

```go
func (l *GetOrderListLogic) GetOrderList(req *types.GetOrderListReq) (resp *types.GetOrderListResp, err error) {

    orderList, err := client.GetOrderList(l.ctx, &orderpb.GetOrderListByUserReq{
        UserId:     cast.ToInt64(uid),
        StatusList: statusList,
        PageSize:   req.PageSize,
        Id:         cast.ToInt64(req.Id),
    })
    if err != nil {
        return nil, err
    }
    orderInfoList := make([]*types.OrderInfo, 0, len(orderList.OrderList))
    for _, v := range orderList.OrderList {
        orderInfo := &types.OrderInfo{
            Id:             cast.ToString(v.Id),
            OrderId:        v.OrderId,
            UserId:         v.UserId,
            SymbolName:     v.SymbolName,
            Price:          v.Price,
            Qty:            v.Qty,
            Amount:         v.Amount,
            Side:           int32(v.Side),
            Status:         int32(v.Status),
            OrderType:      int32(v.OrderType),
            FilledQty:      v.FilledQty,
            FilledAmount:   v.FilledAmount,
            FilledAvgPrice: v.FilledAvgPrice,
            CreatedAt:      v.CreatedAt,
        }
        orderInfoList = append(orderInfoList, orderInfo)
    }
    resp = &types.GetOrderListResp{OrderList: orderInfoList, Total: orderList.Total}
    return
}
```

推荐解决方法

1、使用grpc-gateway，apisix 让rpc服务直接提供接口。grpc-gateway抽取为一个单独的服务，比较推荐。

2、或者在api服务中做业务逻辑。

## 可观测

**日志，链路追踪，指标采集。**

**日志**   elk 。elk 依赖统一的日志格式，一定要提前封装日志组件.。一定要搞。

**链路追踪** jaeger。如何你的调用链复杂，且业务比较重要，可以搞。一般的curd没啥必要，搞了一般也不看。

**指标采集** prometheus。如果有专业的运维团队，且时间充裕可以搞。如果时间不充裕，且运维人手不够。不推荐。



## 错误体系

**以下是我在个人在日常开发中觉得比较好的实践。**

1、预定义错误使用grpc的status来预定义,预定义错误，错误码按照模块区分。api和rpc的业务错误都使用预定的grpc status错误，这样api不用关心grpc的错误直接抛出就可以，在api统一处理。

2、后端的grpc返回业务错误的时候使用error来返回，而不是在message中定义code,message的形式

3、多语言场景，翻译在api层统一对错误处理。如果需要翻译支持动态修改，修改无需重启程序可以使用fsnotify，etcd，nacos来实现。

4、错误统一存放，统一格式，最好使用工具生成。如go:generate 加上stringer，如果不能满足自己的需求可以使用ast包来自己自定义。





## 目录结构

如果使用了开源框架就按开源的来。

api  服务参考 gin-vue-admin

rpc 团队有要求就按照团队的来，否则基本也就是  api/controler--> service/logic-->dao。
