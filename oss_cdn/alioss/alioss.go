//docs: https://help.aliyun.com/document_detail/31980.html?spm=a2c4g.11186623.6.1681.64af3ea7FJnNAh
package alioss

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"net/http"
)

type ossManager struct {
	accessKeyID     string
	accessKeySecret string
	endPoint        string
	bucketName      string
	client          *oss.Client
}

func NewOSSManager(accessKeyID, accessKeySecret, endPoint, bucketName string) (*ossManager, error) {
	mgr := ossManager{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		endPoint:        endPoint,
		bucketName:      bucketName,
		client:          nil,
	}
	client, err := oss.New(endPoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	mgr.client = client
	return &mgr, nil
}

func (this *ossManager) GetObject(objectKey string) (b []byte, err error) {
	client, err := oss.New(this.endPoint, this.accessKeyID, this.accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(this.bucketName)
	if err != nil {
		return nil, err
	}

	reader, err := bucket.GetObject(objectKey)
	if err != nil {
		return nil, err
	}

	b, err = ioutil.ReadAll(reader)
	return b, err
}

func (this *ossManager) PutObjectFromFile(objectKey string, file string) error {
	client, err := oss.New(this.endPoint, this.accessKeyID, this.accessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(this.bucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(objectKey, file)
	return err
}

//获取Last-Modified等信息
func (this *ossManager) GetObjectMeta(objectKey string) (meta http.Header, err error) {
	client, err := oss.New(this.endPoint, this.accessKeyID, this.accessKeySecret)
	if err != nil {
		return meta, err
	}

	bucket, err := client.Bucket(this.bucketName)
	if err != nil {
		return meta, err
	}

	meta, err = bucket.GetObjectMeta(objectKey)
	return meta, err
}

// cToken 游标，首次为空串
// 查询完毕 IsTruncated=false
// prefix 如 "testdir/"
func (this *ossManager) ListObjects(objectKeyPrefix string, cToken string) (result oss.ListObjectsResult, err error) {
	//todo careful
	if len(objectKeyPrefix) < 2 {
		return result, errors.New("objectKeyPrefix too short")
	}

	client, err := oss.New(this.endPoint, this.accessKeyID, this.accessKeySecret)
	if err != nil {
		return result, err
	}

	bucket, err := client.Bucket(this.bucketName)
	if err != nil {
		return result, err
	}

	result, err = bucket.ListObjects(oss.Prefix(objectKeyPrefix), oss.Marker(cToken), oss.MaxKeys(1000))
	return result, err
}

//遍历所有符合prefix的objectKey，大量数据慎用
func (this *ossManager) ListObjectsAll(objectKeyPrefix string) (objs []oss.ObjectProperties, err error) {
	objs = make([]oss.ObjectProperties, 0)
	nextCToken := ""

	//todo careful
	if len(objectKeyPrefix) < 3 {
		return objs, errors.New("objectKeyPrefix too short")
	}

	for {
		result, _ := this.ListObjects(objectKeyPrefix, nextCToken)
		objs = append(objs, result.Objects...)
		nextCToken = result.NextMarker

		if !result.IsTruncated {
			break
		}
	}

	return objs, nil
}
