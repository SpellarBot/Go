# 一, ABtest框架设计
## 将流量切分为独占区和并行区
   1. 并行区提供并行实验环境，层与层之间流量重新hash，使上层对下层所有实验的影响均衡，从而使同一层间多个实验效果可比；并行层可提供人工运营优化，算法策略优化小流量大流量实验
   2. 独占区提供独占实验环境，每份流量只能做一个实验，不受其他实验影响，效果可以更直接地衡量；独占区可提供人工运营灰度发布，算法效果验证

## 提供多种切流hash方式
   1. F(userid,layerid) 按用户和层ID进行hash，在同一个层，用户一直在同一个桶中，方便短期效果对比
   2. F(userid,layerid,date) 按用户和层ID，并考虑上有按时间进行流量轮换的需求，比如每周轮换流量进行效果对比
   3. F(userid,layerid,Random) 直接随机，针对一次性抽样流量对比效果

## 配置中心化
   1. 当前版本定时读取配置文件，更新配置时，需要同步到各台机器上，后期有需求再从共享存储如：web服务器，中定时下载。
   2. 虽然从本地加载，但不需要重启服务。
   3. 独占区的切流比例也是可配置的，甚至为0，独占流量自动从并行区排除。

## 实验流程化
   1. 这一块非必须，可提供一个简单界面，用于创建实验，管理配置参数，实验可以一目了然，后期算法可视化时增加。

# 二, 实现和接口
   1. 创建ABtestConfig，并调用其Init(filePath string)方法，就可使用。filePath为本地绝对文件路径
   2. 使用时调用GetABTag(userId string, layerId string), GetABTag2(userId string, layerId string, option string),GetMutexABTag(userId string) 获取对应的tag, 如果第二个参数返回为false,表明这流量是需要排除。
   3. 各模块按需增加切流代码及日志回流代码，**需要注意的是，每个ABtest必须实现abtag为default的分支代码，以实现基线的功能**
   4.（可选）增加一个管理界面

## 接口说明
   1. 首先需要调用: Init(filePath string)
   参数:
   filePath为本地绝对文件路径
   period: 自动更新配置的周期,-1为不需要自动更新，只获取一次, 单位为分钟。
   2. 获取用户abtag: getABTag， GetABTag2， GetMutexABTag
   参数:
   userId   必选，用户ID
   layerId  为空为获取独占区实验的用户abtag，否则获取并行区相应layer的用户的abtag
   option   可选，如果需要按天或按星期轮换用户，可输入日期或星期字符串，如果需要random，则输入随机生成的字符串

## ABtest 配置说明
   2. XML的配置说明：
   bucketnum 用户桶数，比如10000000即将用户分10000000桶， 桶越大分流分得越好。
   mutex/range 独占区流量比例，比如0.2，即将2000000桶用户流量作为独占区
   layer/id 并行区层ID, 与ABtestUtils.getABTag的layerId参数对应
   domain/name 实验abtag，与ABtestUtils.getABTag返回的abtag对应
   domain/range 实验的流量比例










