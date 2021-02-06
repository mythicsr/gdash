//docs: https://cloud.tencent.com/document/product/228/37871
package tencentcdn

import (
	"errors"
	"fmt"
	cdnlib "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"sort"
	"strings"
	"time"
)

type cdnManager struct {
	secretId  string
	secureKey string
	region    string
	client    *cdnlib.Client
}

func NewCDNManager(secretId string, secureKey string, region string) (*cdnManager, error) {
	obj := cdnManager{
		secretId:  secretId,
		secureKey: secureKey,
		region:    region,
	}
	credential := common.NewCredential(secretId, secureKey)
	cpf := profile.NewClientProfile()
	client, err := cdnlib.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}
	obj.client = client
	return &obj, nil
}

//用于查询提交的 URL 刷新记录及执行进度
//
//keyword 支持域名过滤，或 http(s):// 开头完整 URL 过滤(如 media.meetstarlive.com 或 https://media.meetstarlive.com/testdir/test1.txt)
//
func (m *cdnManager) GetPurgeUrlsTasks(stime time.Time, keyword string) ([]*cdnlib.PurgeTask, error) {
	sTime := stime.Format("2006-01-02 15:04:05")
	eTime := time.Now().Format("2006-01-02 15:04:05")
	purgeType := "url"

	area := "mainland"
	request := cdnlib.NewDescribePurgeTasksRequest()
	request.Area = &area
	request.StartTime = &sTime
	request.EndTime = &eTime
	request.Keyword = &keyword
	request.PurgeType = &purgeType

	resp, err := m.client.DescribePurgeTasks(request)
	if err != nil {
		return nil, err
	}
	if !(resp != nil && resp.Response != nil && resp.Response.PurgeLogs != nil) {
		return nil, errors.New("some error")
	}

	sort.Slice(resp.Response.PurgeLogs, func(i, j int) bool {
		ri := resp.Response.PurgeLogs[i].CreateTime
		rj := resp.Response.PurgeLogs[j].CreateTime
		return *ri > *rj
	})

	return resp.Response.PurgeLogs, nil
}

//用于查询提交的目录刷新记录及执行进度
//
//keyword 支持域名过滤，或 http(s):// 开头完整 URL 过滤(如 media.meetstarlive.com 或 https://media.meetstarlive.com/testdir/ )
//
func (m *cdnManager) GetPurgePathsTasks(startTime time.Time, keyword string) ([]*cdnlib.PurgeTask, error) {
	sTime := startTime.Format("2006-01-02 15:04:05")
	eTime := time.Now().Format("2006-01-02 15:04:05")
	purgeType := "path"

	area := "mainland"
	request := cdnlib.NewDescribePurgeTasksRequest()
	request.Area = &area
	request.StartTime = &sTime
	request.EndTime = &eTime
	request.Keyword = &keyword
	request.PurgeType = &purgeType

	resp, err := m.client.DescribePurgeTasks(request)
	if err != nil {
		return nil, err
	}
	if !(resp != nil && resp.Response != nil && resp.Response.PurgeLogs != nil) {
		return nil, errors.New("some error")
	}

	sort.Slice(resp.Response.PurgeLogs, func(i, j int) bool {
		ri := resp.Response.PurgeLogs[i].CreateTime
		rj := resp.Response.PurgeLogs[j].CreateTime
		return *ri > *rj
	})
	return resp.Response.PurgeLogs, nil
}

// DescribePurgeQuota 用于查询账户刷新配额和每日可用量。
func (m *cdnManager) DescribePurgeQuota() (interface{}, error) {
	request := cdnlib.NewDescribePurgeQuotaRequest()
	resp, err := m.client.DescribePurgeQuota(request)
	if err != nil {
		return nil, err
	}
	if !(resp != nil && resp.Response != nil) {
		return nil, errors.New("some error")
	}
	return resp.Response, nil
}

//用于批量提交 URL 进行刷新，根据 URL 中域名的当前加速区域进行对应区域的刷新。 默认情况下境内、境外加速区域每日 URL 刷新额度各为 10000 条，每次最多可提交 1000 条。
//
//刷新任务提交时需要携带http://或https://协议标识。
//
//url 如 https://media.meetstarlive.com/testdir/test1.txt
func (m *cdnManager) PurgeUrls(urls []string) error {
	if len(urls) > 1000 {
		return fmt.Errorf("too many urls")
	}

	if err := hasHttpPrefix(urls); err != nil {
		return err
	}

	request := cdnlib.NewPurgeUrlsCacheRequest()
	for _, u := range urls {
		url := u
		request.Urls = append(request.Urls, &url)
	}

	_, err := m.client.PurgeUrlsCache(request)
	return err
}

// PurgePathCache 用于批量提交目录刷新，根据域名的加速区域进行对应区域的刷新。 默认情况下境内、境外加速区域每日目录刷新额度为各 100 条，每次最多可提交 20 条
//
//flushType刷新类型 ：flush刷新产生更新的资源，delete：刷新全部资源
//
//path 如 "https://media.meetstarlive.com/testdir/"
func (m *cdnManager) PurgePaths(paths []string, flushType string) error {
	if len(paths) > 20 {
		return fmt.Errorf("too many paths")
	}

	if err := hasHttpPrefix(paths); err != nil {
		return err
	}

	request := cdnlib.NewPurgePathCacheRequest()
	request.FlushType = &flushType

	for _, p := range paths {
		ph := p
		request.Paths = append(request.Paths, &ph)
	}

	_, err := m.client.PurgePathCache(request)
	return err
}

func hasHttpPrefix(urls []string) error {
	for _, u := range urls {
		if !strings.HasPrefix(u, "http") {
			return fmt.Errorf("'%s' has no http or https prefix", u)
		}
	}
	return nil
}
