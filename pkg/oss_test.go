package pkg

import "testing"

func TestOss_OssUpload(t *testing.T) {
	oss := &Oss{
		Endpoint:   "oss-cn-beijing.aliyuncs.com",
		AkId:       "xxx",
		AkSk:       "xxx",
		Bucket:     "cuisongliu",
		Object:     "tekton/release.yaml",
		UploadFile: "/home/cuisongliu/Downloads/release.yaml",
	}
	oss.OssUpload()
}
