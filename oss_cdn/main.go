package main

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	"oss_cdn/alioss"
	"oss_cdn/tencentcdn"
	"time"
)

const (
	//inke
	MEDIA_HOST = "https://media.meetstarlive.com"

	//oss
	OSSEndPoint        = ""
	OSSBucket          = ""
	OSSAccessKeyId     = ""
	OSSAccessKeySecret = ""

	//cdn
	CDNSecretId  = ""
	CDNSecureKey = ""
	CDNRegion    = regions.Guangzhou
)

func getQueryURL(objectKey string) string {
	return MEDIA_HOST + "/" + objectKey
}

func main() {
	keys := []string{"testdir/test1.txt", "testdir/test2.txt"}

	//OSS
	ossMgr, _ := alioss.NewOSSManager(OSSAccessKeyId, OSSAccessKeySecret, OSSEndPoint, OSSBucket)

	for _, file := range keys {
		ossMgr.PutObjectFromFile(file, file)
	}

	objs, _ := ossMgr.ListObjectsAll("testdir")
	_ = objs
	meta, _ := ossMgr.GetObject("testdir/test1.txt")
	_ = meta

	//CDN
	cdn, err := tencentcdn.NewCDNManager(CDNSecretId, CDNSecureKey, CDNRegion)
	if err != nil {
		panic(err)
	}

	//cdn.PurgeUrls([]string{getQueryURL(keys[0])})
	//cdn.PurgePaths([]string{"https://media.meetstarlive.com/testdir/"}, tencentcdn.FlushTypeFlush)

	//logs, _ := cdn.GetPurgeUrlsTasks(time.Now().Add(-30*time.Minute), "media.meetstarlive.com")
	//for _, log := range logs {
	//	fmt.Println(log)
	//}

	name1 := "media.meetstarlive.com"
	logs, _ := cdn.GetPurgePathsTasks(time.Now().Add(-30*time.Minute), name1)
	for _, log := range logs {
		fmt.Println(log)
	}

	cdn.DescribePurgeQuota()

}
