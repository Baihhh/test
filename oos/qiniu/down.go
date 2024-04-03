package qiniu

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/qiniu/go-sdk/v7/storage"
)

// 生成下载的url
func Url() {
	accessKey := "uCYp6h4RVW9F6fuiI1V8QSqLX3BZXHFO6uSQ_aM2"
	secretKey := "lj11G2izSIqp0dYM6yhiJQ891m-3vL3HQGvKvvLg"
	// bucket := "gulimall-baihhh"

	domain := "s9xgwfh7c.hb-bkt.clouddn.com"
	key := "testFile.txt"

	// 公开的
	// 就是做了一个拼接 domain + key
	publicAccessURL := storage.MakePublicURL(domain, key)
	fmt.Println(publicAccessURL)

	// 私有的需要先验证
	mac := auth.New(accessKey, secretKey)
	// token := mac.Sign([]byte(publicAccessURL))
	// fmt.Println(token)

	// putPolicy := storage.PutPolicy{
	// 	Scope: bucket,
	// }
	// token = putPolicy.UploadToken(mac)
	// fmt.Println(token)

	// 这里有点问题  一开始用这个方法的时候生成的url会报错  "error":"download token auth failed"
	// 但是当我用界面中生成的url访问一次之后  这里生成的就可以了
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)
	fmt.Println(privateAccessURL)
	// http://s9xgwfh7c.hb-bkt.clouddn.com/testFile.txt?e=1712110759&token=87t2Qe4CdVy23cyk57kDZEClOv3-xi6Gf6KMF0Pp:3m9AR6vERrKCKYVzv3VbddSdtWQ=
	// http://s9xgwfh7c.hb-bkt.clouddn.com/testFile.txt?e=1712114216&token=87t2Qe4CdVy23cyk57kDZEClOv3-xi6Gf6KMF0Pp:92Q2UcYf0cm2nfJzUKP6_9uwiwE=
}

