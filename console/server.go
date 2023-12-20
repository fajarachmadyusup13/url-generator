package console

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/fajarachmadyusup13/url-generator/config"
	"github.com/fajarachmadyusup13/url-generator/db"
	"github.com/fajarachmadyusup13/url-generator/delivery/httpsvc"
	"github.com/fajarachmadyusup13/url-generator/repository"
	"github.com/fajarachmadyusup13/url-generator/usecase"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  "This subcommand start the server",
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {

	db.InitializeMySQLConn()

	urlRepo := repository.NewUrlRepository(db.MySQL)

	urlUsecase := usecase.NewUrlUsecase(urlRepo)

	httpService := httpsvc.NewHTTPService()
	httpService.RegisterUrlUsecase(&urlUsecase)

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		<-sigCh
		errCh <- errors.New("received an interupt")
		db.StopTickerCh <- true
	}()

	go runHTTPServer(httpService, errCh)
	log.Error(<-errCh)

}

func runHTTPServer(httpService *httpsvc.HTTPService, errCh chan<- error) {
	g := gin.Default()

	httpService.InitRoutes(g)
	errCh <- g.Run(fmt.Sprintf("0.0.0.0:%s", config.HTTPPort()))
}
