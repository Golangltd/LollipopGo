Golang语言社区  开源项目 （持续更新）<br>
============================================
[点击进入社区论坛](http://www.Golang.LTD "悬停显示") 
<br>
<br>
<br>


环境依赖
Go1.9以上

部署步骤
1. 配置数据库(选择其中一种即可)
    config.csv
	config.go
	config.json

2. 部署go执行文件（仅支持linux）
    deploy.sh

<br>
##################################<br>
>>目录结构 如下：

<br>
<br>

├── LollipopGo<br>
│---── config<br>
│------├──config.csv<br>
│------├──config.go<br>
│------├──config.json<br>
│------└──README.txt<br>
│---── deploy<br>
│------├──allclose.sh<br>
│------├──deploy.sh<br>
│------└──restart.sh<br>
│---── lang<br>
│------├──limlt_versoin.go<br>
│------└──zh-cn.go<br>
│---── library<br>
│------└──lollipop<br>
│----------├──cache<br>
│----------├──common<br>
│----------├──controller<br>
│----------├──db<br>
│----------├──globalfun<br>
│----------├──log<br>
│----------├──redis<br>
│----------├──template<br>
│----------├──code.google.com<br>
│----------├──concurrentMap<br>
│----------└──Build.go<br>
│------──traits<br>
│----------└──traits.go<br>
│---── tpl<br>
│------└──default_index.tpl<br>
├── base.go<br>
├── help.go<br>
└── README.md<br>

<br>
<br>

########### V1.0.1 版本内容更新##############
1. 更新base.go 初始化逻辑
2. 更新安全并发map及cache 初始化


<br>
注：2018年由于相对比较忙，预计2019年年初完成，同时用LollipopGo建立3个视频实战项目；具体请关注公众平台文章。<br>  

<br>  
请关注Golang语言社区公众平台ID：Golangweb<br>
![](https://github.com/Golangltd/LollipopGo/blob/master/t7e102owue.png)

