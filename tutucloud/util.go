package tutucloud

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
)

// 参数签名
func signature(params map[string]string, api_secret string) string {
	// 排序
	keys := []string{}
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//拼接 key value, 转换 key 为小写
	signstr := ""
	for _, k := range keys {
		signstr += strings.ToLower(k) + params[k]
	}
	//拼接私有 key
	signstr += api_secret

	//返回 md5 字符串
	hash := md5.New()
	hash.Write([]byte(signstr))
	return hex.EncodeToString(hash.Sum(nil))
}
