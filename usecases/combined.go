package usecases

import (
	"accurics-test/entity"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

//Github defines interface for github related functions
type Github interface {
	GetMasterSHA(user, repo, token string) (string, error)
	CreateBranch(user, repo, branchName, refShaID, token string) error
	GetFileContent(user, repo, branchName, fileName, token string) (*entity.FileContentResp, error)
	CreateORUpdateFile(user, repo, fileName, token string, req *entity.UpdateFileReq) error
	CreatePullRequest(user, repo, token string, req entity.PullReq) error
}

//CreateCombinedUseCase creates instance of CombinedUseCase
func CreateCombinedUseCase(repo Github) CombinedUseCase {
	return CombinedUseCase{
		githubRepo: repo,
	}
}

//CombinedUseCase implements GithubAction interface in handlers
type CombinedUseCase struct {
	githubRepo Github
}

//Perform carries out functions required for a combined usecase
func (cuc CombinedUseCase) Perform(ctx context.Context, user, repo, branchName, token string) error {
	shaID, err := cuc.githubRepo.GetMasterSHA(user, repo, token)
	if err != nil {
		return err
	}

	err = cuc.githubRepo.CreateBranch(user, repo, branchName, shaID, token)
	if err != nil {
		return err
	}

	fileName, ok := ctx.Value("fileName").(string)
	if !ok {
		return errors.New("FileName not present")
	}

	fileContent, err := cuc.githubRepo.GetFileContent(user, repo, branchName, fileName, token)
	if err != nil {
		return err
	}

	updateReq, err := GetUpdateFileReq(fileContent, branchName)
	if err != nil {
		return err
	}

	err = cuc.githubRepo.CreateORUpdateFile(user, repo, fileName, token, updateReq)
	if err != nil {
		return err
	}

	req := GetPullReq(branchName, "master", "Pull Request", "Auto generated pull request")
	err = cuc.githubRepo.CreatePullRequest(user, repo, token, req)
	if err != nil {
		return err
	}

	return nil

}

//GetUpdateFileReq takes in fileContent and sha of a file and returns updated file content
func GetUpdateFileReq(fileContent *entity.FileContentResp, branchName string) (*entity.UpdateFileReq, error) {
	updateReq := &entity.UpdateFileReq{}

	content := fileContent.Content
	lineCount := strings.Count(content, "\n") + 2
	line := fmt.Sprintf("This is statement %d", lineCount)

	updateReq.Message = fmt.Sprintf("Adding Line %d", lineCount)
	updateReq.Content = base64.StdEncoding.EncodeToString([]byte(fileContent.Content + "\n" + line))
	updateReq.Branch = branchName
	updateReq.Sha = fileContent.Sha

	return updateReq, nil
}

//GetPullReq creates pull req json from head, base, body, title
func GetPullReq(head, base, body, title string) entity.PullReq {
	return entity.PullReq{
		Title: title,
		Head:  head,
		Base:  base,
		Body:  body,
	}
}
