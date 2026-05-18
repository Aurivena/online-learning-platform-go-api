package gateway

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"online-learning-platform-go-api/internal/course/dto"
	"online-learning-platform-go-api/internal/course/entity"
)

func objectKeyFromPayload(p entity.PayloadJSON) string {
	if p == nil {
		return ""
	}
	v, ok := p["object_key"].(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(v)
}

func clonePayload(p entity.PayloadJSON) entity.PayloadJSON {
	if p == nil {
		return entity.PayloadJSON{}
	}
	out := entity.PayloadJSON{}
	for k, v := range p {
		out[k] = v
	}
	return out
}

func overlayPayload(dst, src entity.PayloadJSON) {
	if dst == nil || src == nil {
		return
	}
	for k, v := range src {
		dst[k] = v
	}
}

func payloadHasFileReference(p entity.PayloadJSON) bool {
	if p == nil {
		return false
	}
	keys := []string{"object_key", "objectKey", "file_src", "fileSrc", "url"}
	for i := range keys {
		v, ok := p[keys[i]]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if ok && strings.TrimSpace(s) != "" {
			return true
		}
	}
	return false
}

func (g *CourseGateway) publicObjectURL(objectKey string) string {
	base := strings.TrimRight(strings.TrimSpace(g.filesPublicBase), "/")
	if base == "" || g.files == nil || objectKey == "" {
		return ""
	}
	bucket := g.files.BucketName()
	return fmt.Sprintf("%s/%s/%s", base, bucket, objectKey)
}

func (g *CourseGateway) uploadSlideFile(ctx context.Context, moduleID uint64, payload *entity.PayloadJSON, fh *multipart.FileHeader) *netsp.Response[netsp.ErrorDetail] {
	if g.files == nil || fh == nil {
		return nil
	}
	if payload == nil {
		return nil
	}
	if *payload == nil {
		*payload = entity.PayloadJSON{}
	}
	ext := filepath.Ext(fh.Filename)
	if ext == "" {
		ext = ".bin"
	}
	objectKey := fmt.Sprintf("slides/%d/%s%s", moduleID, uuid.New().String(), ext)
	rc, err := fh.Open()
	if err != nil {
		return netsp.BuildError(netstatus.CodeBadRequest, netsp.ErrorDetail{
			Title:    "File Read Error",
			Message:  err.Error(),
			Solution: "Retry the upload with a valid file",
		})
	}
	defer rc.Close()

	size := fh.Size
	if size <= 0 {
		size = -1
	}

	ct := fh.Header.Get("Content-Type")
	if err := g.files.Put(ctx, objectKey, rc, size, ct); err != nil {
		slog.Warn("slide file: object storage upload failed",
			"error", err, "object_key", objectKey, "original_name", fh.Filename)
		return netsp.BuildError(http.StatusInternalServerError, netsp.ErrorDetail{
			Title:    "File Upload Failed",
			Message:  "Could not store slide file in object storage",
			Solution: "Check object storage configuration and retry upload",
		})
	}

	(*payload)["object_key"] = objectKey
	(*payload)["filename"] = filepath.Base(fh.Filename)
	if u := g.publicObjectURL(objectKey); u != "" {
		(*payload)["url"] = u
	}
	return nil
}

func (g *CourseGateway) removeSlideObject(ctx context.Context, p entity.PayloadJSON) {
	if g.files == nil {
		return
	}
	k := objectKeyFromPayload(p)
	if k == "" {
		return
	}
	_ = g.files.Remove(ctx, k)
}

// encodeObjectKeyForURLPath escapes each path segment so keys like "slides/1/name with spaces.docx"
// can be used in a single path wildcard after /api/files/.
func encodeObjectKeyForURLPath(objectKey string) string {
	objectKey = strings.Trim(objectKey, "/")
	if objectKey == "" {
		return ""
	}
	parts := strings.Split(objectKey, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}

// slidePayloadForAPIResponse returns a payload safe to expose to clients; for FILE slides with
// object_key it adds file_src (same-origin URL) so the SPA does not guess /api/files/{filename}.
func slidePayloadForAPIResponse(p entity.PayloadJSON, st entity.SlideType) entity.PayloadJSON {
	if st != entity.SlideTypeFile {
		return p
	}
	if p == nil {
		return nil
	}
	out := clonePayload(p)
	k := objectKeyFromPayload(out)
	if k == "" {
		return out
	}
	enc := encodeObjectKeyForURLPath(k)
	if enc != "" {
		out["file_src"] = "/api/files/" + enc
	}
	return out
}

func slideEntityToResponse(s entity.Slide) dto.SlideResponse {
	return dto.SlideResponse{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		SlideType:   s.SlideType,
		Payload:     slidePayloadForAPIResponse(s.Payload, s.SlideType),
		CreatedAt:   s.CreatedAt,
	}
}

// writeObjectKeyResponse streams an object from storage by key (slides/...). Returns false if response already sent with error.
func (g *CourseGateway) writeObjectKeyResponse(c *gin.Context, objectKey string) bool {
	if g.writeLocalCourseFileResponse(c, objectKey) {
		return true
	}
	if g.files == nil {
		c.Status(http.StatusServiceUnavailable)
		return false
	}
	if objectKey == "" || strings.Contains(objectKey, "..") || !strings.HasPrefix(objectKey, "slides/") {
		c.Status(http.StatusNotFound)
		return false
	}
	rc, meta, err := g.files.Open(c.Request.Context(), objectKey)
	if err != nil {
		c.Status(http.StatusNotFound)
		return false
	}
	defer rc.Close()
	if meta.ContentType != "" {
		c.Header("Content-Type", meta.ContentType)
	}
	if meta.Size > 0 {
		c.Header("Content-Length", fmt.Sprintf("%d", meta.Size))
	}
	c.Header("Content-Disposition", "inline")
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, rc)
	return true
}

func (g *CourseGateway) writeLocalCourseFileResponse(c *gin.Context, objectKey string) bool {
	localPath, ok := resolveLocalCourseFilePath(objectKey)
	if !ok {
		return false
	}

	f, err := os.Open(localPath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return true
	}
	defer f.Close()

	if info, err := f.Stat(); err == nil {
		if ct := mime.TypeByExtension(filepath.Ext(localPath)); ct != "" {
			c.Header("Content-Type", ct)
		}
		c.Header("Content-Length", fmt.Sprintf("%d", info.Size()))
		c.Header("Content-Disposition", "inline")
		http.ServeContent(c.Writer, c.Request, filepath.Base(localPath), info.ModTime(), f)
		return true
	}

	c.Status(http.StatusNotFound)
	return true
}

func resolveLocalCourseFilePath(objectKey string) (string, bool) {
	const prefix = "course_files/"
	objectKey = strings.Trim(strings.TrimSpace(objectKey), "/")
	if objectKey == "" || !strings.HasPrefix(objectKey, prefix) {
		return "", false
	}

	rel := filepath.Clean(strings.TrimPrefix(objectKey, prefix))
	if rel == "." || rel == "" || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", true
	}

	candidates := []string{
		filepath.Join("resources", "course_files"),
		filepath.Join("..", "resources", "course_files"),
	}

	for _, base := range candidates {
		absBase, err := filepath.Abs(base)
		if err != nil {
			continue
		}
		target := filepath.Join(absBase, rel)
		absTarget, err := filepath.Abs(target)
		if err != nil {
			continue
		}
		if absTarget != absBase && !strings.HasPrefix(absTarget, absBase+string(os.PathSeparator)) {
			continue
		}
		if info, err := os.Stat(absTarget); err == nil && !info.IsDir() {
			return absTarget, true
		}
	}

	return "", true
}

