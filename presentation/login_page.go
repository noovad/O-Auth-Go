package presentation

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginPage serves the HTML page
func LoginPage(c *gin.Context) {
	htmlIndex := `<html>
<head>
	<style>
		body {
			background-color: #ffffff;
			color: #1E1E1E;
			font-family: Arial, sans-serif;
		}
		.center {
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
		}
		.button {
			background-color: #1E1E1E;
			border: none;
			color: white;
			padding: 15px 32px;
			text-align: center;
			text-decoration: none;
			display: inline-block;
			font-size: 16px;
			margin: 4px 2px;
			cursor: pointer;
			border-radius: 12px;
			transition: background-color 0.3s ease;
		}
		.button:hover {
			background-color: #ffffff;
			border: 2px solid #1E1E1E;
			color: #1E1E1E;
		}
	</style>
</head>
<body>
	<div class="center">
		<a href="/login" class="button">Google Log In</a>
	</div>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlIndex))
}
