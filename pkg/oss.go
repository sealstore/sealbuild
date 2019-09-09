package pkg

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wonderivan/logger"
)

// 定义进度条监听器。
type OssProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		logger.Alert("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		logger.Alert("Transfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.\n",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		logger.Alert("Transfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		logger.Alert("Transfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

//对象存储结构
type Oss struct {
	Endpoint   string
	AkId       string
	AkSk       string
	Bucket     string
	Object     string
	UploadFile string
}

func (o *Oss) OssUpload() {
	client, err := oss.New(o.Endpoint, o.AkId, o.AkSk)
	if err != nil {
		logger.Error(err)
	}

	bucketName := o.Bucket
	objectName := o.Object
	localFilename := o.UploadFile

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		logger.Error(err)
	}

	err = bucket.PutObjectFromFile(objectName, localFilename, oss.Progress(&OssProgressListener{}))
	if err != nil {
		logger.Error(err)
	}
}
