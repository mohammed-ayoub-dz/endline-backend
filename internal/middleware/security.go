package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func SecurityHeaders() fiber.Handler {
	return helmet.New(helmet.Config{
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
		OriginAgentCluster:        "?1",
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		XDNSPrefetchControl:       "off",
		XFrameOptions:             "DENY",
	})
}