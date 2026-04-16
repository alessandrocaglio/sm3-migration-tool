package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/checkers"
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/tools/clientcmd"
	"regexp"
	"github.com/spf13/viper"
)

type Server struct {
	router *gin.Engine
	mode   string
}

// ScanRequest represents the body of a scan request.
type ScanRequest struct {
	Namespace string `json:"namespace"`
}

// NewServer creates a new API server.
func NewServer(mode string) *Server {
	r := gin.Default()

	s := &Server{
		router: r,
		mode:   mode,
	}
	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	staticFS, err := StaticFS()
	if err == nil {
		s.router.StaticFS("/ui", staticFS)
		s.router.NoRoute(func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/ui/")
		})
	}

	api := s.router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		api.GET("/namespaces", s.handleListNamespaces)
		api.GET("/scan", s.handleScan)
		api.POST("/scan", s.handleScan)
	}
}

func (s *Server) handleScan(c *gin.Context) {
	var namespace string
	if c.Request.Method == "GET" {
		namespace = c.Query("namespace")
	} else {
		var req ScanRequest
		if err := c.ShouldBindJSON(&req); err == nil {
			namespace = req.Namespace
		}
	}

	// Validate namespace against allowed regex
	allowedRegex := viper.GetString("allowed-namespaces-regex")
	if namespace != "" && allowedRegex != "" {
		matched, _ := regexp.MatchString(allowedRegex, namespace)
		if !matched {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Namespace %s is not allowed by configuration pattern", namespace)})
			return
		}
	}

	config := discovery.DiscoveryConfig{
		ControlPlaneNamespace: namespace,
	}

	// Discovery engine
	var disc discovery.Discovery
	var err error

	if s.mode == "mock" {
		disc = discovery.NewMockDiscovery("testdata/mock-cluster.yaml", config)
	} else {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		restConfig, err := kubeConfig.ClientConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load kubeconfig: " + err.Error()})
			return
		}
		disc, err = discovery.NewLiveDiscovery(config, restConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create live discovery: " + err.Error()})
			return
		}
	}

	state, err := disc.Discover(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	allCheckers := []checkers.Checker{
		checkers.NewGatewayChecker(),
		checkers.NewServiceEntryChecker(),
		checkers.NewMtlsChecker(),
		checkers.NewAddonsChecker(),
		checkers.NewNetworkPolicyChecker(),
		checkers.NewRouteChecker(),
	}

	var allFindings []checkers.Finding
	for _, chk := range allCheckers {
		findings, err := chk.Check(context.Background(), state)
		if err != nil {
			continue
		}
		allFindings = append(allFindings, findings...)
	}

	c.JSON(http.StatusOK, gin.H{
		"findings": allFindings,
		"count":    len(allFindings),
		"scanned_namespace": namespace,
	})
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) handleListNamespaces(c *gin.Context) {
	// Discovery engine
	var disc discovery.Discovery
	var err error

	config := discovery.DiscoveryConfig{}

	if s.mode == "mock" {
		disc = discovery.NewMockDiscovery("testdata/mock-cluster.yaml", config)
	} else {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		restConfig, err := kubeConfig.ClientConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load kubeconfig: " + err.Error()})
			return
		}
		disc, err = discovery.NewLiveDiscovery(config, restConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create live discovery: " + err.Error()})
			return
		}
	}

	state, err := disc.Discover(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	allowedRegex := viper.GetString("allowed-namespaces-regex")
	var filtered []string
	seen := make(map[string]bool)

	for _, ns := range state.Namespaces {
		if allowedRegex != "" {
			matched, _ := regexp.MatchString(allowedRegex, ns.Name)
			if !matched {
				continue
			}
		}
		if !seen[ns.Name] {
			filtered = append(filtered, ns.Name)
			seen[ns.Name] = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"namespaces": filtered,
		"regex":      allowedRegex,
	})
}
