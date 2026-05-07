package file_app

import (
	"bytes"
	file_domain "cloud/internal/features/file/domain"
	"context"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type FileUploadedEventHandler interface {
	HandleUploaded(
		ctx context.Context,
		event FileUploadedEvent,
	) error
}

type fileUploadedEventHandler struct {
	fileRepo FileRepository
	storage  FileStorage
	metrics  FileMetrics
}

func NewFileUploadedEventHandler(
	fileRepo FileRepository,
	storage FileStorage,
	metrics FileMetrics,
) *fileUploadedEventHandler {
	return &fileUploadedEventHandler{
		fileRepo: fileRepo,
		storage:  storage,
		metrics:  metrics,
	}
}

func (h *fileUploadedEventHandler) HandleUploaded(
	ctx context.Context,
	event FileUploadedEvent,
) error {
	file, err := h.fileRepo.FindByIDAndUserID(ctx, event.FileID, event.UserID)
	if err != nil {
		return err
	}

	if file.Status.IsTerminal() {
		return nil
	}

	startedAt := time.Now()
	h.metrics.ProcessingStarted()

	if !isPreviewSupported(file) {
		_, err := h.fileRepo.UpdateStatus(ctx, file.ID, file.UserID, file_domain.FileStatusProcessed)
		if err != nil {
			h.metrics.ObserveProcessingFinished(FileProcessingResultFailed, time.Since(startedAt))
			return err
		}

		h.metrics.ObserveProcessingFinished(FileProcessingResultSuccess, time.Since(startedAt))
		return err
	}

	if err := h.createPreview(ctx, file); err != nil {
		_, statusErr := h.fileRepo.UpdateStatus(ctx, file.ID, file.UserID, file_domain.FileStatusFailed)
		if statusErr != nil {
			h.metrics.ObserveProcessingFinished(FileProcessingResultFailed, time.Since(startedAt))
			return statusErr
		}

		h.metrics.ObserveProcessingFinished(FileProcessingResultFailed, time.Since(startedAt))
		return err
	}

	previewPath := previewPath(file.StoragePath)
	previewMimetype := "image/jpeg"
	_, err = h.fileRepo.UpdatePreview(
		ctx,
		file.ID,
		file.UserID,
		file_domain.FileStatusProcessed,
		&previewPath,
		&previewMimetype,
	)
	if err != nil {
		h.metrics.ObserveProcessingFinished(FileProcessingResultFailed, time.Since(startedAt))
		return err
	}

	h.metrics.ObserveProcessingFinished(FileProcessingResultSuccess, time.Since(startedAt))
	return err
}

func isPreviewSupported(
	file file_domain.File,
) bool {
	mimetype := ""
	if file.Mimetype != nil {
		mimetype = strings.ToLower(*file.Mimetype)
	}

	if strings.HasPrefix(mimetype, "image/") || mimetype == "application/pdf" {
		return true
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".pdf"
}

func (h *fileUploadedEventHandler) createPreview(
	ctx context.Context,
	file file_domain.File,
) error {
	if isPDF(file) {
		return h.createPDFPreview(ctx, file)
	}

	return h.createImagePreview(ctx, file)
}

func isPDF(
	file file_domain.File,
) bool {
	if file.Mimetype != nil && strings.EqualFold(*file.Mimetype, "application/pdf") {
		return true
	}

	return strings.EqualFold(filepath.Ext(file.Filename), ".pdf")
}

func (h *fileUploadedEventHandler) createImagePreview(
	ctx context.Context,
	file file_domain.File,
) error {
	content, err := h.storage.Open(ctx, file.StoragePath)
	if err != nil {
		return err
	}
	defer content.Close()

	srcImage, _, err := image.Decode(content)
	if err != nil {
		return err
	}

	preview := resizeImage(srcImage, 320)

	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, preview, &jpeg.Options{Quality: 85}); err != nil {
		return err
	}

	return h.storage.SaveAtPath(ctx, previewPath(file.StoragePath), bytes.NewReader(buffer.Bytes()))
}

func (h *fileUploadedEventHandler) createPDFPreview(
	ctx context.Context,
	file file_domain.File,
) error {
	content, err := h.storage.Open(ctx, file.StoragePath)
	if err != nil {
		return err
	}
	defer content.Close()

	tempDir, err := os.MkdirTemp("", "cloud-file-preview-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	sourcePath := filepath.Join(tempDir, "source.pdf")
	sourceFile, err := os.Create(sourcePath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(sourceFile, content); err != nil {
		sourceFile.Close()
		return err
	}

	if err := sourceFile.Close(); err != nil {
		return err
	}

	outputPrefix := filepath.Join(tempDir, "preview")
	command := exec.CommandContext(
		ctx,
		"pdftoppm",
		"-f", "1",
		"-l", "1",
		"-singlefile",
		"-jpeg",
		"-scale-to", "320",
		sourcePath,
		outputPrefix,
	)

	if err := command.Run(); err != nil {
		return err
	}

	previewFile, err := os.Open(outputPrefix + ".jpg")
	if err != nil {
		return err
	}
	defer previewFile.Close()

	return h.storage.SaveAtPath(ctx, previewPath(file.StoragePath), previewFile)
}

func previewPath(
	storagePath string,
) string {
	return filepath.Join("previews", storagePath+".jpg")
}

func resizeImage(
	src image.Image,
	maxDimension int,
) image.Image {
	bounds := src.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	if srcWidth <= maxDimension && srcHeight <= maxDimension {
		return src
	}

	dstWidth := srcWidth
	dstHeight := srcHeight
	if srcWidth >= srcHeight {
		dstWidth = maxDimension
		dstHeight = srcHeight * maxDimension / srcWidth
	} else {
		dstHeight = maxDimension
		dstWidth = srcWidth * maxDimension / srcHeight
	}

	if dstWidth < 1 {
		dstWidth = 1
	}
	if dstHeight < 1 {
		dstHeight = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))
	for y := 0; y < dstHeight; y++ {
		srcY := bounds.Min.Y + y*srcHeight/dstHeight
		for x := 0; x < dstWidth; x++ {
			srcX := bounds.Min.X + x*srcWidth/dstWidth
			dst.Set(x, y, color.RGBAModel.Convert(src.At(srcX, srcY)))
		}
	}

	return dst
}
