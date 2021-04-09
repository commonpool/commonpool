package handler

import (
	"cp/pkg/api"
	"fmt"
	form "github.com/go-playground/form/v4"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
	uuid "github.com/satori/go.uuid"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"time"
)

type ExistingImages struct {
	ID      string `form:"id"`
	GroupID string `form:"groupId"`
	PostID  string `form:"postId"`
	Delete  bool   `form:"delete"`
}

type SubmitPost struct {
	GroupID        string           `param:"groupId"`
	Title          string           `form:"title"`
	Description    string           `form:"description"`
	Type           api.PostType     `form:"type"`
	ValueFrom      string           `form:"valueFrom"`
	ValueTo        string           `form:"valueTo"`
	ExistingImages []ExistingImages `form:"existingImages"`
}

func (h *Handler) handlePostEdit(c echo.Context) error {

	authenticatedUser, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	membership, err := h.getAuthenticatedUserMembership(c)
	if err != nil {
		return err
	}

	group, err := h.getGroup(c)
	if err != nil {
		return err
	}

	post, err := h.getPost(c)
	if err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "post_form", map[string]interface{}{
			"Title":      "New Post",
			"Group":      group,
			"Post":       post,
			"Membership": membership,
		})
	}

	var payload SubmitPost
	if err := c.Bind(&payload); err != nil {
		return err
	}

	req := c.Request()
	if err := req.ParseForm(); err != nil {
		return err
	}
	decoder := form.NewDecoder()
	if err := decoder.Decode(&payload, req.Form); err != nil {
		return err
	}

	var valueFromPtr *time.Duration
	var valueToPtr *time.Duration

	if payload.Type != api.CommentPost {
		valueFrom, err := time.ParseDuration(payload.ValueFrom)
		if err != nil {
			return err
		}
		valueTo, err := time.ParseDuration(payload.ValueTo)
		if err != nil {
			return err
		}
		if valueFrom < 0 {
			return fmt.Errorf("invalid value from")
		}
		if valueTo < valueFrom {
			return fmt.Errorf("value to cannot be smaller than value from")
		}
		valueFromPtr = &valueFrom
		valueToPtr = &valueTo
	}

	if payload.Type == api.CommentPost {
		if payload.Description == "" {
			return echo.ErrBadRequest
		}
	}

	if payload.Type != api.OfferPost && payload.Type != api.RequestPost && payload.Type != api.CommentPost {
		return fmt.Errorf("invalid post type")
	}

	id := uuid.NewV4().String()
	var isNewPost = true
	if post != nil {
		isNewPost = false
		id = post.ID
		if post.GroupID != payload.GroupID {
			return echo.ErrBadRequest
		}
		if post.AuthorID != authenticatedUser.ID {
			return echo.ErrBadRequest
		}
	}

	post = &api.Post{
		ID:          id,
		GroupID:     payload.GroupID,
		AuthorID:    authenticatedUser.ID,
		Title:       payload.Title,
		Description: payload.Description,
		ValueFrom:   valueFromPtr,
		ValueTo:     valueToPtr,
		Type:        payload.Type,
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	uploadDir := os.Getenv("PUBLIC_DIR")
	if uploadDir == "" {
		uploadDir = "public"
	}

	imgsPath := fmt.Sprintf("%s/images/full/groups/%s/posts/%s", uploadDir, group.ID, post.ID)
	if err := os.MkdirAll(imgsPath, 0700); err != nil {
		return err
	}
	mediumPath := fmt.Sprintf("%s/images/medium/groups/%s/posts/%s", uploadDir, group.ID, post.ID)
	if err := os.MkdirAll(mediumPath, 0700); err != nil {
		return err
	}
	thumbsPath := fmt.Sprintf("%s/images/thumb/groups/%s/posts/%s", uploadDir, group.ID, post.ID)
	if err := os.MkdirAll(thumbsPath, 0700); err != nil {
		return err
	}

	for _, existingImage := range payload.ExistingImages {
		if existingImage.Delete {
			if existingImage.GroupID != group.ID || existingImage.PostID != post.ID {
				return echo.ErrBadRequest
			}
			if err := h.imageStore.Delete(existingImage.ID); err != nil {
				return err
			}
			os.Remove(fmt.Sprintf("%s/images/full/groups/%s/posts/%s/%s.jpg", uploadDir, existingImage.GroupID, existingImage.PostID, existingImage.ID))
			os.Remove(fmt.Sprintf("%s/public/images/thumb/groups/%s/posts/%s/%s.jpg", uploadDir, existingImage.GroupID, existingImage.PostID, existingImage.ID))
			os.Remove(fmt.Sprintf("%s/images/medium/groups/%s/posts/%s/%s.jpg", uploadDir, existingImage.GroupID, existingImage.PostID, existingImage.ID))
		}
	}

	var images []*api.Image

	var i = 0
	for {
		files := form.File[fmt.Sprintf("image-%d", i)]
		if len(files) == 0 {
			break
		}
		for _, file := range files {

			id := uuid.NewV4().String()

			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			img, _, err := image.Decode(src)
			if err != nil {
				return err
			}

			thumb := resize.Resize(160, 0, img, resize.Lanczos3)
			thumbDest, err := os.Create(fmt.Sprintf("%s/%s.jpg", thumbsPath, id))
			if err != nil {
				return err
			}
			defer thumbDest.Close()
			if err := jpeg.Encode(thumbDest, thumb, &jpeg.Options{
				Quality: 60,
			}); err != nil {
				return err
			}

			medium := resize.Resize(400, 0, img, resize.Lanczos3)
			mediumDest, err := os.Create(fmt.Sprintf("%s/%s.jpg", mediumPath, id))
			if err != nil {
				return err
			}
			defer mediumDest.Close()
			if err := jpeg.Encode(mediumDest, medium, &jpeg.Options{
				Quality: 60,
			}); err != nil {
				return err
			}

			imgDest, err := os.Create(fmt.Sprintf("%s/%s.jpg", imgsPath, id))
			if err != nil {
				return err
			}
			defer imgDest.Close()

			if err := jpeg.Encode(imgDest, img, &jpeg.Options{
				Quality: 60,
			}); err != nil {
				return err
			}

			images = append(images, &api.Image{
				ID:      id,
				PostID:  post.ID,
				GroupID: group.ID,
			})

		}
		i++
	}

	if !isNewPost {
		if err := h.postStore.Update(post); err != nil {
			return err
		}
	} else {
		if err := h.postStore.Create(post); err != nil {
			return err
		}
	}

	if err := h.imageStore.Add(images); err != nil {
		return err
	}

	c.Response().Header().Set("Location", fmt.Sprintf("%s://%s/groups/%s/posts/%s", c.Scheme(), c.Request().Host, group.ID, post.ID))
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil

}
