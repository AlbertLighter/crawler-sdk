package xhs

import (
	"fmt"
	"os"

	"resty.dev/v3"
)

// 1. 获取上传许可
func (c *XhsClient) getUploadPermit(client *resty.Client) (*resty.Response, error) {
	fmt.Println("Step 1: Getting upload permit...")
	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"biz_name":   "spectrum",
			"scene":      "image", // 如果是视频，这里可能是 "video"
			"file_count": "1",
			"version":    "1",
			"source":     "web",
		}).
		SetHeaders(map[string]string{
			"Accept":        "application/json, text/plain, */*",
			"Authorization": "",
			"Origin":        "https://creator.xiaohongshu.com",
			"Referer":       "https://creator.xiaohongshu.com/",
			"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0",
		}).
		Get("https://creator.xiaohongshu.com/api/media/v1/upload/creator/permit")

	if err != nil {
		return nil, fmt.Errorf("failed to get upload permit: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("get upload permit request failed with status: %s", resp.Status())
	}

	fmt.Println("Successfully got upload permit.")
	// 在实际应用中，你需要解析 resp.Body() 来获取 uploadAddr, fileId 和 token
	// fmt.Println("Response Body:", resp.String())
	return resp, nil
}

// 2. 上传文件
func uploadFile(permitResp *resty.Response, filePath string) (*resty.Response, error) {
	fmt.Println("\nStep 2: Uploading file...")

	// 占位符：这里需要从 getUploadPermit 的响应中解析出真实的上传地址和凭证
	// 这是一个示例结构，你需要根据实际返回的 JSON 进行调整
	// 例如: uploadAddr := gjson.Get(permitResp.String(), "data.uploadTempPermits.0.uploadAddr").String()
	uploadURL := "https://ros-upload-d4.xhscdn.com/spectrum/wc7dm3ayWRWG0399QncgK1AAjBMdETEUhZReftWx49SAtrQ" // 从响应中解析
	uploadToken := "GODLoWDc60IQpItfzv8AOTnIBNM:eyJkZWFkbGluZSI6..."                                         // 从响应中解析

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(fileContent).
		SetHeaders(map[string]string{
			// 这个 Authorization 头是 PUT 请求特有的，格式和 API 请求的不同
			"Authorization": "q-sign-algorithm=sha1&q-ak=null&q-sign-time=...", // 从响应中解析或重新生成
			"Content-Type":  "image/png",                                       // 根据你的文件类型设置
			"Origin":        "https://creator.xiaohongshu.com",
			"Referer":       "https://creator.xiaohongshu.com/",
			"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0",
			// 这个 token 也是从上一步响应中获取的
			"x-cos-security-token": uploadToken,
		}).
		Put(uploadURL)

	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("upload file request failed with status: %s", resp.Status())
	}

	fmt.Println("Successfully uploaded file.")
	// 上传成功后，响应头里通常会包含文件的最终 URL 或 ID (如 Etag, X-Ros-Preview-Url)
	// fmt.Println("Upload Response Headers:", resp.Header())
	return resp, nil
}

// 3. 发布笔记 (函数和结构体已根据 69_Full.txt 的日志更新)

// PublishNotePayload 是根据抓包结果创建的发布笔记的完整请求体结构
type PublishNotePayload struct {
	Common    CommonInfo  `json:"common"`
	ImageInfo interface{} `json:"image_info"` // 根据日志为 null，具体结构未知
	VideoInfo VideoInfo   `json:"video_info"`
}

type CommonInfo struct {
	Type          string        `json:"type"`
	NoteID        string        `json:"note_id"`
	Source        string        `json:"source"`
	Title         string        `json:"title"`
	Desc          string        `json:"desc"`
	Ats           []string      `json:"ats"`
	HashTag       []string      `json:"hash_tag"`
	BusinessBinds string        `json:"business_binds"`
	PrivacyInfo   PrivacyInfo   `json:"privacy_info"`
	GoodsInfo     interface{}   `json:"goods_info"` // 根据日志为 {}
	BizRelations  []string      `json:"biz_relations"`
	CapaTraceInfo CapaTraceInfo `json:"capa_trace_info"`
}

type PrivacyInfo struct {
	OpType  int      `json:"op_type"`
	Type    int      `json:"type"`
	UserIDs []string `json:"user_ids"`
}

type CapaTraceInfo struct {
	ContextJson string `json:"contextJson"`
}

type VideoInfo struct {
	Fileid            string            `json:"fileid"`
	FileID            string            `json:"file_id"`
	FormatWidth       int               `json:"format_width"`
	FormatHeight      int               `json:"format_height"`
	VideoPreviewType  string            `json:"video_preview_type"`
	CompositeMetadata CompositeMetadata `json:"composite_metadata"`
	Timelines         []string          `json:"timelines"`
	Cover             CoverInfo         `json:"cover"`
	Chapters          []string          `json:"chapters"`
	ChapterSyncText   bool              `json:"chapter_sync_text"`
	Segments          SegmentsInfo      `json:"segments"`
}

