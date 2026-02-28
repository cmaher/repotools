package github

type PRData struct {
	Number            int             `json:"number"`
	Title             string          `json:"title"`
	Body              string          `json:"body"`
	State             string          `json:"state"`
	Author            Author          `json:"author"`
	BaseRefName       string          `json:"baseRefName"`
	HeadRefName       string          `json:"headRefName"`
	HeadRefOid        string          `json:"headRefOid"`
	URL               string          `json:"url"`
	Labels            []Label         `json:"labels"`
	Assignees         []Author        `json:"assignees"`
	ReviewRequests    []ReviewRequest `json:"reviewRequests"`
	CreatedAt         string          `json:"createdAt"`
	UpdatedAt         string          `json:"updatedAt"`
	MergedAt          string          `json:"mergedAt"`
	ClosedAt          string          `json:"closedAt"`
	Additions         int             `json:"additions"`
	Deletions         int             `json:"deletions"`
	ChangedFiles      int             `json:"changedFiles"`
	Mergeable         string          `json:"mergeable"`
	ReviewDecision    string          `json:"reviewDecision"`
	IsDraft           bool            `json:"isDraft"`
	Comments          []Comment       `json:"comments"`
	Reviews           []Review        `json:"reviews"`
	Commits           []Commit        `json:"commits"`
	Files             []FileChange    `json:"files"`
	StatusCheckRollup []CheckStatus   `json:"statusCheckRollup"`
}

type Author struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Label struct {
	Name string `json:"name"`
}

type ReviewRequest struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Comment struct {
	Author    Author `json:"author"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
}

type Review struct {
	Author      Author `json:"author"`
	State       string `json:"state"`
	Body        string `json:"body"`
	SubmittedAt string `json:"submittedAt"`
	CreatedAt   string `json:"createdAt"`
}

type Commit struct {
	Oid             string   `json:"oid"`
	MessageHeadline string   `json:"messageHeadline"`
	Authors         []Author `json:"authors"`
}

type FileChange struct {
	Path      string `json:"path"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type CheckStatus struct {
	Name       string `json:"name"`
	Context    string `json:"context"`
	Conclusion string `json:"conclusion"`
	Status     string `json:"status"`
	State      string `json:"state"`
}

type ReviewComment struct {
	User         Author `json:"user"`
	Path         string `json:"path"`
	Line         *int   `json:"line"`
	OriginalLine *int   `json:"original_line"`
	CreatedAt    string `json:"created_at"`
	Body         string `json:"body"`
	DiffHunk     string `json:"diff_hunk"`
	InReplyToID  *int   `json:"in_reply_to_id"`
}
