package httpsvc

import (
	"github.com/fajarachmadyusup13/url-generator/delivery/httpsvc/middleware"
	"github.com/fajarachmadyusup13/url-generator/model"
	"github.com/gin-gonic/gin"
)

type HTTPService struct {
	urlUsecase model.UrlUsecase
}

func NewHTTPService() *HTTPService {
	return new(HTTPService)
}

func (s *HTTPService) InitRoutes(route *gin.Engine) {
	route.ForwardedByClientIP = true
	route.SetTrustedProxies([]string{"youtube.com"})
	route.Use(middleware.CustomErrorMiddleware)

	shortener := route.Group("/shortener")
	shortener.POST("generate", s.GenShortUrl)
	shortener.POST("update", s.UpdateShortUrl)
	shortener.GET("findByID", s.FindByID)
	shortener.GET("findBySlug", s.FindBySlug)

	route.GET("/:slug", s.RedirectUrl)

}

func (s *HTTPService) RegisterUrlUsecase(urlUsecase *model.UrlUsecase) {
	s.urlUsecase = *urlUsecase
}
