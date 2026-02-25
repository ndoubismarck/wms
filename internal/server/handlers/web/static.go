package web

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"server/internal/core/shared/utils/stringutil"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func (h *handler) serveStaticImage(context *gin.Context) {
	var width uint32
	var height uint32
	widthStr := context.Query("width")
	heightStr := context.Query("height")
	if len(widthStr) > 0 {
		width = stringutil.ToUint32(widthStr)
	}
	if len(heightStr) > 0 {
		height = stringutil.ToUint32(heightStr)
	}
	filePath := filepath.Join("./files/images", context.Param("any"))
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		context.Status(404)
		return
	}
	modifiedSince := strings.TrimSpace(context.GetHeader("If-Modified-Since"))
	if len(modifiedSince) > 0 {
		modifiedTime, err := time.Parse(http.TimeFormat, modifiedSince)
		if err == nil {
			if !fileInfo.ModTime().UTC().After(modifiedTime.UTC()) {
				context.Status(http.StatusNotModified)
				return
			}
		}
	}
	fileObj, err := os.Open(filePath)
	if err != nil {
		context.Status(404)
		h.ctx.Logger().Error(err)
		return
	}
	fileReader, err := io.ReadAll(fileObj)
	if err != nil {
		context.Status(404)
		h.ctx.Logger().Error(err)
		return
	}
	var outputBuff bytes.Buffer
	bytesReader := bytes.NewReader(fileReader)
	imgFile, imgFormat, err := image.Decode(bytesReader)
	if err != nil {
		context.Status(404)
		h.ctx.Logger().Error(err)
		return
	}
	var resizedImg *image.NRGBA
	if width == 0 && height == 0 {
		resizedImg = imaging.Fit(imgFile, 1000, 1000, imaging.Lanczos)
	} else {
		resizedImg = imaging.Resize(imgFile, int(width), int(height), imaging.Lanczos)
	}
	switch imgFormat {
	case "png":
		err = png.Encode(&outputBuff, resizedImg)
	case "jpeg":
		err = jpeg.Encode(&outputBuff, resizedImg, nil)
	default:
		context.Status(404)
		return
	}
	if err != nil {
		context.Status(404)
		h.ctx.Logger().Error(err)
		return
	}
	outputBytes := outputBuff.Bytes()
	context.Header("Cache-Control", "public, max-age=31536000, immutable")
	context.Header("Content-Length", fmt.Sprintf("%d", len(outputBytes)))
	context.Data(http.StatusOK, "image/"+imgFormat, outputBytes)
}
