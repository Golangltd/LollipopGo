# LollipopGo 
Golang语言情怀  全球服游戏服务器框架,目前协议支持websocket、http及RPC，采用状态同步，愿景：打造竞技实时【比赛】对战游戏平台框架！ 功能持续更新中... ...
>微信订阅号：Golang语言情怀<Br/>
>微信服务号：Golang语言游戏服务器<Br/>
>商业定制版：联系彬哥(微信：cserli)<Br/>


论坛
--------------
WwW.Golang.Ltd

LollipopGo框架交流群
-----------
955259501

Go语言交流群
----------
221273219(已满)

713828896


框架示例
----------
[游戏服务器示例(登录服务器，代理服务器等)](https://github.com/Golangltd/LollipopGo_Server_Client)


腾讯云+社区专栏
-----------
[腾讯专栏](https://cloud.tencent.com/developer/column/2170)


Golang语言情怀
-----------

<ol>
<li>希望更多喜欢Go语言的同学及想从事Go语言开发游戏服务器的同学一个方向的指引</li>
<li>课程多维度教学，lollipopGo游戏框架实战课程等等</li>
<li>LollipopGo架构 最新版本: v2.8.X </li>
<li>LollipopGo架构 直接下载就可以使用（彬哥维护），无需依赖管理，否则导致部分官方接口无法使用 </li>
<li>LollipopGo架构 手机对战游戏视频：<a href="https://www.bilibili.com/video/av52239498" target="_blank">点击访问</a></li>
<li>LollipopGo架构 PC端游对战游戏视频：<a href="https://www.bilibili.com/video/av54726431" target="_blank">点击访问</a></li>
<li>同时我们的免费课程也在持续更新中; 点击访问：<a href="http://gopher.ke.qq.com" target="_blank">腾讯课堂</a></li>
<li>同时我们的免费课程也在持续更新中; 点击访问：<a href="https://study.163.com/provider/400000000538037/index.htm?share=2&shareId=400000000538037" target="_blank">网易云课堂</a></li>
<li>同时我们的免费课程也在持续更新中; 点击访问：<a href="http://space.bilibili.com/389368547?" target="_blank">B站(bilibili.com)</a></li>
<li>同时我们的免费课程也在持续更新中; 点击访问：<a href="http://www.gameais.com" target="_blank">LollipopGo框架文档地址：GameAIs平台(GameAIs.com)</a></li>
</ol>





架构目录说明
-----------
```
├── encrypt
│   ├── AES           # AES加密           
│   ├── DES           # DES加密
│   └── RSA           # RSA加密
├── error             # 错误统一格式管理模块，全局错误码定义
├── global_Interface  # 网络接口定义，分布式服务器需要单独实现接口
├── heartbeat         # 通用心跳模块
├── leaf              # leaf的一些扩展函数，包括自定义的protobuf消息解析器
├── log               # 通用的log封装
├── network           # 网络处理封装，目前支持：http、rpc、websocket
├── PhysX             # 3D游戏Physx支持，目前内测
│   ├── DLLAndSo      # PhysX动态链接库，支持windows系统及linux系统
│   ├── PhysXAPI      # Go语言封装的PhysXAPI调用接口
│   └── example       # Go语言封装接口测试实例
├── Proxy_Server      
│   └──  Proto        # 反向代理消息公用模块，框架标准
├── SM                # 游戏AI处理,目前支持有限状态机
├── standard          # 框架函数、变量命名规则(建议驼峰命名)
├── timer             # 通用定时器
├── tools
│   ├── collection    # 集合类的扩展方法
│   ├── database      # 快速初始化数据库连接
│   ├── deepcopy      # 通用深拷贝（使用反射）
│   ├── DFA           # 过滤敏感字
│   ├── fs            # 文件系统/配置解析
│   ├── ip            # ip地址库
│   ├── jsonutils     # json工具库
│   ├── mem           # 常用的内存缓存类
│   ├── num           # 基础数字类型工具函数
│   ├── sample        # 随机抽样函数
│   └── tz            # 时间函数
└── util              # 随机数，并发安全map、排序等相关公用接口
```
 <div class="footer">

 </div>
