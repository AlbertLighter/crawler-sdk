package xhs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	webHost   = "https://www.xiaohongshu.com"
	apiHost   = "https://edith.xiaohongshu.com"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
)

// UserProfile 代表完整的用户个人信息，包括基本信息和笔记
type UserProfile struct {
	User  UserDetail `json:"user"`
	Notes []Note     `json:"notes"`
}

// UserDetail 代表从HTML中解析出的用户详细信息
type UserDetail map[string]interface{}

// Note 代表用户的单条笔记
type Note map[string]interface{}

// getUserDetail 获取小红书用户详细信息
func (c *Client) GetUserDetail(ctx context.Context, userID string) (UserDetail, error) {
	reqURL := fmt.Sprintf("%s/user/profile/%s", webHost, userID)
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("User-Agent", userAgent).
		Get(reqURL)

	if err != nil {
		return nil, fmt.Errorf("请求用户详情页失败: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("请求用户详情页返回错误状态: %s", resp.Status())
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	re := regexp.MustCompile(`window\.__INITIAL_STATE__=`)
	scriptContent := doc.Find("script").FilterFunction(func(i int, s *goquery.Selection) bool {
		return re.MatchString(s.Text())
	}).First().Text()

	if scriptContent == "" {
		return nil, fmt.Errorf("在HTML中未找到 __INITIAL_STATE__")
	}

	jsonStr := strings.TrimPrefix(scriptContent, "window.__INITIAL_STATE__=")
	jsonStr = strings.ReplaceAll(jsonStr, "undefined", "null")

	var data UserDetail
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return data, nil
}

// getUserNotes 获取用户发布的笔记列表
func (c *Client) GetUserNotes(ctx context.Context, userID string, offset, limit int) ([]Note, error) {
	var allNotes []Note
	var cursor string
	hasMore := true
	endLength := offset + limit

	for hasMore && len(allNotes) < endLength {
		apiURL := "/api/sns/web/v1/user_posted"
		params := map[string]interface{}{
			"num":           30,
			"cursor":        cursor,
			"user_id":       userID,
			"image_formats": []string{"jpg", "webp", "avif"},
		}

		fullURL := fmt.Sprintf("%s%s", apiHost, apiURL)
		resp, err := c.client.R().
			SetContext(ctx).
			SetQueryParamsFromValues(toURLValues(params)).
			Get(fullURL)

		if err != nil {
			return nil, fmt.Errorf("请求用户笔记API失败: %w", err)
		}
		if resp.IsError() {
			return nil, fmt.Errorf("请求用户笔记API返回错误状态: %s, body: %s", resp.Status(), resp.String())
		}

		var result struct {
			Data struct {
				Notes   []Note `json:"notes"`
				Cursor  string `json:"cursor"`
				HasMore bool   `json:"has_more"`
			} `json:"data"`
		}

		if err := json.Unmarshal([]byte(resp.String()), &result); err != nil {
			return nil, fmt.Errorf("解析笔记API响应失败: %w", err)
		}

		allNotes = append(allNotes, result.Data.Notes...)
		cursor = result.Data.Cursor
		hasMore = result.Data.HasMore
	}

	if len(allNotes) > offset+limit {
		return allNotes[offset : offset+limit], nil
	}
	if len(allNotes) > offset {
		return allNotes[offset:], nil
	}

	return []Note{}, nil
}

// getCookieValue 从cookie字符串中提取特定键的值
func getCookieValue(cookie, key string) string {
	parts := strings.Split(cookie, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, key+"=") {
			return strings.TrimPrefix(part, key+"=")
		}
	}
	return ""
}

// toURLValues 将 map[string]interface{} 转换为 url.Values
func toURLValues(params map[string]interface{}) url.Values {
	values := url.Values{}
	for k, v := range params {
		switch val := v.(type) {
		case string:
			values.Set(k, val)
		case []string:
			for _, item := range val {
				values.Add(k, item)
			}
		default:
			values.Set(k, fmt.Sprintf("%v", v))
		}
	}
	return values
}