type CompositeMetadata struct {
	Video MediaMetadata `json:"video"`
	Audio MediaMetadata `json:"audio"`
}

type MediaMetadata struct {
	Bitrate                 int    `json:"bitrate"`
	ColourPrimaries         string `json:"colour_primaries,omitempty"`
	Duration                int    `json:"duration"`
	Format                  string `json:"format"`
	FrameRate               int    `json:"frame_rate,omitempty"`
	Height                  int    `json:"height,omitempty"`
	MatrixCoefficients      string `json:"matrix_coefficients,omitempty"`
	Rotation                int    `json:"rotation,omitempty"`
	TransferCharacteristics string `json:"transfer_characteristics,omitempty"`
	Width                   int    `json:"width,omitempty"`
	Channels                int    `json:"channels,omitempty"`
	SamplingRate            int    `json:"sampling_rate,omitempty"`
}

type CoverInfo struct {
	Fileid        string       `json:"fileid"`
	FileID        string       `json:"file_id"`
	Height        int          `json:"height"`
	Width         int          `json:"width"`
	Frame         FrameInfo    `json:"frame"`
	Stickers      StickersInfo `json:"stickers"`
	Fonts         []string     `json:"fonts"`
	ExtraInfoJson string       `json:"extra_info_json"`
}

type FrameInfo struct {
	Ts           int  `json:"ts"`
	IsUserSelect bool `json:"is_user_select"`
	IsUpload     bool `json:"is_upload"`
}

type StickersInfo struct {
	Version int      `json:"version"`
	Neptune []string `json:"neptune"`
}

type SegmentsInfo struct {
	Count     int           `json:"count"`
	NeedSlice bool          `json:"need_slice"`
	Items     []SegmentItem `json:"items"`
}

type SegmentItem struct {
	Mute             int               `json:"mute"`
	Speed            int               `json:"speed"`
	Start            int               `json:"start"`
	Duration         float64           `json:"duration"`
	Transcoded       int               `json:"transcoded"`
	MediaSource      int               `json:"media_source"`
	OriginalMetadata CompositeMetadata `json:"original_metadata"`
}

