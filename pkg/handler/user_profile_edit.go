package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
	uuid "github.com/satori/go.uuid"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
)

type SubmitUserProfile struct {
	Name        string `form:"name"`
	ContactInfo string `form:"contactInfo"`
	About       string `form:"about"`
}

func (h *Handler) handleEditUserProfile(c echo.Context) error {

	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "user_profile_edit_view", map[string]interface{}{
			"Title": "Hello",
		})
	}

	authenticatedUser, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	if authenticatedUser.ID != user.ID {
		return echo.ErrForbidden
	}

	var payload SubmitUserProfile
	if err := c.Bind(&payload); err != nil {
		return err
	}

	user.Name = payload.Name
	user.ContactInfo = payload.ContactInfo
	user.About = payload.About


	file, err := c.FormFile("profilePicture")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		img, _, err := image.Decode(src)
		if err != nil {
			return err
		}

		uploadDir := os.Getenv("PUBLIC_DIR")
		if uploadDir == "" {
			uploadDir = "public"
		}

		id := uuid.NewV4().String()
		imgPath := fmt.Sprintf("%s/images/users/%s/%s", uploadDir, user.ID, id)
		if err := os.MkdirAll(imgPath, 0700); err != nil {
			return err
		}

		imgDest, err := os.Create(fmt.Sprintf("%s/full.jpg", imgPath))
		if err != nil {
			return err
		}
		defer imgDest.Close()
		resized := resize.Thumbnail(400, 400, img, resize.Lanczos3)
		if err := jpeg.Encode(imgDest, resized, &jpeg.Options{
			Quality: 60,
		}); err != nil {
			return err
		}

		bounds := img.Bounds()
		size := bounds.Size()
		var min int
		var offsetX int
		var offsetY int
		if size.X < size.Y {
			min = size.X
			offsetX = 0
			offsetY = (size.Y - size.X) / 2
		} else {
			min = size.Y
			offsetY = 0
			offsetX = (size.X - size.Y) / 2
		}
		thumbR := image.Rect(0, 0, min, min)
		thumb := image.NewRGBA(thumbR)
		draw.Draw(thumb, thumbR, img, image.Point{
			X: offsetX,
			Y: offsetY,
		}, draw.Src)
		thumbDest, err := os.Create(fmt.Sprintf("%s/thumb.jpg", imgPath))
		if err != nil {
			return err
		}
		defer thumbDest.Close()
		resized = resize.Thumbnail(60, 60, thumb, resize.Lanczos3)
		if err := jpeg.Encode(thumbDest, resized, &jpeg.Options{
			Quality: 60,
		}); err != nil {
			return err
		}

		if user.ProfilePictureID != "" {
			_ = os.Remove(fmt.Sprintf("%s", fmt.Sprintf("%s/images/users/%s/%s/full.jpg", uploadDir, user.ID, user.ProfilePictureID)))
			_ = os.Remove(fmt.Sprintf("%s", fmt.Sprintf("%s/images/users/%s/%s/thumb.jpg", uploadDir, user.ID, user.ProfilePictureID)))
			_ = os.Remove(fmt.Sprintf("%s", fmt.Sprintf("%s/images/users/%s/%s", uploadDir, user.ID, user.ProfilePictureID)))
		}

		user.ProfilePictureID = id

	}

	if err := h.userStore.Save(user); err != nil {
		return err
	}

	c.Response().Header().Set("Location", fmt.Sprintf("%s://%s/users/%s/profile", c.Scheme(), c.Request().Host, user.ID))
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil
}
