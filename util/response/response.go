package response

import "github.com/labstack/echo/v4"

// FieldError dipakai buat ngasih tahu field mana yang gagal validasi.
type FieldError struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

// Success bentuk respons sukses yang konsisten.
func Success(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// Error bentuk respons gagal yang konsisten dan fieldErrors boleh nil kalau bukan error validasi.
func Error(c echo.Context, status int, message string, fieldErrors []FieldError) error {
	body := map[string]interface{}{
		"success": false,
		"message": message,
	}
	if len(fieldErrors) > 0 {
		body["errors"] = fieldErrors
	}
	return c.JSON(status, body)
}
