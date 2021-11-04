package route


import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 跨域处理中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func InitRouter(engine *gin.Engine, fs embed.FS, flag bool) *gin.Engine {

	engine.Use(Cors())

	//resourceRouter(engine)

	group := engine.Group("v1")
	{
		//groupRouter(group)     // 服务分组接口
		//customizeRouter(group) // 自定义服务接口
		//websocketRouter(group)
	}
	return engine

}

// resourceRouter 静态资源配置
func resourceRouter(engine *gin.Engine) {
	html := NewHtmlHandler()
	group := engine.Group("/ui")
	{
		group.GET("", html.Index)
	}
	// 解决刷新404问题
	engine.NoRoute(html.RedirectIndex)
}

func groupRouter(engine *gin.RouterGroup) {
	api := api.NewGroupApi()
	group := engine.Group("group")
	{
		group.GET("pages", api.Pages)
		group.GET("delete", api.Delete)
		group.POST("insert", api.Create)
	}
}


// customizeRouter 自定义服务相关
func customizeRouter(engine *gin.RouterGroup) {
	customize := api.NewCustomizeApi()
	group := engine.Group("customize")
	{
		group.POST("list", customize.ServerList)
		group.GET("running", customize.Run)
		group.POST("add", customize.Create)
		group.GET("delete", customize.Delete)
	}
}


func websocketRouter(engine *gin.RouterGroup) {
	ws := api.NewWebSocketApi()
	group := engine.Group("ws")
	{
		group.GET("", ws.RealTimeLog)  // 读取实时日志
	}
}