package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
	"strings"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	c.JSON(http.StatusOK, token)
	//actionType := c.Query("action_type")
	params := &dto.CommentActionInput{}
	out := &dto.CommentActionResponse{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	if user, err := CheckToken(params.Token); err == nil {
		//发布评论
		if params.ActionType == "1" {
			//text := c.Query("comment_text")
			//dao写入数据库
			com := &dao.Comment{}
			userObject, err := com.CommentAdd(int64(user.ID), params)
			if err != nil {
				out.ResponseError(common.SqlAddErr, common.SqlAddErrMsg)
				c.JSON(http.StatusOK, out)
				return
			}
			userIdString := fmt.Sprintf("%s#", strconv.Itoa(int(user.Model.ID)))
			isFollow := strings.Contains(userObject.FollowerList, userIdString)
			//fmt.Println(isFollow)
			timeLayout := "01-02"
			createTime := com.CreatedAt.Format(timeLayout)
			fmt.Println("CommentAction")
			fmt.Println(int64(userObject.ID))
			fmt.Println(userObject.Name)
			//dto返回json
			out = &dto.CommentActionResponse{
				Response: dto.Response{
					StatusCode: 0,
					StatusMsg:  "Msg sent",
				},

				Comment: dto.Comment{
					Id: int64(com.ID),
					User: dto.User{
						Id:            int64(user.ID),
						Name:          user.Name,
						FollowCount:   user.FollowCount,
						FollowerCount: user.FollowerCount,
						IsFollow:      isFollow,
					},
					Content:    com.CommentText,
					CreateDate: createTime,
				},
			}
		}
		//删除评论
		if params.ActionType == "2" {
			//删除数据库中的某些值
			com := &dao.Comment{}
			userObject, err := com.CommentDelte(int64(user.ID), params)
			userIdString := fmt.Sprintf("%s#", strconv.Itoa(int(user.Model.ID)))
			isFollow := strings.Contains(userObject.FollowerList, userIdString)
			timeLayout := "01-02"
			createTime := com.CreatedAt.Format(timeLayout)
			if err != nil {
				out.ResponseError(common.SqlDelErr, common.SqlDelErrMsg)
				c.JSON(http.StatusOK, out)
				return
			}

			out = &dto.CommentActionResponse{
				Response: dto.Response{
					StatusCode: 0,
					StatusMsg:  "Msg deleted",
				},
				Comment: dto.Comment{
					Id: 0,
					User: dto.User{
						Id:            int64(userObject.ID),
						Name:          user.Name,
						FollowCount:   user.FollowCount,
						FollowerCount: user.FollowerCount,
						IsFollow:      isFollow,
					},
					Content:    com.CommentText,
					CreateDate: createTime,
				},
			}
		}
		c.JSON(http.StatusOK, out)
		return
	}
	c.JSON(http.StatusOK, out)
	return

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	//获取参数，校验
	token := c.Query("token")
	videoid := c.Query("video_id")
	videoidInt, _ := strconv.Atoi(videoid)
	out := &dto.CommentListResponse{}

	if user, err := CheckToken(token); err == nil {

		var com dao.Comment
		var params dto.CommentListRequire
		com.VideoId = int64(videoidInt)
		fmt.Println("com.VideoId:", com.VideoId)

		//检查video_id
		if _, err := com.VideoIdCheck(com); err == nil {
			//text := c.Query("comment_text")
			params.VideoId = com.VideoId
			commentlist, userObject, err := com.VideoCommentList(int64(user.ID), params)
			if err != nil && commentlist == nil {
				out.ResponseError(common.SqlFindErr, common.SqlFindErrMsg)
				c.JSON(http.StatusOK, out)
				return
			}
			if err != nil && commentlist != nil {
				out.ResponseError(common.SqlFindErr, common.SqlFindErrMsg)
				c.JSON(http.StatusOK, out)
				return

			}

			userIdString := fmt.Sprintf("%s#", strconv.Itoa(int(user.Model.ID)))
			isFollow := strings.Contains(userObject.FollowerList, userIdString)

			var outcommentlist []dto.Comment
			for _, item := range *commentlist {
				//timeLayout := "2006-01-02 15:04:05"
				timeLayout := "01-02"
				createtime := item.Model.CreatedAt.Format(timeLayout)
				outcommentlist = append(outcommentlist,
					dto.Comment{
						Id: int64(item.Model.ID),
						User: dto.User{
							Id:            int64(item.User.Model.ID),
							Name:          item.User.Name,
							FollowCount:   item.User.FollowCount,
							FollowerCount: item.User.FollowerCount,
							IsFollow:      isFollow, //没有关注列表，待完善

						},
						Content:    item.CommentText,
						CreateDate: createtime,
					})
			}
			out.ResponseSuccess(&outcommentlist)
			c.JSON(http.StatusOK, out)
			//c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			//	Comment: Comment{
			//		Id: 1,
			//		//User:       "1",
			//		Content:    "1234",
			//		CreateDate: "05-01",
			//	}})

			return
		}
		//_, err := com.VideoIdCheck(com)
		//fmt.Println("err", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video doesn't exist"})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	//c.JSON(http.StatusOK, CommentListResponse{
	//	Response:    Response{StatusCode: 0},
	//	CommentList: DemoComments,
	//})
}
