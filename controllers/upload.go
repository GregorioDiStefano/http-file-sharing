package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
)

const (
	LOCAL = "local"
	S3    = "S3"
)

func SimpleUpload(c *gin.Context) {

	if _, err := checkUploadSize(c); err != nil {
		sendError(c, "Upload too large")
		return
	}

	fn := c.Param("filename")
	key := models.DB.FindUnsedKey()
	deleteKey := helpers.RandomString(helpers.Config.DeleteKeySize)

	if helpers.Config.StorageMethod == LOCAL {
		processUpload(c.Request.Body, key, fn)
	} else if helpers.Config.StorageMethod == S3 {
		if err := processUploadS3(c.Request.Body, key, fn); err != nil {
			sendError(c, "Uploading file to S3 bucket failed!")
			return
		}
	}

	simpleStoredFiled := models.StoredFile{
		Key:           key,
		DeleteKey:     deleteKey,
		FileName:      fn,
		FileSize:      c.Request.ContentLength,
		UploadTime:    time.Now().UTC(),
		LastAccess:    time.Now().UTC(),
		StorageMethod: helpers.Config.StorageMethod,
	}

	models.DB.WriteStoredFile(simpleStoredFiled)

	returnJSON := make(map[string]string)
	returnJSON["downloadURL"] = fmt.Sprintf("%s/%s/%s", helpers.Config.Domain, key, fn)
	returnJSON["deleteURL"] = fmt.Sprintf("%s/%s/%s/%s", helpers.Config.Domain, key, deleteKey, fn)

	c.JSON(http.StatusOK, returnJSON)
}
