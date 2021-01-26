package vkapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	awpconf "github.com/kAbramenko/AnonWallPoster/internal/awp"
)

func GetStats() (*VkResponseStats, error) {
	data := url.Values{
		GroupID:     {*awpconf.Cfg.GetGroup()},
		AccessToken: {awpconf.Cfg.Token},
		StatsGroups: {"visitor"},
		Interval:    {"all"},
		Version:     {StatsGetVer},
	}
	resp, err := http.PostForm(BuildURL(StatsGet), data)
	if err != nil {
		return nil, err
	}
	vkResp := VkResponseStats{}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&vkResp); err != nil {
		return nil, err
	}
	return &vkResp, nil
}

// UploadPhoto ...
func UploadPhoto(part *multipart.FileHeader) (*VkPhotoObject, error) {
	vkResp, err := photosGetWallUploadServer()
	if err != nil {
		return nil, err
	}
	if vkResp.HasError() {
		return nil, errors.New(vkResp.Error.ErrorMsg)
	}
	uploadResp, err := uploadPhotoToServer(vkResp.Response.UploadURL, part)
	if err != nil {
		return nil, err
	}
	if vkResp.HasError() {
		return nil, errors.New(vkResp.Error.ErrorMsg)
	}

	vkRespArr, err := photosSaveWallPhoto(uploadResp)
	if err != nil {
		return nil, err
	}
	if len(vkRespArr.Response) == 0 {
		return nil, errors.New("Response is empty")
	}

	return &(vkRespArr.Response[0]), nil
}

func photosGetWallUploadServer() (*VkResponse, error) {
	data := url.Values{
		GroupID:     {*awpconf.Cfg.GetGroup()},
		AccessToken: {*&awpconf.Cfg.Token},
		Version:     {PhotosGetWallUploadServerVer},
	}
	resp, err := http.PostForm(BuildURL(PhotosGetWallUploadServer), data)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var vkResp VkResponse
	err = decoder.Decode(&vkResp)
	if err != nil {
		return nil, err
	}
	return &vkResp, nil
}

func photosSaveWallPhoto(vkParam *VkUploadResponse) (*VkResponseArray, error) {
	data := url.Values{
		GroupID:     {*awpconf.Cfg.GetGroup()},
		AccessToken: {awpconf.Cfg.Token},
		Photo:       {vkParam.Photo},
		Hash:        {vkParam.Hash},
		Server:      {fmt.Sprint(vkParam.Server)},
		Version:     {PhotosSaveWallPhotoVer},
	}
	resp, err := http.PostForm(BuildURL(PhotosSaveWallPhoto), data)
	if err != nil {
		return nil, err
	}
	var vkResp VkResponseArray
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&vkResp); err != nil {
		return nil, err
	}

	return &vkResp, nil
}

func uploadPhotoToServer(url string, f *multipart.FileHeader) (*VkUploadResponse, error) {
	rf, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rf.Close()

	var buff bytes.Buffer
	data := multipart.NewWriter(&buff)
	wf, err := data.CreatePart(f.Header)
	if err != nil {
		return nil, err
	}
	io.Copy(wf, rf)
	if err := data.Close(); err != nil {
		return nil, err
	}

	resp, err := http.Post(url, data.FormDataContentType(), &buff)
	if err != nil {
		return nil, err
	}
	var vkResp VkUploadResponse

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&vkResp); err != nil {
		return nil, err
	}

	return &vkResp, nil
}

// Fun ...
func Fun() {
	resp, err := photosGetWallUploadServer()
	fmt.Println(resp, err)
}
