package awpapi

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"

	awpconf "github.com/kAbramenko/AnonWallPoster/internal/awp"
	"github.com/kAbramenko/AnonWallPoster/internal/vkapi"
	fh "github.com/valyala/fasthttp"
)

const (
	vkHost   = "api.vk.com"
	vkPath   = "method/wall.post"
	vkScheme = "https"
	vkVer    = "5.73"
)

var (
	templateGlobs = []string{
		"web/template/*.html",
		"web/template/*/*.html",
	}
)

func readTemplates() (*template.Template, error) {
	var err error
	tpl := template.New("spletnica")
	for _, tpath := range templateGlobs {
		tpl, err = tpl.ParseGlob(tpath)
		if err != nil {
			return nil, err
		}
	}

	return tpl, nil
}

type TemplateCtx struct {
	Stats    *vkapi.VkResponseStats
	Health   string
	Name     string
	Adequacy float64
}

// Index ...
func Index(ctx *fh.RequestCtx) {
	ctx.Logger().Printf("Request path %s", ctx.Request.URI().Path())
	ctx.SetContentType(mime.TypeByExtension(".html"))
	tpl, err := readTemplates()
	if err != nil {
		ctx.Error(err.Error(), fh.StatusInternalServerError)
	}
	stats, err := vkapi.GetStats()
	if err != nil {
		ctx.Logger().Printf("Cannot get statistic: %v", err)
	}
	tctx := TemplateCtx{
		Stats:    stats,
		Health:   awpconf.Cfg.Stats.Health,
		Adequacy: awpconf.Cfg.Stats.Adequacy,
		Name:     awpconf.Cfg.Name,
	}
	if err := tpl.ExecuteTemplate(ctx, "index", tctx); err != nil {
		ctx.Logger().Printf("Render error: %v", err)
		ctx.SetStatusCode(fh.StatusInternalServerError)
		return
	}
}

// Post ...
func Post(ctx *fh.RequestCtx) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Logger().Printf("Error: %v", err)
		ctx.SetStatusCode(fh.StatusBadRequest)
		return
	}
	msg, has := form.Value["msg"]
	if has == false && len(msg) == 0 {
		ctx.Logger().Printf("Not found key: 'msg'")
		ctx.SetStatusCode(fh.StatusBadRequest)
		return
	}
	ctx.Logger().Printf("Message: '%v'", msg)

	u := url.URL{
		Host:   vkHost,
		Scheme: vkScheme,
		Path:   vkPath,
	}
	log.Println("REQ", u.String())

	data := url.Values{
		"owner_id":     {fmt.Sprint("-", *awpconf.Cfg.GetGroup())},
		"from_group":   {*awpconf.Cfg.GetAsGroup()},
		"message":      msg,
		"access_token": {awpconf.Cfg.Token},
		"publish_date": {fmt.Sprintf("%d", time.Now().Add(time.Hour*24*90).Unix())},
		"v":            {vkVer},
	}
	photoList, has := form.File["photo"]
	photoListAttach := ""
	if has {
		photoArr := make([]string, 0)
		for _, photo := range photoList {
			// ctx.Logger().Printf("Has photo: %v", photo)
			ctx.Logger().Printf("Has photo: %v", photo.Filename)
			ctx.Logger().Printf("Has photo: %v", photo.Header)
			ctx.Logger().Printf("Has photo: %v", photo.Size)
			photoInfo, err := vkapi.UploadPhoto(photo)
			if err != nil {
				ctx.Logger().Printf("Cannot upload photo: %v", err)
			} else {
				ctx.Logger().Printf("Uploaded photo: %v", photoInfo.ID)
				photoArr = append(photoArr, fmt.Sprint("photo", photoInfo.OwnerID, "_", photoInfo.ID))
			}
		}
		if len(photoArr) != 0 {
			photoListAttach = strings.Join(photoArr, ",")
		}
	}
	if photoListAttach != "" {
		data.Add("attachments", photoListAttach)
	}
	log.Debug("VK Post:", data)
	resp, err := http.PostForm(u.String(), data)
	if err != nil {
		ctx.Logger().Printf("Warning: %v", err)
	}
	ans := json.NewDecoder(resp.Body)
	var vkResp vkapi.VkResponse
	if err := ans.Decode(&vkResp); err != nil {
		log.Warn("HTTP:", err)
		log.Debug(vkResp)
		ctx.SetStatusCode(fh.StatusInternalServerError)
		return
	}
	if vkResp.HasError() {
		log.Warn("Vk response:", vkResp.Error)
		ctx.SetStatusCode(fh.StatusServiceUnavailable)
	} else {
		log.Debug("Response WallPost:", vkResp)
	}
}

// BadRequest ..
func BadRequest(ctx *fh.RequestCtx) {
	ctx.Logger().Printf("ff")
	ctx.SetStatusCode(fh.StatusBadRequest)

}
