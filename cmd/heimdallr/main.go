package main

import (
	"flag"
	"log"
	"os"
	"time"

	heimdallrv1 "github.com/jeromefroe/heimdallr/pkg/apis/heimdallr/v1alpha1"
	clientset "github.com/jeromefroe/heimdallr/pkg/client/clientset/versioned"
	"github.com/jeromefroe/heimdallr/pkg/controller"
	"github.com/jeromefroe/heimdallr/pkg/pingdom"

	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func main() {
	var (
		username = flag.String("username", os.Getenv("PINGDOM_USERNAME"), "Pingdom Username")
		password = flag.String("password", os.Getenv("PINGDOM_PASSWORD"), "Pingdom Password")
		appkey   = flag.String("appkey", os.Getenv("PINGDOM_APPKEY"), "Pingdom Application Key")
	)
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		logger.Fatal("unable to create in cluster config", zap.Error(err))
	}

	cli, err := clientset.NewForConfig(cfg)
	if err != nil {
		logger.Fatal("unable to create heimdallr client", zap.Error(err))
	}

	lw := cache.NewListWatchFromClient(
		cli.HeimdallrV1alpha1().RESTClient(), heimdallrv1.ResourcePlural, v1.NamespaceAll, fields.Everything(),
	)
	sw := cache.NewSharedInformer(lw, new(heimdallrv1.HTTPCheck), time.Duration(0)) // resync timer disabled

	pc, err := pingdom.New(*username, *password, *appkey, logger)
	if err != nil {
		logger.Fatal("unable to create pingdom client", zap.Error(err))
	}
	logger.Info("successfully created Pingdom client")

	ctrl := controller.New(pc, logger)
	sw.AddEventHandler(ctrl)

	logger.Info("starting controller")
	sw.Run(nil)
}
