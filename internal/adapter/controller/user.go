package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/DESA/internal/adapter/validator"
	"github.com/wang900115/DESA/internal/application/usecase"
	"github.com/wang900115/DESA/lib/common"
	iresponse "github.com/wang900115/DESA/lib/common/response"
	"github.com/wang900115/DESA/lib/domain"
)

type UserController struct {
	user usecase.UserUsecase
	resp iresponse.IResponse
	p2p  usecase.P2PUsecase
}

func NewUserController(user *usecase.UserUsecase, p2p *usecase.P2PUsecase, resp iresponse.IResponse) *UserController {
	return &UserController{user: *user, p2p: *p2p, resp: resp}
}

// 登入
func (u *UserController) Login(c *gin.Context) {
	var req validator.LoginAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	accessToken, refreshToken, user, err := u.user.Login(c, req.Username, req.Password)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	u.resp.SuccessWithData(c, common.LOGIN_SUCESS, map[string]interface{}{
		"aaccessToken": accessToken,
		"refreshToken": refreshToken,
		"user":         user})
}

// 登出
func (u *UserController) Logout(c *gin.Context) {
	userPeerID := c.GetString("user_id")
	if err := u.user.Logout(c, userPeerID); err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	u.resp.Success(c, common.LOGOUT_SUCCESS)
}

// 註冊用戶節點
func (u *UserController) Register(c *gin.Context) {
	var req validator.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	user := domain.User{
		Username:    req.Username,
		Password:    &req.Password,
		FirstEmail:  req.FirstEmail,
		SecondEmail: req.SecondEmail,
		Phone:       req.Phone,
		NickName:    req.NickName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Birth:       req.Birth,
		Country:     req.Country,
		City:        req.City,
	}

	peerID, multiAddr, pri, pub, err := u.p2p.CreateUserHost(user)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	user.PeerID = peerID
	user.MultiAddr = multiAddr
	user.PublicKey = &pub

	created, err := u.user.CreateUser(c, user)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	u.resp.SuccessWithData(c, common.CREATE_SUCCESS, map[string]interface{}{
		"user":       created,
		"privateKey": pri,
		"publicKey":  pub,
	})
}

// 獲取所有用戶資訊
func (u *UserController) Query(c *gin.Context) {
	users, err := u.user.QueryUser(c)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	u.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"users": users,
	})
}

// 更新用戶資訊(不可更新的有:PeerID, MultiAddr, PublicKey)
func (u *UserController) Update(c *gin.Context) {
	var req validator.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	peerID := c.GetString("user_id")
	user := domain.User{
		PeerID:   peerID,
		Username: req.Username,
		Password: &req.Password,

		FirstEmail:  req.FirstEmail,
		SecondEmail: req.SecondEmail,
		Phone:       req.Phone,
		NickName:    req.NickName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Birth:       req.Birth,
		Country:     req.Country,
		City:        req.City,
	}

	updated, err := u.user.UpdateUser(c, user)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	u.resp.SuccessWithData(c, common.UPDATE_SUCCESS, map[string]interface{}{
		"user": updated,
	})
}

// 刪除用戶節點
func (u *UserController) Delete(c *gin.Context) {
	var req validator.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	userPeerID := c.GetString("user_id")
	if err := u.user.DeleteUser(c, userPeerID); err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	if err := u.p2p.ShutDownUserHost(userPeerID); err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	u.resp.Success(c, common.DELETE_SUCCESS)
}
