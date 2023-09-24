package client

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/yanmhlv/pcloud/pkg/util"
)

// uploadprogress
// downloadfile
// checksumfile

// DownloadFile; https://docs.pcloud.com/methods/file/downloadfile.html
func (c *pCloudClient) DownloadFile(urlStr string, path string, folderid int, target string) error {
	values := url.Values{
		"url":  {urlStr},
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case folderid >= 0:
		values.Add("folderid", strconv.Itoa(folderid))
	}

	if target != "" {
		values.Add("target", target)
	}

	return util.CheckResult(c.Client.Get(util.UrlBuilder("downloadfile", values)))
}

// UploadFile; https://docs.pcloud.com/methods/file/uploadfile.html
func (c *pCloudClient) UploadFile(reader io.Reader, path string, folderID int, filename string, noPartial int, progressHash string, renameIfExists int) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case folderID >= 0:
		values.Add("folderid", strconv.Itoa(folderID))
	default:
		return errors.New("bad params")
	}

	if filename == "" {
		return errors.New("bad params")
	}

	if noPartial > 0 {
		values.Add("nopartial", strconv.Itoa(noPartial))
	}
	if progressHash != "" {
		values.Add("progresshash", progressHash)
	}
	if renameIfExists > 0 {
		values.Add("renameifexists", strconv.Itoa(renameIfExists))
	}

	fw, err := w.CreateFormFile(filename, filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fw, reader); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", util.UrlBuilder("uploadfile", values), &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return util.CheckResult(c.Client.Do(req))
}

// CopyFile; https://docs.pcloud.com/methods/file/copyfile.html
func (c *pCloudClient) CopyFile(fileID int, path string, toFolderID int, toName string, toPath string) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case fileID > 0:
		values.Add("fileid", strconv.Itoa(fileID))
	case path != "":
		values.Add("path", path)
	default:
		return errors.New("bad params")
	}

	switch {
	case toFolderID > 0 && toName != "":
		values.Add("tofolderid", strconv.Itoa(toFolderID))
		values.Add("toname", toName)
	case toPath != "":
		values.Add("topath", toPath)
	default:
		return errors.New("bad params")
	}

	return util.CheckResult(c.Client.Get(util.UrlBuilder("copyfile", values)))
}

// DeleteFile; https://docs.pcloud.com/methods/file/deletefile.html
func (c *pCloudClient) DeleteFile(fileID int, path string) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case fileID > 0:
		values.Add("fileid", strconv.Itoa(fileID))
	case path != "":
		values.Add("path", path)
	default:
		return errors.New("bad params")
	}

	return util.CheckResult(c.Client.Get(util.UrlBuilder("deletefile", values)))
}

// RenameFile; https://docs.pcloud.com/methods/file/renamefile.html
func (c *pCloudClient) RenameFile(fileID int, path string, toPath string, toFolderID int, toName string) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case fileID > 0:
		values.Add("fileid", strconv.Itoa(fileID))
	case path != "":
		values.Add("path", path)
	default:
		return errors.New("bad params")
	}

	switch {
	case toPath != "":
		values["topath"] = []string{toPath}
	case toFolderID > 0 && toName != "":
		values["toname"] = []string{toName}
		values["tofolderid"] = []string{strconv.Itoa(toFolderID)}
	default:
		return errors.New("bad params")
	}

	return util.CheckResult(c.Client.Get(util.UrlBuilder("renamefile", values)))
}