// publishNote 函数根据提供的参数构建并发送发布笔记的请求。
// title: 笔记标题
// desc: 笔记描述
// videoFileID: 视频文件的 fileId (来自上传步骤)
// coverFileID: 封面图片的 fileId (来自上传步骤)
func publishNote(client *resty.Client, title, desc, videoFileID, coverFileID string) (*resty.Response, error) {
	fmt.Println("\nStep 3: Publishing note...")

	publishURL := "https://edith.xiaohongshu.com/web_api/sns/v2/note"

	// 根据 69_Full.txt 的日志文件内容构建请求体
	// 注意：许多字段（如 metadata）是硬编码的示例值，实际应用中你可能需要动态填充
	payload := PublishNotePayload{
		Common: CommonInfo{
			Type:          "video",
			NoteID:        "",
			Source:        "{\"type\":\"web\",\"ids\":\"\",\"extraInfo\":\"{\\\"subType\\\":\\\"official\\\",\\\"systemId\\\":\\\"web\\\"}\"}",
			Title:         title,
			Desc:          desc,
			Ats:           []string{},
			HashTag:       []string{},
			BusinessBinds: "{\"version\":1,\"noteId\":0,\"bizType\":0,\"noteOrderBind\":{},\"notePostTiming\":{},\"noteCollectionBind\":{\"id\":\"\"},\"noteSketchCollectionBind\":{\"id\":\"\"},\"coProduceBind\":{\"enable\":true},\"noteCopyBind\":{\"copyable\":true},\"interactionPermissionBind\":{\"commentPermission\":0},\"optionRelationList\":[]}",
			PrivacyInfo: PrivacyInfo{
				OpType:  1,
				Type:    0,
				UserIDs: []string{},
			},
			GoodsInfo:    map[string]interface{}{},
			BizRelations: []string{},
			CapaTraceInfo: CapaTraceInfo{
				ContextJson: "{\"longTextToImage\":{\"imageFileIds\":[]},\"recommend_title\":{\"recommend_title_id\":\"\",\"is_use\":3,\"used_index\":-1},\"recommendTitle\":[],\"recommend_topics\":{\"used\":[]}}",
			},
		},
		ImageInfo: nil,
		VideoInfo: VideoInfo{
			Fileid:           videoFileID,
			FileID:           videoFileID,
			FormatWidth:      480,
			FormatHeight:     854,
			VideoPreviewType: "full_vertical_screen",
			CompositeMetadata: CompositeMetadata{
				Video: MediaMetadata{Bitrate: 3718964, ColourPrimaries: "BT.709", Duration: 4567, Format: "HEVC", FrameRate: 30, Height: 854, MatrixCoefficients: "BT.709", Rotation: 0, TransferCharacteristics: "BT.709", Width: 480},
				Audio: MediaMetadata{Bitrate: 178919, Channels: 2, Duration: 4575, Format: "AAC", SamplingRate: 44100},
			},
			Timelines: []string{},
			Cover: CoverInfo{
				Fileid:        coverFileID,
				FileID:        coverFileID,
				Height:        854,
				Width:         480,
				Frame:         FrameInfo{Ts: 0, IsUserSelect: false, IsUpload: false},
				Stickers:      StickersInfo{Version: 2, Neptune: []string{}},
				Fonts:         []string{},
				ExtraInfoJson: "{}",
			},
			Chapters:        []string{},
			ChapterSyncText: false,
			Segments: SegmentsInfo{
				Count:     1,
				NeedSlice: false,
				Items: []SegmentItem{
					{
						Mute:        0,
						Speed:       1,
						Start:       0,
						Duration:    4.567,
						Transcoded:  0,
						MediaSource: 1,
						OriginalMetadata: CompositeMetadata{
							Video: MediaMetadata{Bitrate: 3718964, ColourPrimaries: "BT.709", Duration: 4567, Format: "HEVC", FrameRate: 30, Height: 854, MatrixCoefficients: "BT.709", Rotation: 0, TransferCharacteristics: "BT.709", Width: 480},
							Audio: MediaMetadata{Bitrate: 178919, Channels: 2, Duration: 4575, Format: "AAC", SamplingRate: 44100},
						},
					},
				},
			},
		},
	}

	resp, err := client.R().
		SetBody(payload). // 使用 SetBody 发送 JSON 数据
		SetHeaders(map[string]string{
			"Accept":        "application/json, text/plain, */*",
			"Authorization": authorization,
			"Cookie":        cookie,
			"Content-Type":  "application/json;charset=UTF-8",
			"Origin":        "https://creator.xiaohongshu.com",
			"Referer":       "https://creator.xiaohongshu.com/",
			"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0",
			"x-s":           xs,
			"x-t":           xt,
		}).
		Post(publishURL)

	if err != nil {
		return nil, fmt.Errorf("failed to publish note: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("publish note request failed with status: %s, body: %s", resp.Status(), resp.String())
	}

	fmt.Println("Successfully published note!")
	fmt.Println("Response:", resp.String())
	return resp, nil
}

func test() {
	// 创建一个 resty 客户端用于后续的 API 请求
	client := resty.New()

	// 1. 获取上传许可
	permitResp, err := getUploadPermit(client)
	if err != nil {
		panic(err)
	}
	// 在实际应用中，你需要从 permitResp 中解析出上传所需的信息
	// 例如: fileID := gjson.Get(permitResp.String(), "data.uploadTempPermits.0.fileIds.0").String()
	fmt.Println("Permit Response Status:", permitResp.Status())

	// 2. 上传文件
	// 占位符：将 "path/to/your/image.png" 替换为你的文件路径
	uploadResp, err := uploadFile(permitResp, "path/to/your/image.png")
	if err != nil {
		panic(err)
	}
	// 从 getUploadPermit 的响应中获取 fileId，而不是 uploadResp 的响应头
	// 这是一个占位符，你需要用真实逻辑替换
	uploadedFileID := "spectrum/wc7dm3ayWRWG0399QncgK1AAjBMdETEUhZReftWx49SAtrQ" // 从 getUploadPermit 响应中解析
	fmt.Println("Upload Response Status:", uploadResp.Status())
	fmt.Printf("Uploaded File ID: %s\n", uploadedFileID)

	// 3. 发布笔记
	// 占位符：请替换成你自己的笔记标题和描述
	noteTitle := "我的视频笔记标题"
	noteDesc := "这是通过 Go 脚本发布的视频笔记描述！ #Go #自动化"
	// 占位符：videoFileID 和 coverFileID 都来自第一步 getUploadPermit 的返回，你需要从该响应中解析
	// 通常，你会为视频和封面分别调用 getUploadPermit，或者在一次调用中获取两个许可
	videoFileID := "spectrum/ovWswFtPGlBalii4HplaQzFL_DsPYOO2Mr64dstTmR8EHUg" // 示例 video file ID
	coverFileID := "spectrum/wc7dm3ayWRWG0399QncgK1AAjBMdETEUhZReftWx49SAtrQ" // 示例 cover file ID

	publishResp, err := publishNote(client, noteTitle, noteDesc, videoFileID, coverFileID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Publish Response Status:", publishResp.Status())
}
