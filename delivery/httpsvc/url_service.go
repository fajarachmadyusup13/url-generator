package httpsvc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fajarachmadyusup13/url-generator/delivery/httpsvc/middleware"
	httpsvcmodel "github.com/fajarachmadyusup13/url-generator/delivery/httpsvc/model"
	"github.com/fajarachmadyusup13/url-generator/model"
	"github.com/gin-gonic/gin"
)

func (s *HTTPService) GenShortUrl(c *gin.Context) {
	ctx := c.Request.Context()
	body := httpsvcmodel.GenerateShortUrlRequest{}

	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		err = middleware.NewHTTPError(http.StatusBadRequest, err.Error())
		c.Error(err)
	}

	res, err := s.urlUsecase.GenerateShortURL(ctx, body.Url)
	if err != nil {
		err = middleware.NewHTTPError(http.StatusInternalServerError, err.Error())
		c.Error(err)
	}

	c.JSON(http.StatusCreated, res)
}

func (s *HTTPService) UpdateShortUrl(c *gin.Context) {
	ctx := c.Request.Context()
	body := httpsvcmodel.UpdateShortUrlRequest{}

	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		err = middleware.NewHTTPError(http.StatusBadRequest, err.Error())
		c.Error(err)
	}

	url := &model.Url{
		ID:   body.ID,
		Url:  body.Url,
		Slug: body.Slug,
	}

	res, err := s.urlUsecase.UpdateURL(ctx, url)
	if err != nil {
		err = middleware.NewHTTPError(http.StatusInternalServerError, err.Error())
		c.Error(err)
	}

	c.JSON(http.StatusCreated, res)
}

func (s *HTTPService) FindByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		err = middleware.NewHTTPError(http.StatusBadRequest, err.Error())
		c.Error(err)
	}

	url, err := s.urlUsecase.FindByID(ctx, int64(intID))
	if err != nil {
		err = middleware.NewHTTPError(http.StatusInternalServerError, err.Error())
		c.Error(err)
	}

	c.JSON(http.StatusOK, &url)
}

func (s *HTTPService) FindBySlug(c *gin.Context) {
	ctx := c.Request.Context()

	slug := c.Query("slug")

	url, err := s.urlUsecase.FindBySlug(ctx, slug)
	if err != nil {
		err = middleware.NewHTTPError(http.StatusInternalServerError, err.Error())
		c.Error(err)
	}

	c.JSON(http.StatusOK, &url)
}

func (s *HTTPService) RedirectUrl(c *gin.Context) {
	ctx := c.Request.Context()

	slug := c.Param("slug")

	fmt.Println("SLUGGG: ", slug)

	url, err := s.urlUsecase.FindBySlug(ctx, slug)
	if err != nil {
		err = middleware.NewHTTPError(http.StatusInternalServerError, err.Error())
		c.Error(err)
	}

	fmt.Println("APAPAP: ", url)
	c.Redirect(302, url.Url)
	c.Abort()
}
