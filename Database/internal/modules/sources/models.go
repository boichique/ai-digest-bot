package sources

type Source struct {
	UserID       int      `json:"userID"`
	Source       string   `json:"source"`
	NewVids      []string `json:"newVids"`
	TodaysDigest string   `json:"todaysDigest"`
}

type Video struct {
	Title       string `json:"title"`
	VideoID     string `json:"videoID"`
	PublishedAt string `json:"publishedAt"`
}

func (v Video) CompareTo(other interface{}) int {
	otherVideo, ok := other.(Video)
	if !ok {
		panic("Cannot compare Video with non-Video type")
	}

	if v.Title < otherVideo.Title {
		return -1
	} else if v.Title > otherVideo.Title {
		return 1
	} else {
		return 0
	}
}

type SearchResult struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	ID   struct {
		Kind      string `json:"kind"`
		ChannelID string `json:"channelId,omitempty"`
		VideoID   string `json:"videoId,omitempty"`
	} `json:"id"`
	Snippet struct {
		PublishedAt string `json:"publishedAt"`
		ChannelID   string `json:"channelId"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnails  map[string]struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumbnails"`
		ChannelTitle         string `json:"channelTitle"`
		LiveBroadcastContent string `json:"liveBroadcastContent"`
		PublishTime          string `json:"publishTime"`
	} `json:"snippet"`
}

type SearchListResponse struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []SearchResult `json:"items"`
}
