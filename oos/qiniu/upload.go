package qiniu

// 存储相关功能的引入包只有这两个，后面不再赘述
import (
	"context"
	"fmt"
	"os"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

// https://portal.qiniu.com/kodo/bucket/overview?bucketName=gulimall-baihhh

var domain = "s9xgwfh7c.hb-bkt.clouddn.com"

func NewAuth() string {
	accessKey := "uCYp6h4RVW9F6fuiI1V8QSqLX3BZXHFO6uSQ_aM2"
	secretKey := "lj11G2izSIqp0dYM6yhiJQ891m-3vL3HQGvKvvLg"
	bucket := "gulimall-baihhh"
	mac := auth.New(accessKey, secretKey)
	// mac = qbox.NewMac(accessKey, secretKey)

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	token := putPolicy.UploadToken(mac)
	fmt.Println(token)
	return token
}

func SecurityAuth() string {
	accessKey := "uCYp6h4RVW9F6fuiI1V8QSqLX3BZXHFO6uSQ_aM2"
	secretKey := "lj11G2izSIqp0dYM6yhiJQ891m-3vL3HQGvKvvLg"
	bucket := "gulimall-baihhh"
	mac := auth.New(accessKey, secretKey)
	// mac = qbox.NewMac(accessKey, secretKey)

	/**
	Key  为了防止恶意文件上传-上传策略中 scope 指定 key 值，只允许用户上传指定 key 的文件
	指定上传的目标资源空间 Bucket 和资源键 Key（最大为 750 字节）。有三种格式：
	<bucket>，表示允许用户上传文件到指定的 bucket。在这种格式下文件只能新增（分片上传需要指定insertOnly为1才是新增，否则也为覆盖上传），若已存在同名资源（且文件内容/etag不一致），上传会失败；若已存在资源的内容/etag一致，则上传会返回成功。
	<bucket>:<key>，表示只允许用户上传指定 key 的文件。在这种格式下文件默认允许修改，若已存在同名资源则会被覆盖。如果只希望上传指定 key 的文件，并且不允许修改，那么可以将下面的 insertOnly 属性值设为 1。
	<bucket>:<keyPrefix>，表示只允许用户上传指定以 keyPrefix 为前缀的文件，当且仅当 isPrefixalScope 字段为 1 时生效，isPrefixalScope 为 1 时无法覆盖上传。
	*/
	putPolicy := storage.PutPolicy{
		/**
		如果上传的时候key 和 这边指定的不同，回报这个错误PutFile 文件路径上传错误
		key doesn't match with scope
		*/
		Scope:        bucket + ":test.txt",
		ForceSaveKey: true,
		SaveKey:      "test.txt", // 上传之后的文件命名就是这个
		// 为了防止恶意文件上传-上传策略中 insertOnly 指定为 1 或者指定 forceInsertOnly 为 true，只允许用户上传指定 key ，并且不允许修改
		InsertOnly: 1, // 不能修改   默认是0可以修改
		/**
		限定用户上传的文件类型。指定本字段值，七牛服务器会侦测文件内容以判断 MimeType，再用判断值跟指定值进行匹配，匹配成功则允许上传，匹配失败则返回 403 状态码。示例：
		image/*表示只允许上传图片类型
		image/jpeg;image/png表示只允许上传jpg和png类型的图片
		!application/json;text/plain表示禁止上传json文本和纯文本ß。注意最前面的感叹号！
		*/
		MimeLimit: "image/*",
	}
	token := putPolicy.UploadToken(mac)
	fmt.Println(token)
	return token
}

// 基本上传
func Upload() {
	// token := NewAuth()
	token := SecurityAuth()

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuabei
	// 是否使用https域名
	// cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 可选配置
	putExtra := storage.PutExtra{
		// Params: map[string]string{
		// 	"x:name": "github logo",
		// },
	}
	// 根据路径上传文件
	localFile := "/Users/bai/Desktop/testFile.txt"
	key := "ii.txt"
	err := formUploader.PutFile(context.Background(), &ret, token, key, localFile, &putExtra)
	if err != nil {
		fmt.Println("PutFile 文件路径上传错误")
		fmt.Println(err)
		return
	}

	// // 通过子节数组上传
	// data := []byte("hello, this is qiniu cloud")
	// dataLen := int64(len(data))
	// keyByte := "testByte.txt"
	// err = formUploader.Put(context.Background(), &ret, token, keyByte, bytes.NewReader(data), dataLen, &putExtra)
	// if err != nil {
	// 	fmt.Println("Put 子节数组上传错误")
	// 	fmt.Println(err)
	// 	return
	// }

	// // 分片上传 v2通过将一个文件切割为块，然后通过上传块的方式来进行文件的上传
	// resumeUploader := storage.NewResumeUploaderV2(&cfg)
	// putExtraV2 := storage.RputV2Extra{}
	// keyV2 := "testV2.txt"
	// err = resumeUploader.PutFile(context.Background(), &ret, token, keyV2, localFile, &putExtraV2)
	// if err != nil {
	// 	fmt.Println("PutFile 分片上传错误")
	// 	fmt.Println(err)
	// 	return
	// }
	fmt.Println(ret.Key, ret.Hash)
}

// 断点上传   是基于分片上传实现的
func RecoderUpload() {
	token := NewAuth()
	localFile := "/Users/bai/Desktop/testFile.txt"
	key := "testFile.txt"

	cfg := storage.Config{}
	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	// PutRet 为七牛标准的上传回复内容。 如果使用了上传回调或者自定义了returnBody，那么需要根据实际情况，自己自定义一个返回值结构体
	ret := storage.PutRet{}
	recorder, err := storage.NewFileRecorder(os.TempDir())
	if err != nil {
		fmt.Println(err)
		return
	}
	putExtra := storage.RputV2Extra{
		Recorder: recorder,
	}
	err = resumeUploader.PutFile(context.Background(), &ret, token, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}

// 业务服务器验证存储服务回调

// 在上传策略里面设置了上传回调相关参数的时候，七牛在文件上传到服务器之后，会主动地向callbackUrl发送POST请求的回调，回调的内容为callbackBody模版所定义的内容，如果这个模版里面引用了魔法变量或者自定义变量，那么这些变量会被自动填充对应的值，然后在发送给业务服务器。

// 业务服务器在收到来自七牛的回调请求的时候，可以根据请求头部的Authorization字段来进行验证，查看该请求是否是来自七牛的未经篡改的请求。

// Go SDK 提供了一个方法 qbox.VerifyCallback 用于验证回调的请求：

// VerifyCallback 验证上传回调请求是否来自存储服务
// func VerifyCallback(mac *Mac, req *http.Request) (bool, error) {
// 	return mac.VerifyCallback(req)
// }
