package repositories

import (
	"accurics-test_backup/httpclient"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"accurics-test/config"
	"accurics-test/entity"
	"encoding/base64"
)

const (
	createBranchURL   = "%s/repos/%s/%s/git/refs"
	contentURL        = "%s/repos/%s/%s/contents/%s"
	pullReqURL        = "%s/repos/%s/%s/pulls"
	shaReqURL         = "%s/repos/%s/%s/branches/%s"
	masterBranch      = "master"
	refConst          = "refs/heads/%s"
	authHeaderKey     = "Authorization"
	authToken         = "token %s"
	acceptKey         = "Accept"
	githubAcceptValue = "application/vnd.github.v3+json"
)

//CreateGithubImpl returns instance of CreateGithubImpl
func CreateGithubImpl() GithubImpl {
	return GithubImpl{}
}

//GithubImpl implements Github interface
type GithubImpl struct{}

//GetMasterSHA gets sha of master branch
func (g GithubImpl) GetMasterSHA(user, repo, token string) (string, error) {
	url := fmt.Sprintf(shaReqURL, config.ServerConfig.GithubAPIURL, user, repo, masterBranch)

	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(authToken, token)
	headers[acceptKey] = githubAcceptValue

	resp, err := httpclient.CallRestAPI(http.MethodGet, url, headers, nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		// config.AppLogger.ErrorLogger.Printf("response status for url %s is %s", url, resp.Status)
		return "", errors.New("file content not obtained")
	}

	apiOutput, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in reading response for url - %s, error - %s", url, err.Error())
		return "", err
	}

	shaResp := &entity.SHAResp{}
	err = json.Unmarshal(apiOutput, shaResp)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in unmarshaling response from url - %s, response - %s, error - %s", url, string(apiOutput), err.Error())
		return "", err
	}

	return shaResp.Commit.Sha, nil

}

//CreateBranch creates a branch
func (g GithubImpl) CreateBranch(user, repo, branchName, refShaID, token string) error {
	url := fmt.Sprintf(createBranchURL, config.ServerConfig.GithubAPIURL, user, repo)

	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(authToken, token)

	reqBody := entity.BranchReq{
		Ref: fmt.Sprintf(refConst, branchName),
		Sha: refShaID,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in marshaling request from url - %s, request - %s, error - %s", url, reqBody, err.Error())
		return err
	}
	resp, err := httpclient.CallRestAPI(http.MethodPost, url, headers, reqJSON)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in reading response for url - %s, error - %s", url, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		// config.AppLogger.ErrorLogger.Printf("response status for url %s is %s", url, resp.Status)
		return errors.New("branch not created")
	}
	return nil
}

//GetFileContent gets file content in plain text
func (g GithubImpl) GetFileContent(user, repo, branchName, fileName, token string) (*entity.FileContentResp, error) {
	url := fmt.Sprintf(contentURL, config.ServerConfig.GithubAPIURL, user, repo, fileName)

	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(authToken, token)
	resp, err := httpclient.CallRestAPI(http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		// config.AppLogger.ErrorLogger.Printf("response status for url %s is %s", url, resp.Status)
		return nil, errors.New("file content not obtained")
	}
	apiOutput, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in reading response for url - %s, error - %s", url, err.Error())
		return nil, err
	}
	fileContent := &entity.FileContentResp{}
	err = json.Unmarshal(apiOutput, fileContent)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in unmarshaling response from url - %s, response - %s, error - %s", url, string(apiOutput), err.Error())
		return nil, err
	}

	Content, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error - %s,  occured while decoding base64 file content - %s", err.Error(), fileContent.Content)
		return nil, err
	}
	fileContent.Content = string(Content)

	return fileContent, nil
}

//CreateORUpdateFile updates file
func (g GithubImpl) CreateORUpdateFile(user, repo, fileName, token string, req *entity.UpdateFileReq) error {
	url := fmt.Sprintf(contentURL, config.ServerConfig.GithubAPIURL, user, repo, fileName)

	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(authToken, token)
	headers[acceptKey] = githubAcceptValue

	reqJSON, err := json.Marshal(req)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in marshaling request from url - %s, request - %s, error - %s", url, req, err.Error())
		return err
	}

	resp, err := httpclient.CallRestAPI(http.MethodPut, url, headers, reqJSON)
	if err != nil {
		return err
	}

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusCreated) {
		// config.AppLogger.ErrorLogger.Printf("response status for url %s is %s", url, resp.Status)
		return err
	}

	return nil
}

//CreatePullRequest creates a pull request
func (g GithubImpl) CreatePullRequest(user, repo, token string, req entity.PullReq) error {
	url := fmt.Sprintf(pullReqURL, config.ServerConfig.GithubAPIURL, user, repo)

	headers := make(map[string]string)
	headers[authHeaderKey] = fmt.Sprintf(authToken, token)
	headers[acceptKey] = githubAcceptValue

	reqJSON, err := json.Marshal(req)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in marshaling request from url - %s, request - %s, error - %s", url, req, err.Error())
		return err
	}

	resp, err := httpclient.CallRestAPI(http.MethodPost, url, headers, reqJSON)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// config.AppLogger.ErrorLogger.Printf("error in reading response for url - %s, error - %s", url, err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// config.AppLogger.ErrorLogger.Printf("response status for url %s is %s", url, resp.Status)
		return errors.New("some erroneous status")
	}

	return nil
}
