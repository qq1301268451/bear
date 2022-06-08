package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	params := &dto.FavoriteInput{}
	out := &dto.FavoriteOutput{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	_, exist := usersLoginInfo[params.Token]
	if !exist {
		out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	video := &dao.Video{}
	err := video.UpdateVideoByFavorite(params)
	if err != nil {
		out.ResponseError(common.SqlFindErr, common.SqlFindErrMsg)
	}

	out.ResponseSuccess()
	c.JSON(http.StatusOK, out)
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	params := &dto.FavoriteListInput{}
	out := &dto.FavoriteListOutput{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	_, exist := usersLoginInfo[params.Token]
	if !exist {
		out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	video := &dao.Video{}
	videoList, err := video.VideoListByFavorite(params)
	if err != nil {
		out.ResponseError(common.SqlFindErr, common.SqlFindErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}

	var outVideoList []dto.Video
	for _, item := range *videoList {
		isFollow := strings.Contains(item.User.FollowList, params.UserID)
		isFavorite := strings.Contains(item.FavoriteList, params.UserID)
		outVideoList = append(outVideoList, dto.Video{
			Id: int64(item.Model.ID),
			Author: dto.User{
				Id:            int64(item.User.Model.ID),
				Name:          item.User.Name,
				FollowCount:   item.User.FollowCount,
				FollowerCount: item.User.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: item.FavoriteCount,
			CommentCount:  item.CommentCount,
			IsFavorite:    isFavorite,
			Title:         item.Title,
		})
	}
	out.ResponseSuccess(outVideoList)
	c.JSON(http.StatusOK, out)

}
