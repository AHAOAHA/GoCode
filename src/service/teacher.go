/******************************************************************
 * Copyright(C) 2020-2020. All right reserved.
 *
 * Filename: teacher.go
 * Author: ahaoozhang
 * Date: 2020-03-09 22:35:18 (Monday)
 * Describe:
 ******************************************************************/
package service

import (
	"GradeManager/src/context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TeacherIndexHandler(c *gin.Context) {
	var t context.TeacherContext
	// check cookie
	if err := t.CheckCookies(c, "user_cookie"); err != nil {
		c.HTML(http.StatusBadRequest, "401.html", nil)
		return
	}

	// login status ok
	c.HTML(http.StatusOK, "teacher_index.html", gin.H{
		"loginer_name": t.Info.GetName(),
	})
}
