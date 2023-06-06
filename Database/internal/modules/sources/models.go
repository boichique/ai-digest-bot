package sources

type Source struct {
	UserID int    `json:"userID"`
	Source string `json:"source"`
}

type Video struct {
	Title       string `json:"title"`
	VideoID     string `json:"videoID"`
	PublishedAt string `json:"publishedAt"`
}

type SearchResult struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	Id   struct {
		Kind      string `json:"kind"`
		ChannelId string `json:"channelId,omitempty"`
		VideoId   string `json:"videoId,omitempty"`
	} `json:"id"`
	Snippet struct {
		PublishedAt string `json:"publishedAt"`
		ChannelId   string `json:"channelId"`
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
