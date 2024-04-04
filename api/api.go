package api

import (
	"github.com/skyrocketOoO/AuthNet/internal/delivery/rest"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Binding(r *gin.Engine, d *rest.Delivery) {
	r.GET("/ping", d.Ping)
	r.GET("/healthy", d.Healthy)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	relRouter := r.Group("/relation")
	{
		relRouter.GET("/", d.Get)
		relRouter.POST("/", d.Create)
		relRouter.DELETE("/", d.Delete)
		relRouter.DELETE("/all", d.ClearAll)

		relRouter.POST("/check", d.CheckAuth)
		relRouter.POST("/obj-auths", d.GetObjAuths)
		relRouter.POST("/sbj-who-has-auth", d.GetSbjsWhoHasAuth)
		relRouter.POST("/get-tree", d.GetTree)
		relRouter.GET("/see-tree", d.SeeTree)
	}
}