func Down() {
	accessKey := "uCYp6h4RVW9F6fuiI1V8QSqLX3BZXHFO6uSQ_aM2"
	secretKey := "lj11G2izSIqp0dYM6yhiJQ891m-3vL3HQGvKvvLg"
	key := "testFile.txt"
	bucket := "gulimall-baihhh"

	mac := auth.New(accessKey, secretKey)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuabei
	// 是否使用https域名
	// cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	bm := storage.NewBucketManager(mac, &cfg)

	// err 和 resp 可能同时有值，当 err 有值时，下载是失败的，此时如果 resp 也有值可以通过 resp 获取响应状态码等其他信息
	resp, err := bm.Get(bucket, key, &storage.GetObjectInput{
		DownloadDomains: []string{
			// "bai", // 当前仅支持配置一个，不配置时，使用源站域名进行下载，会对下载的 URL 进行签名
		},
		PresignUrl: true,        // 下载 URL 是否进行签名，源站域名或者私有空间需要配置为 true
		Range:      "bytes=2-5", // 下载文件时 HTTP 请求的 Range 请求头
	})
	if err != nil || resp == nil {
		return
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

}

func getBM() *storage.BucketManager {
	accessKey := "uCYp6h4RVW9F6fuiI1V8QSqLX3BZXHFO6uSQ_aM2"
	secretKey := "lj11G2izSIqp0dYM6yhiJQ891m-3vL3HQGvKvvLg"
	// key := "testFile.txt"
	// bucket := "gulimall-baihhh"

	mac := auth.New(accessKey, secretKey)
	cfg := storage.Config{}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Region=&storage.ZoneHuabei
	// 是否使用https域名
	// cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	return storage.NewBucketManager(mac, &cfg)

}

func Info() {

	bucket := "gulimall-baihhh"
	bucketManager := getBM()

	// 文件相关的很多东西都在bucketManager 中包括信息获取、元信息修改、移动等
	// https://developer.qiniu.com/kodo/1238/go#rs-stat

	// bucketManager.ListFiles获取指定前缀的文件列表（迭代器）

	// 抓取网络资源到空间
	resURL := "http://devtools.qiniu.com/qiniu.png"
	// 指定保存的key
	fetchRet, err := bucketManager.Fetch(resURL, bucket, "qiniu.png")
	if err != nil {
		fmt.Println("fetch error,", err)
	} else {
		fmt.Println(fetchRet.String())
	}
}

func Batch() {
	bucket := "gulimall-baihhh"
	bucketManager := getBM()
	//每个batch的操作数量不可以超过1000个，如果总数量超过1000，需要分批发送
	keys := []string{
		"qiniu.png",
		"testFile.txt",
	}
	statOps := make([]string, 0, len(keys))
	for _, key := range keys {
		// 	拼接  批量获取文件信息
		// 这里是关键  storage.URIChangeMime  他就是修改文件信息
		statOps = append(statOps, storage.URIStat(bucket, key))
	}

	rets, err := bucketManager.Batch(statOps)
	if len(rets) == 0 {
		// 处理错误
		if e, ok := err.(*storage.ErrorInfo); ok {
			fmt.Printf("batch error, code:%d", e.Code)
		} else {
			fmt.Printf("batch error, %s", err)
		}
		return
	}

	// 返回 rets，先判断 rets 是否
	for _, ret := range rets {
		// 200 为成功
		if ret.Code == 200 {
			fmt.Printf("%d\n", ret.Code)
		} else {
			fmt.Printf("%s\n", ret.Data.Error)
		}
	}

	/**

	常见的MIME类型

	超文本标记语言文本 .html,.html text/html
	普通文本 .txt text/plain
	RTF文本 .rtf application/rtf
	GIF图形 .gif image/gif
	JPEG图形 .ipeg,.jpg image/jpeg
	au声音文件 .au audio/basic
	MIDI音乐文件 mid,.midi audio/midi,audio/x-midi
	RealAudio音乐文件 .ra, .ram audio/x-pn-realaudio
	MPEG文件 .mpg,.mpeg video/mpeg
	AVI文件 .avi video/x-msvideo
	GZIP文件 .gz application/x-gzip
	TAR文件 .tar application/x-tar
	*/
	// 每个MIME类型由两部分组成，前面是数据的大类别，例如声音audio、图象image等，后面定义具体的种类
	// 为什么是“text/HTML”而不是“HTML/text”或者别的什么？MIME Type 不是个人指定的，是经过 ietf 组织协商，以 RFC 的形式作为建议的标准发布在网上的，大多数的 Web 服务器和用户代理都会支持这个规范
	// storage.URIChangeMime()
	// 这里的type 是自定义的？
	// storage.URIChangeType()

}

func getOM() *storage.OperationManager {
	accessKey := "uCYp6h4RVW9F6fuiI1V8QSqLX3BZXHFO6uSQ_aM2"
	secretKey := "lj11G2izSIqp0dYM6yhiJQ891m-3vL3HQGvKvvLg"
	// key := "testFile.txt"
	// bucket := "gulimall-baihhh"

	mac := auth.New(accessKey, secretKey)
	cfg := storage.Config{}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Region=&storage.ZoneHuabei
	// 是否使用https域名
	// cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	return storage.NewOperationManager(mac, &cfg)
}

func Data() {
	operationManager := getOM()
	saveBucket := "gulimall-baihhh"
	bucket := "gulimall-baihhh"
	key := "qiniu.mp4" // 对视频进行处理
	// 处理指令集合
	// avthumb转码的结果会保存在原文件的空间中，但是文件名按照默认规则生成，为了方便获取转码后资源链接，建议自定义处理结果资源的名称，请参考处理结果另存 (saveas)。

	// 普通音视频转码接口方便用户对音频、视频资源进行编码和格式转换   /vb/<VideoBitRate>
	// https://developer.qiniu.com/dora/1248/audio-and-video-transcoding-avthumb
	fopAvthumb := fmt.Sprintf("avthumb/mp4/s/480x320/vb/500k|saveas/%s",
		storage.EncodedEntry(saveBucket, "pfop_test_qiniu.mp4"))
	// 视频单帧缩略图接口(vframe)用于从视频流中截取指定时刻的单帧画面并按指定大小缩放成图片
	// https://developer.qiniu.com/dora/1313/video-frame-thumbnails-vframe
	fopVframe := fmt.Sprintf("vframe/jpg/offset/10|saveas/%s",
		storage.EncodedEntry(saveBucket, "pfop_test_qiniu.jpg"))
	// 视频采样缩略图接口(vsample)用于从视频文件中截取多帧画面并按指定大小缩放成图片
	// https://developer.qiniu.com/dora/1315/video-sampling-thumbnails-vsample
	fopVsample := fmt.Sprintf("vsample/jpg/interval/20/pattern/%s",
		base64.URLEncoding.EncodeToString([]byte("pfop_test_$(count).jpg")))
	fopBatch := []string{fopAvthumb, fopVframe, fopVsample}
	fops := strings.Join(fopBatch, ";")
	// 强制重新执行数据处理任务
	force := true
	// 数据处理指令全部完成之后，通知该地址
	notifyURL := "http://api.example.com/pfop/callback"
	// 数据处理的私有队列，必须指定以保障处理速度
	pipeline := "jemy"
	persistentId, err := operationManager.Pfop(bucket, key, fops, pipeline, notifyURL, force)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(persistentId)
}

func getCM() *cdn.CdnManager {
	accessKey := "uCYp6h4RVW9F6fuiI1V8QSqLX3BZXHFO6uSQ_aM2"
	secretKey := "lj11G2izSIqp0dYM6yhiJQ891m-3vL3HQGvKvvLg"
	// key := "testFile.txt"
	// bucket := "gulimall-baihhh"

	mac := auth.New(accessKey, secretKey)
	return cdn.NewCdnManager(mac)
}

func CDN() {
	cdnManager := getCM()
	//刷新链接，单次请求链接不可以超过100个，如果超过，请分批发送请求
	urlsToRefresh := []string{
		"http://if-pbl.qiniudn.com/qiniu.png",
		"http://if-pbl.qiniudn.com/github.png",
	}
	ret, err := cdnManager.RefreshUrls(urlsToRefresh)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Code)
	fmt.Println(ret.RequestID)

	// 刷新目录，刷新目录需要联系七牛技术支持开通权限
	// 单次请求链接不可以超过10个，如果超过，请分批发送请求
	dirsToRefresh := []string{
		"http://if-pbl.qiniudn.com/images/",
		"http://if-pbl.qiniudn.com/static/",
	}
	ret, err = cdnManager.RefreshDirs(dirsToRefresh)

	//  获取域名流量
	domains := []string{
		"if-pbl.qiniudn.com",
		"qdisk.qiniudn.com",
	}
	startDate := "2023-07-30"
	endDate := "2024-07-31"
	granularity := "day"
	data, err := cdnManager.GetFluxData(startDate, endDate, granularity, domains)
	fmt.Printf("%v\n", data)
}
