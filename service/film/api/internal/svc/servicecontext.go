package svc

import (
	"encoding/json"
	"mall/service/film/api/internal/config"
	"mall/service/film/model"

	"github.com/golang-jwt/jwt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	FilmOrderModel     model.FilmOrderModel
	FilmOrderItemModel model.FilmOrderItemModel
	FilmPhotoModel     model.FilmPhotoModel
}

// 用于调试JWT令牌的函数
func ParseJwtToken(tokenString string, secret string) {
	// 移除Bearer前缀
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// 解析但不验证签名
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if token == nil {
		logx.Error("无法解析JWT令牌")
		return
	}

	// 检查签名方法
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		logx.Errorf("非预期的签名方法: %v", token.Header["alg"])
	}

	// 检查并打印claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		claimsJSON, _ := json.Marshal(claims)
		logx.Infof("JWT令牌内容: %s", string(claimsJSON))

		// 检查是否有uid
		if uid, exists := claims["uid"]; exists {
			logx.Infof("找到uid字段: %v (类型: %T)", uid, uid)
		} else {
			logx.Error("JWT令牌中缺少uid字段")
			// 尝试找到可能的用户ID字段
			for k, v := range claims {
				if k == "userId" || k == "user_id" || k == "userID" || k == "id" {
					logx.Infof("找到可能的用户ID字段: %s = %v", k, v)
				}
			}
		}

		// 验证是否过期
		if exp, ok := claims["exp"].(float64); ok {
			logx.Infof("令牌过期时间: %v", exp)
		}
	} else {
		logx.Error("无法获取JWT令牌中的claims")
	}
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// 打印服务上下文初始化信息
	logx.Infof("初始化FilmAPI服务上下文，MySQL配置: %s", c.Mysql.DataSource)
	logx.Infof("Auth配置: SecretKey=%s", c.Auth.AccessSecret)

	return &ServiceContext{
		Config:             c,
		FilmOrderModel:     model.NewFilmOrderModel(conn),
		FilmOrderItemModel: model.NewFilmOrderItemModel(conn),
		FilmPhotoModel:     model.NewFilmPhotoModel(conn),
	}
}
