package vkapi

import "net/url"

const (
	VkHost = "api.vk.com"
	//VkHost = "localhost:8081"
	VkScheme = "https"
	//VkScheme = "http"
)

const (
	WallPost                     = "wall.post"
	WallPostVer                  = "5.73"
	PhotosGetWallUploadServer    = "photos.getWallUploadServer"
	PhotosGetWallUploadServerVer = "5.124"
	PhotosSaveWallPhoto          = "photos.saveWallPhoto"
	PhotosSaveWallPhotoVer       = "5.124"
	StatsGet                     = "stats.get"
	StatsGetVer                  = "5.124"
)

// &stats_groups=visitors&interval=all&v=5.124
const (
	GroupID     = "group_id"
	AccessToken = "access_token"
	Photo       = "photo"
	Server      = "server"
	Hash        = "hash"
	StatsGroups = "stats_group"
	Interval    = "interval"
	Version     = "v"
)

// BuildURL format https://{VkHost}/method/{path}
func BuildURL(path string) string {
	u := url.URL{
		Host:   VkHost,
		Scheme: VkScheme,
		Path:   "method/" + path,
	}
	return u.String()
}
