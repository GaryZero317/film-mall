package logic

import (
	"context"
	"fmt"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkLogic {
	return &CreateWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWorkLogic) CreateWork(req *types.CreateWorkReq) (resp *types.CreateWorkResp, err error) {
	// 初始化默认响应（失败状态）
	resp = &types.CreateWorkResp{
		Code: 500,
		Msg:  "创建失败",
		Data: types.CreateWorkData{},
	}

	// 获取用户ID
	uid, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok || uid <= 0 {
		resp.Code = 401
		resp.Msg = "未登录或登录已过期"
		return resp, nil
	}

	logx.Infof("尝试获取用户信息，用户ID: %d", uid)

	// 查询用户信息获取用户名
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uid)
	if err != nil {
		logx.Errorf("查询用户信息失败: %v, 用户ID: %d", err, uid)
		// 检查数据库中是否真的有该用户
		logx.Infof("直接检查数据库中用户ID为%d的记录", uid)
		var checkUser struct {
			Id   int64  `db:"id"`
			Name string `db:"name"`
		}
		// 获取数据库连接并直接执行查询
		conn := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource)
		checkErr := conn.QueryRowCtx(
			l.ctx,
			&checkUser,
			"SELECT id, name FROM user WHERE id = ? LIMIT 1",
			uid,
		)
		if checkErr != nil {
			logx.Errorf("直接查询用户失败: %v", checkErr)
		} else {
			logx.Infof("直接查询用户成功: ID=%d, Name=%s", checkUser.Id, checkUser.Name)
			// 既然直接查询成功了，让我们使用这个结果
			if checkUser.Name != "" {
				user = &model.User{
					Id:       checkUser.Id,
					Nickname: checkUser.Name,
				}
				logx.Infof("使用直接查询获取的用户名称: %s", user.Nickname)
			}
		}
		// 即使查询失败也继续创建作品，只是没有作者名
	} else {
		logx.Infof("成功获取用户信息: %+v", user)
	}

	// 获取用户昵称，如果查询失败则使用默认名称
	var nickname string
	if user != nil && user.Nickname != "" {
		nickname = user.Nickname
		logx.Infof("使用查询到的用户昵称: %s", nickname)
	} else {
		nickname = fmt.Sprintf("用户%d", uid)
		logx.Infof("使用默认用户昵称: %s", nickname)
	}

	// 创建作品
	name := nickname // 使用变量存储昵称，这样可以获取指针
	work := &model.Work{
		Uid:          uid,
		Name:         &name,
		Title:        req.Title,
		Description:  req.Description,
		CoverUrl:     req.CoverUrl,
		FilmType:     req.FilmType,
		FilmBrand:    req.FilmBrand,
		Camera:       req.Camera,
		Lens:         req.Lens,
		ExifInfo:     req.ExifInfo,
		ViewCount:    0,
		LikeCount:    0,
		CommentCount: 0,
		Status:       req.Status,
	}

	// 执行数据库插入操作
	logx.Infof("开始插入作品: %+v", work)
	result, err := l.svcCtx.WorkModel.Insert(l.ctx, work)
	if err != nil {
		resp.Code = 500
		resp.Msg = "创建作品失败: " + err.Error()
		logx.Errorf("插入作品失败: %v", err)
		return resp, nil
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		resp.Code = 500
		resp.Msg = "获取作品ID失败: " + err.Error()
		logx.Errorf("获取插入ID失败: %v", err)
		return resp, nil
	}

	// 验证插入是否成功
	logx.Infof("作品创建成功，ID: %d", insertId)

	// 设置成功响应
	resp.Code = 0
	resp.Msg = "创建成功"
	resp.Data = types.CreateWorkData{
		Id: insertId,
	}

	// 作为额外验证，尝试查询刚刚创建的作品
	createdWork, findErr := l.svcCtx.WorkModel.FindOne(l.ctx, insertId)
	if findErr != nil {
		logx.Errorf("创建后无法查询到作品(ID=%d): %v", insertId, findErr)
	} else {
		logx.Infof("验证创建: 成功查询到新作品: ID=%d, 标题=%s", createdWork.Id, createdWork.Title)
	}

	return resp, nil
}
