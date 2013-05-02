
# Fetch Weibo

## 简介

使用新浪微博的API，抓取微博的朋友和用户的timeline。用于社交网络中的朋友推荐算法研究。

## 如何工作

1. 抓取某一个用户的timeline以及朋友ID。
2. 取这个用户朋友的前10个，如果没有抓取过，则进行以上操作。
3. 如果总数达到了要求，则中断。

## 错误处理

1. 如果抓取发生错误，间隔errorWait重试。
2. 如果错误发生了errorTimes次，则退出程序。
3. 因为新浪微博API对抓取频率有所限制，因此如果出现`error_code 10023`, 不视为错误，间隔overloadWait时间后继续抓取。

## 怎样使用

### 1. 安装go语言

你必须首先安装了go语言。

### 2. 导入新浪微博的appkey以及cookies

appkey需要在新浪申请一个应用，cookies在浏览器里找到新浪微博的cookies。

	export APPKEY="your app key here"
	export COOKIES="your sina weibo cookies" 

### 3. 下载并运行代码

	git clone https://github.com/wb14123/fetch_weibo.git
	go run fetch_weibo.go

程序会把用户的timeline放入`timeline_go`文件夹，文件名为用户uid。把用户的friend关系放入`friend_go`文件夹，文件名也是用户uid。

