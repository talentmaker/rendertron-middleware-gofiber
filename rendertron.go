// Package rendertron provides a re-write of the Rendertron middleware
// https://github.com/GoogleChrome/rendertron/blob/main/middleware/src/middleware.ts
//
// Copyright 2017 Google Inc. All rights reserved.
// Copyright 2021 Luke Zhang
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.
package rendertronmiddleware

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func contains(arr []string, target string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}

	return false
}

type Options struct {
	// Base URL of the Rendertron proxy service. Required.
	ProxyUrl string `json:"proxyUrl" binding:"required"`

	// Regular expression to match user agent to proxy.
	//
	// Default: set of bots that do not perform well with pages that require
	// JavaScript + configured extra bot user agents.
	UserAgentPattern string `json:"userAgentPattern"`

	// Extra bot user agents to include on top of existing ones, such as
	//
	//	[]string{
	//		"curl",
	//		"googlebot",
	//		"bingbot",
	//		"linkedinbot",
	//		"mediapartners-google",
	// 	}
	ExtraBotUserAgents []string `json:"extraBotUserAgents"`

	// Regular expression used to exclude request URL paths.
	//
	// Default: set of typical static asset file extensions + configured extra
	// exclude urls.
	ExcludeUrlPattern string `json:"excludeUrlPattern"`

	// Extra file extensions to ignore on top of existing ones
	ExtraExcludeUrls []string `json:"extraExcludeUrls"`

	// Force web components polyfills to be loaded and enabled
	//
	// Default: false.
	InjectShadyDom bool `json:"injectShadyDom"`

	// Millisecond timeout for proxy requests.
	//
	// Defaults: 11000 milliseconds.
	Timeout int `json:"timeout"`

	// If a forwarded host header is found and matches one of the hosts in this
	// array, then that host will be used for the request to the rendertron
	// server instead of the actual host of the request.
	//
	// This is usedful if this middleware is running on a different host which
	// is proxied behind the actual site, and the rendertron server should
	// request the main site.
	AllowedForwardedHosts []string `json:"allowedForwardedHosts"`

	// Header used to determine the forwarded host that should be used when
	// building the URL to be rendered.
	// Only applicable if allowedForwardedHosts is not empty.
	//
	// Defaults to "X-Forwarded-Host".
	ForwardedHostHeader string `json:"forwardedHostHeader"`
}

const (
	defaultForwardedHostHeader = "X-Forwarded-Host"
	defaultTimeout             = 11000
)

// New rendertron middleware proxies requests from the server to a Rendertron bot
// rendering service.
func New(options ...Options) fiber.Handler {
	config := options[0]

	if config.ProxyUrl == "" {
		log.Fatal("Must set options.proxyUrl.")
	}

	proxyUrl := config.ProxyUrl

	if !strings.HasSuffix(proxyUrl, "/") {
		proxyUrl += "/"
	}

	userAgentPattern := config.UserAgentPattern

	if userAgentPattern == "" {
		userAgentPattern = strings.Join(
			append(BotUserAgents, config.ExtraBotUserAgents...),
			"|",
		)
	}

	excludeUrlPattern := config.ExcludeUrlPattern

	if excludeUrlPattern == "" {
		excludeUrlPattern = fmt.Sprintf(
			"\\.(%s)$",
			strings.Join(
				append(StaticFileExtensions, config.ExtraExcludeUrls...),
				"|",
			),
		)
	}

	injectShadyDom := config.InjectShadyDom
	timeout := config.Timeout

	if timeout == 0 {
		timeout = defaultTimeout
	}

	allowedForwardedHosts := config.AllowedForwardedHosts
	var forwardedHostHeader string

	if len(config.AllowedForwardedHosts) == 0 {
		forwardedHostHeader = ""
	} else if config.ForwardedHostHeader == "" {
		forwardedHostHeader = defaultForwardedHostHeader
	} else {
		forwardedHostHeader = config.ForwardedHostHeader
	}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	return func(ctx *fiber.Ctx) error {
		userAgent := string(ctx.Request().Header.UserAgent())

		if userAgent == "" ||
			!regexp.MustCompile("(?i)"+userAgentPattern).MatchString(userAgent) ||
			regexp.MustCompile("(?i)"+excludeUrlPattern).MatchString(ctx.Path()) {

			return ctx.Next()
		}

		var forwardedHost string

		if forwardedHostHeader != "" {
			forwardedHost = ctx.Get(forwardedHostHeader)
		}

		var host string

		if forwardedHost != "" && contains(allowedForwardedHosts, forwardedHost) {
			host = forwardedHost
		} else {
			host = string(ctx.Request().Header.Host())
		}

		incomingUrl := ctx.Protocol() + "://" + host + ctx.OriginalURL()
		renderUrl := proxyUrl + url.QueryEscape(incomingUrl)

		if injectShadyDom {
			renderUrl += "?wc-inject-shadydom=true"
		}

		response, err := client.Get(renderUrl)

		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return err
		}

		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

		return ctx.Status(200).SendString(string(body))
	}
}
