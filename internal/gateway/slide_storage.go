package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/google/uuid"

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
		// MinIO misconfiguration (e.g. Access Denied) should not block saving the slide:
		// filename is already on the payload from multipart binding; media URL may be absent until storage works.
		slog.Warn("slide file: object storage upload failed, saving slide metadata only",
			"error", err, "object_key", objectKey, "original_name", fh.Filename)
		return nil
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
