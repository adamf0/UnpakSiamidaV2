package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	"UnpakSiamida/common/infrastructure/emailhub"
	commonpresentation "UnpakSiamida/common/presentation"

	CreateUser "UnpakSiamida/modules/user/application/CreateUser"
	DeleteUser "UnpakSiamida/modules/user/application/DeleteUser"
	getEmailHistories "UnpakSiamida/modules/user/application/GetEmailHistories"
	getEmailStatus "UnpakSiamida/modules/user/application/GetEmailStatus"
	getGoogleAuthStatus "UnpakSiamida/modules/user/application/GetGoogleAuthStatus"
	GetAllUserOptions "UnpakSiamida/modules/user/application/GetAllUserOptions"
	GetAllUsers "UnpakSiamida/modules/user/application/GetAllUsers"
	GetUser "UnpakSiamida/modules/user/application/GetUser"
	MergeUser "UnpakSiamida/modules/user/application/MergeUser"
	SendAccountEmail "UnpakSiamida/modules/user/application/SendAccountEmail"
	SetupUuidUser "UnpakSiamida/modules/user/application/SetupUuidUser"
	UpdateUser "UnpakSiamida/modules/user/application/UpdateUser"
	verifyEmailHistory "UnpakSiamida/modules/user/application/VerifyEmailHistory"
	userdomain "UnpakSiamida/modules/user/domain"
	eventUser "UnpakSiamida/modules/user/event"
)

// =======================================================
// POST /user
// =======================================================

