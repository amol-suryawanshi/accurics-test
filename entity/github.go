package entity

//BranchReq repsents request body to create branch
type BranchReq struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

//FileContentResp represnts file Content Resp
type FileContentResp struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Content string `json:"content"`
	Sha string `json:"sha"`
}

//UpdateFileReq represents file updation request
type UpdateFileReq struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch string `json:"branch"`
	Sha string `json:"sha"`
}

//PullReq represents req body to create pull request
type PullReq struct{
	Title string `json:"title"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	Body  string `json:"body"`
}

//SHAResp represents SHA response
type SHAResp struct {
	Name string `json:"name"`
	Commit CommitType `json:"commit"`
}

type CommitType struct {
	Sha string `json:"sha"`
}