// ServeUploadedObject streams a file from object storage. Path must be the object_key (e.g. slides/<moduleId>/<uuid>.ext).
func (g *CourseGateway) ServeUploadedObject(c *gin.Context) {
	raw := strings.TrimPrefix(c.Param("filepath"), "/")
	key, err := url.PathUnescape(raw)
	if err != nil || key == "" {
		c.Status(http.StatusNotFound)
		return
	}
	_ = g.writeObjectKeyResponse(c, key)
}

// GetSlideFile streams the FILE slide binary for GET .../slides/:slideId/file (avoids collision with .../slides/:slideId/:optionId).
func (g *CourseGateway) GetSlideFile(c *gin.Context) {
	moduleID, errMod := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	slideID, errSlide := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if errSlide != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID урока"})
		return
	}
	if errMod != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID модуля"})
		return
	}

	mod, errResp := g.moduleUC.GetModule(c, moduleID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}
	found := false
	for i := range mod.Slides {
		if mod.Slides[i].ID == slideID {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Урок не найден в этом модуле"})
		return
	}

	slide, errResp := g.slideUC.GetSlide(c, slideID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}
	if slide.SlideType != entity.SlideTypeFile {
		c.Status(http.StatusNotFound)
		return
	}
	key := objectKeyFromPayload(slide.Payload)
	if key == "" {
		c.Status(http.StatusNotFound)
		return
	}
	_ = g.writeObjectKeyResponse(c, key)
}
