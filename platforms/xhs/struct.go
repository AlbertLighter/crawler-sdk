package xhs

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