// CreateUserHandler godoc
// @Summary Create new User
// @Tags User
// @Param name formData string true "Nama"
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param email formData string true "Email"
// @Param level formData string true "Level User"
// @Param FakultasUnit formData string false "Fakultas / Unit UUID"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created user"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /user [post]
func CreateUserHandler(c *fiber.Ctx) error {
	name := strings.TrimSpace(c.FormValue("name"))
	username := strings.TrimSpace(c.FormValue("username"))
	password := strings.TrimSpace(c.FormValue("password"))
	email := strings.TrimSpace(c.FormValue("email"))
	level := strings.ToLower(strings.TrimSpace(c.FormValue("level")))
	if level != "admin" {
		level = "user"
	}

	fu := strings.TrimSpace(c.FormValue("FakultasUnit"))
	if fu == "" {
		fu = strings.TrimSpace(c.FormValue("uuid_fakultas_unit"))
	}

	cmd := CreateUser.CreateUserCommand{
		Name:             name,
		Username:         username,
		Password:         password,
		Email:            email,
		Level:            level,
		UuidFakultasUnit: helper.StrPtr(fu),
	}

	uuid, err := mediatr.Send[
		CreateUser.CreateUserCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /user/{uuid}
// =======================================================

// UpdateUserHandler godoc
// @Summary Update existing User
// @Tags User
// @Param uuid path string true "User UUID" format(uuid)
// @Param name formData string true "Nama"
// @Param username formData string true "Username"
// @Param password formData string false "Password (optional)"
// @Param email formData string true "Email"
// @Param level formData string true "Level User"
// @Param FakultasUnit formData string false "Fakultas / Unit UUID"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated user"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /user/{uuid} [put]
func UpdateUserHandler(c *fiber.Ctx) error {
	name := strings.TrimSpace(c.FormValue("name"))
	username := strings.TrimSpace(c.FormValue("username"))
	password := strings.TrimSpace(c.FormValue("password"))
	email := strings.TrimSpace(c.FormValue("email"))
	level := strings.ToLower(strings.TrimSpace(c.FormValue("level")))
	if level != "admin" {
		level = "user"
	}

	fu := strings.TrimSpace(c.FormValue("FakultasUnit"))
	if fu == "" {
		fu = strings.TrimSpace(c.FormValue("uuid_fakultas_unit"))
	}

	var passwordPtr *string
	if password != "" {
		passwordPtr = &password
	}

	cmd := UpdateUser.UpdateUserCommand{
		Uuid:             c.Params("uuid"),
		Name:             name,
		Username:         username,
		Password:         passwordPtr,
		Email:            email,
		Level:            level,
		UuidFakultasUnit: helper.StrPtr(fu),
	}

	updatedID, err := mediatr.Send[
		UpdateUser.UpdateUserCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /user/{uuid}
// =======================================================

// DeleteUserHandler godoc
// @Summary Delete User
// @Tags User
// @Param uuid path string true "User UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted user"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /user/{uuid} [delete]
func DeleteUserHandler(c *fiber.Ctx) error {

	cmd := DeleteUser.DeleteUserCommand{
		Uuid: c.Params("uuid"),
	}

	deletedID, err := mediatr.Send[
		DeleteUser.DeleteUserCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /user/{uuid}
// =======================================================

// GetUserHandler godoc
// @Summary Get User by UUID
// @Tags User
// @Param uuid path string true "User UUID" format(uuid)
// @Produce json
// @Success 200 {object} userdomain.User
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
//
//	@Router /user/{uuid} [get]
func GetUserHandler(c *fiber.Ctx) error {

	query := GetUser.GetUserByUuidQuery{
		Uuid: c.Params("uuid"),
	}

	user, err := mediatr.Send[
		GetUser.GetUserByUuidQuery,
		*userdomain.User,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// =======================================================
// GET /users
// =======================================================

// GetAllUsersHandler godoc
// @Summary Get All Users
// @Tags User
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Param filters query string false "Search filters (field:op:value;...)"
// @Produce json
// @Success 200 {object} commondomain.Paged[userdomain.User]
// @Router /users [get]
func GetAllUsersHandler(c *fiber.Ctx) error {

	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	filtersRaw := c.Query("filters", "")
	var filters []commondomain.SearchFilter

	if filtersRaw != "" {
		parts := strings.Split(filtersRaw, ";")
		for _, p := range parts {
			tokens := strings.SplitN(p, ":", 3)
			if len(tokens) != 3 {
				continue
			}

			field := strings.TrimSpace(tokens[0])
			op := strings.TrimSpace(tokens[1])
			rawValue := strings.TrimSpace(tokens[2])

			var ptr *string
			if rawValue != "" && rawValue != "null" {
				ptr = &rawValue
			}

			filters = append(filters, commondomain.SearchFilter{
				Field:    field,
				Operator: op,
				Value:    ptr,
			})
		}
	}

	query := GetAllUsers.GetAllUsersQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[userdomain.User]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[userdomain.User]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[userdomain.User]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[userdomain.User]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[userdomain.User]{}
	}

	result, err := mediatr.Send[
		GetAllUsers.GetAllUsersQuery,
		commondomain.Paged[userdomain.User],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

// =======================================================
// GET /user-options
// =======================================================

// GetAllUserOptionsHandler godoc
// @Summary Get All User Options
// @Tags User
// @Produce json
// @Success 200 {object} userdomain.UserOptions
// @Router /user-options [get]
func GetAllUserOptionsHandler(c *fiber.Ctx) error {
	query := GetAllUserOptions.GetAllUserOptionsQuery{}
	result, err := mediatr.Send[
		GetAllUserOptions.GetAllUserOptionsQuery,
		[]userdomain.UserOptions,
	](c.Context(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func SetupUuidUsersHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidUser.SetupUuidUserCommand{}

	message, err := mediatr.Send[SetupUuidUser.SetupUuidUserCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

// =======================================================
// POST /users/merge
// =======================================================

// MergeUserHandler godoc
// @Summary Merge multiple secondary users into a primary user
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} map[string]string "message"
// @Router /users/merge [post]
func MergeUserHandler(c *fiber.Ctx) error {
	primaryUUID := strings.TrimSpace(c.FormValue("primary_uuid"))
	var primaryID *uint
	if val := c.FormValue("primary_id"); val != "" {
		if id, err := strconv.ParseUint(val, 10, 32); err == nil {
			uid := uint(id)
			primaryID = &uid
		}
	}

	var secondaryUUIDs []string
	if val := c.FormValue("secondary_uuids"); val != "" {
		if err := json.Unmarshal([]byte(val), &secondaryUUIDs); err != nil {
			for _, s := range strings.Split(val, ",") {
				if trimmed := strings.TrimSpace(s); trimmed != "" {
					secondaryUUIDs = append(secondaryUUIDs, trimmed)
				}
			}
		}
	}

	var secondaryIDs []uint
	if val := c.FormValue("secondary_ids"); val != "" {
		_ = json.Unmarshal([]byte(val), &secondaryIDs)
	}

	cmd := MergeUser.MergeUserCommand{
		PrimaryUUID:    primaryUUID,
		SecondaryUUIDs: secondaryUUIDs,
		PrimaryID:      primaryID,
		SecondaryIDs:   secondaryIDs,
	}

	message, err := mediatr.Send[MergeUser.MergeUserCommand, string](c.Context(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": message,
	})
}

const (
	maxEmailAttachments     = 5
	maxEmailAttachmentSize  = 10 << 20 // 10 MB per file
	maxEmailAttachmentTotal = 15 << 20 // 15 MB total (safe under Gmail 25 MB MIME limit)
)

var allowedEmailAttachmentExts = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
	".txt":  true,
	".csv":  true,
	".zip":  true,
}

func sanitizeUploadFilename(name string) string {
	base := filepath.Base(strings.ReplaceAll(strings.TrimSpace(name), "\\", "/"))
	var b strings.Builder
	for _, r := range base {
		switch {
		case r < 32 || r == 127 || r == '"' || r == '/' || r == '\\':
			continue
		default:
			b.WriteRune(r)
		}
	}
	cleaned := strings.TrimSpace(b.String())
	if cleaned == "" || cleaned == "." || cleaned == ".." {
		return "lampiran"
	}
	return cleaned
}

func stageEmailAttachments(c *fiber.Ctx) ([]eventUser.UserAccountEmailAttachment, string, error) {
	form, err := c.MultipartForm()
	if err != nil || form == nil || len(form.File) == 0 {
		return nil, "", nil
	}

	var fileHeaders []*multipart.FileHeader
	for _, fieldName := range []string{"attachments", "attachments[]", "files", "files[]", "file"} {
		if list, ok := form.File[fieldName]; ok && len(list) > 0 {
			fileHeaders = append(fileHeaders, list...)
		}
	}
	if len(fileHeaders) == 0 {
		return nil, "", nil
	}
	if len(fileHeaders) > maxEmailAttachments {
		return nil, "", fmt.Errorf("maksimal %d file lampiran per pengiriman email", maxEmailAttachments)
	}

	var totalSize int64
	for _, fh := range fileHeaders {
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		if !allowedEmailAttachmentExts[ext] {
			return nil, "", fmt.Errorf("tipe file %q tidak diizinkan untuk lampiran email", fh.Filename)
		}
		if fh.Size > maxEmailAttachmentSize {
			return nil, "", fmt.Errorf("ukuran file %q melebihi batas maksimal 10 MB", fh.Filename)
		}
		totalSize += fh.Size
	}
	if totalSize > maxEmailAttachmentTotal {
		return nil, "", fmt.Errorf("total ukuran seluruh lampiran melebihi batas maksimal 15 MB")
	}

	batchDir := filepath.Join(eventUser.EmailAttachmentBaseDir, uuid.NewString())
	if err := os.MkdirAll(batchDir, 0o755); err != nil {
		return nil, "", fmt.Errorf("gagal membuat direktori lampiran: %w", err)
	}

	attachments := make([]eventUser.UserAccountEmailAttachment, 0, len(fileHeaders))
	for i, fh := range fileHeaders {
		safeName := sanitizeUploadFilename(fh.Filename)
		ext := strings.ToLower(filepath.Ext(safeName))
		diskName := fmt.Sprintf("%d_%s", i+1, safeName)
		destPath := filepath.Join(batchDir, diskName)

		src, err := fh.Open()
		if err != nil {
			_ = os.RemoveAll(batchDir)
			return nil, "", fmt.Errorf("gagal membaca file %q: %w", safeName, err)
		}

		sniffBuf := make([]byte, 512)
		n, readErr := src.Read(sniffBuf)
		if readErr != nil && readErr != io.EOF {
			_ = src.Close()
			_ = os.RemoveAll(batchDir)
			return nil, "", fmt.Errorf("gagal membaca isi file %q: %w", safeName, readErr)
		}

		dst, err := os.Create(destPath)
		if err != nil {
			_ = src.Close()
			_ = os.RemoveAll(batchDir)
			return nil, "", fmt.Errorf("gagal menyimpan file %q: %w", safeName, err)
		}

		written := int64(0)
		if n > 0 {
			w1, err := dst.Write(sniffBuf[:n])
			if err != nil {
				_ = dst.Close()
				_ = src.Close()
				_ = os.RemoveAll(batchDir)
				return nil, "", fmt.Errorf("gagal menulis file %q: %w", safeName, err)
			}
			written += int64(w1)
		}
		w2, err := io.Copy(dst, src)
		_ = dst.Close()
		_ = src.Close()
		if err != nil {
			_ = os.RemoveAll(batchDir)
			return nil, "", fmt.Errorf("gagal menyalin file %q: %w", safeName, err)
		}
		written += w2

		contentType := mime.TypeByExtension(ext)
		if contentType == "" && n > 0 {
			contentType = http.DetectContentType(sniffBuf[:n])
		}
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		// Ambil hanya MIME type utama tanpa parameter charset tambahan
		if semi := strings.Index(contentType, ";"); semi != -1 {
			contentType = strings.TrimSpace(contentType[:semi])
		}

		attachments = append(attachments, eventUser.UserAccountEmailAttachment{
			FileName:    safeName,
			ContentType: contentType,
			FilePath:    destPath,
			Size:        written,
		})
	}

	return attachments, batchDir, nil
}

// =======================================================
// POST /user/:uuid/send-email & POST /users/send-email
// =======================================================

// SendAccountEmailHandler godoc
// @Summary Send account credentials email via Outbox Pattern
// @Tags User
// @Param uuid path string false "User UUID"
// @Param user_uuids formData string false "JSON string or comma-separated list of User UUIDs for batch"
// @Param attachments formData file false "Optional attachment files (max 5 files, 10MB/file, 15MB total)"
// @Produce json
// @Success 200 {object} SendAccountEmail.SendAccountEmailResult
// @Router /user/{uuid}/send-email [post]
// @Router /users/send-email [post]
func SendAccountEmailHandler(c *fiber.Ctx) error {
	var userUUIDs []string

	// Ambil dari path param jika ada (/user/:uuid/send-email)
	if paramUuid := strings.TrimSpace(c.Params("uuid")); paramUuid != "" {
		userUUIDs = append(userUUIDs, paramUuid)
	}

	// Ambil dari FormValue ("user_uuids" atau "uuids" atau "uuid")
	if val := c.FormValue("user_uuids"); val != "" {
		var formArr []string
		if err := json.Unmarshal([]byte(val), &formArr); err == nil {
			userUUIDs = append(userUUIDs, formArr...)
		} else {
			for _, s := range strings.Split(val, ",") {
				if trimmed := strings.TrimSpace(s); trimmed != "" {
					userUUIDs = append(userUUIDs, trimmed)
				}
			}
		}
	}

	if val := c.FormValue("uuids"); val != "" {
		var extraUUIDs []string
		if err := json.Unmarshal([]byte(val), &extraUUIDs); err == nil {
			userUUIDs = append(userUUIDs, extraUUIDs...)
		} else {
			for _, s := range strings.Split(val, ",") {
				if trimmed := strings.TrimSpace(s); trimmed != "" {
					extraUUIDs = append(extraUUIDs, trimmed)
				}
			}
		}
		userUUIDs = append(userUUIDs, extraUUIDs...)
	}

	if val := strings.TrimSpace(c.FormValue("uuid")); val != "" {
		userUUIDs = append(userUUIDs, val)
	}

	// Deduplikasi UUID
	seen := make(map[string]bool)
	var uniqueUUIDs []string
	for _, u := range userUUIDs {
		uClean := strings.TrimSpace(u)
		if uClean != "" && !seen[uClean] {
			seen[uClean] = true
			uniqueUUIDs = append(uniqueUUIDs, uClean)
		}
	}

	if len(uniqueUUIDs) == 0 {
		return commoninfra.HandleError(c, userdomain.EmptyData())
	}

	attachments, batchDir, attErr := stageEmailAttachments(c)
	if attErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Code":    "UserSendEmail.InvalidAttachment",
			"Message": attErr.Error(),
			"message": attErr.Error(),
		})
	}

	cmd := SendAccountEmail.SendAccountEmailCommand{
		UserUUIDs:   uniqueUUIDs,
		Attachments: attachments,
	}

	res, err := mediatr.Send[SendAccountEmail.SendAccountEmailCommand, *SendAccountEmail.SendAccountEmailResult](c.Context(), cmd)
	if err != nil {
		if batchDir != "" {
			_ = os.RemoveAll(batchDir)
		}
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

// =======================================================
// EMAIL HISTORY & REALTIME CHECK HANDLERS
// =======================================================

func GetEmailStatusHandler(c *fiber.Ctx) error {
	outboxID := strings.TrimSpace(c.Query("outbox_id"))
	if outboxID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "InvalidParam",
			"message": "Parameter outbox_id wajib diisi",
		})
	}
	verifyLive := c.Query("verify") == "true"

	res, err := mediatr.Send[getEmailStatus.GetEmailStatusQuery, *getEmailStatus.GetEmailStatusResult](
		c.Context(),
		getEmailStatus.GetEmailStatusQuery{
			OutboxID:   outboxID,
			VerifyLive: verifyLive,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

func GetEmailHistoriesHandler(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	res, err := mediatr.Send[getEmailHistories.GetEmailHistoriesQuery, *getEmailHistories.GetEmailHistoriesResult](
		c.Context(),
		getEmailHistories.GetEmailHistoriesQuery{
			Search: search,
			Status: status,
			Page:   page,
			Limit:  limit,
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

func VerifyEmailHistoryHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "InvalidID",
			"message": "ID riwayat email tidak valid",
		})
	}

	res, err := mediatr.Send[verifyEmailHistory.VerifyEmailHistoryCommand, *verifyEmailHistory.VerifyEmailHistoryResult](
		c.Context(),
		verifyEmailHistory.VerifyEmailHistoryCommand{
			ID: uint(id),
		},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

func GetGoogleAuthStatusHandler(c *fiber.Ctx) error {
	res, err := mediatr.Send[getGoogleAuthStatus.GetGoogleAuthStatusQuery, *getGoogleAuthStatus.GetGoogleAuthStatusResult](
		c.Context(),
		getGoogleAuthStatus.GetGoogleAuthStatusQuery{},
	)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(res)
}

func ModuleUser(app *fiber.App) {
	admin := []string{"admin"}
	whoamiURL := "http://localhost:3000/whoami"

	app.Get("/user/setupuuid", SetupUuidUsersHandlerfunc)

	app.Post("/user", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), CreateUserHandler)
	app.Put("/user/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), UpdateUserHandler)
	app.Delete("/user/:uuid", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), DeleteUserHandler)
	app.Get("/user/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetUserHandler)
	app.Get("/users", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllUsersHandler)

	app.Get("/user-options", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllUserOptionsHandler)

	app.Post("/users/merge", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), MergeUserHandler)
	app.Post("/user/merge", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), MergeUserHandler)

	app.Post("/user/:uuid/send-email", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SendAccountEmailHandler)
	app.Post("/users/send-email", commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), SendAccountEmailHandler)

	// Status pengiriman email & Realtime check
	app.Get("/users/email-status", commonpresentation.JWTMiddleware(), GetEmailStatusHandler)
	app.Get("/users/email-history", commonpresentation.JWTMiddleware(), GetEmailHistoriesHandler)
	app.Post("/users/email-history/:id/verify", commonpresentation.JWTMiddleware(), VerifyEmailHistoryHandler)

	// Live Streaming Realtime Email Status (WebSocket & SSE)
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/email-status", websocket.New(func(conn *websocket.Conn) {
		hub := emailhub.GetHub()
		client := hub.RegisterWS(conn)
		defer hub.UnregisterWS(client)

		_ = conn.WriteJSON(fiber.Map{
			"type":    "connected",
			"message": "Connected to Email Status Live Stream",
		})

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}))
	app.Get("/users/email-history/live", emailhub.HandleSSE)

	app.Get("/auth/google/status", commonpresentation.JWTMiddleware(), GetGoogleAuthStatusHandler)

	ModuleUserAbac(app)
}
