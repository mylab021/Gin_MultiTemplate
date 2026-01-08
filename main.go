package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// 创建模板渲染器
func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// 前台页面
	r.AddFromFiles("frontend_index",
		"templates/frontend/base.html",
		"templates/frontend/header.html",
		"templates/frontend/index.html",
		"templates/frontend/footer.html",
	)

	r.AddFromFiles("frontend_about",
		"templates/frontend/base.html",
		"templates/frontend/header.html",
		"templates/frontend/about.html",
		"templates/frontend/footer.html",
	)

	// 后台页面
	r.AddFromFiles("admin_dashboard",
		"templates/admin/base.html",
		"templates/admin/sidebar.html",
		"templates/admin/dashboard.html",
	)

	r.AddFromFiles("admin_users",
		"templates/admin/base.html",
		"templates/admin/sidebar.html",
		"templates/admin/users.html",
		"templates/admin/footer.html",
	)

	// 认证页面
	r.AddFromFiles("auth_login",
		"templates/auth/base.html",
		"templates/auth/login.html",
	)

	r.AddFromFiles("auth_register",
		"templates/auth/base.html",
		"templates/auth/register.html",
	)

	// 认证页面
	r.AddFromFiles("auth_logout",
		"templates/auth/base.html",
		"templates/auth/logout.html",
	)
	return r
}

func main() {
	router := gin.Default()

	// 设置静态文件目录
	router.Static("/static", "./static")

	// 使用 multitemplate
	router.HTMLRender = createRender()

	// 前台路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "frontend_index", gin.H{
			"title": "首页",
			"data":  "欢迎来到我们的网站！",
		})
	})

	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "frontend_about", gin.H{
			"title": "关于我们",
		})
	})

	// 后台路由
	router.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_dashboard", gin.H{
			"title":       "仪表板",
			"username":    "管理员",
			"current":     "dashboard",
			"server_time": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	router.GET("/admin/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_users", gin.H{
			"title":    "用户管理",
			"username": "管理员",
			"current":  "users",
		})
	})

	// 登录页面
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth_login", gin.H{
			"title":    "用户登录",
			"redirect": c.Query("redirect"),
		})
	})

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		remember := c.PostForm("remember") == "on"
		redirect := c.PostForm("redirect")

		// 简单验证
		if username == "admin" && password == "123456" {
			if redirect != "" {
				c.Redirect(http.StatusFound, redirect)
			} else {
				c.Redirect(http.StatusFound, "/admin")
			}
			return
		}

		c.HTML(http.StatusOK, "auth_login", gin.H{
			"title":    "用户登录",
			"username": username,
			"remember": remember,
			"redirect": redirect,
			"error":    "用户名或密码错误",
		})
	})

	// 注册页面
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth_register", gin.H{
			"title": "用户注册",
		})
	})

	router.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")
		confirmPassword := c.PostForm("confirm_password")
		agreeTerms := c.PostForm("agree_terms") == "on"

		// 验证表单
		errors := make(map[string]string)

		if len(username) < 3 {
			errors["username"] = "用户名至少需要3个字符"
		}

		if len(password) < 6 {
			errors["password"] = "密码至少需要6个字符"
		}

		if password != confirmPassword {
			errors["confirm_password"] = "两次输入的密码不一致"
		}

		if !agreeTerms {
			errors["agree_terms"] = "请同意服务条款"
		}

		if len(errors) > 0 {
			c.HTML(http.StatusOK, "auth_register", gin.H{
				"title": "用户注册",
				"form": gin.H{
					"username":    username,
					"email":       email,
					"agree_terms": agreeTerms,
				},
				"errors": errors,
			})
			return
		}

		// 注册成功，跳转到登录页面
		c.Redirect(http.StatusFound, "/login?msg=注册成功，请登录")
	})

	router.GET("/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth_logout", gin.H{})
	})

	// 启动服务器
	router.Run(":8080")
}
