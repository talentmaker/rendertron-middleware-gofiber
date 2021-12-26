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

// A default set of user agent patterns for bots/crawlers that do not perform well with pages that require JavaScript.
var BotUserAgents = []string{
	"Baiduspider",
	"bingbot",
	"Embedly",
	"facebookexternalhit",
	"LinkedInBot",
	"outbrain",
	"pinterest",
	"quora link preview",
	"rogerbot",
	"showyoubot",
	"Slackbot",
	"TelegramBot",
	"Twitterbot",
	"vkShare",
	"W3C_Validator",
	"WhatsApp",
}

// A default set of file extensions for static assets that do not need to be proxied.
var StaticFileExtensions = []string{
	"ai",
	"avi",
	"css",
	"dat",
	"dmg",
	"doc",
	"doc",
	"exe",
	"flv",
	"gif",
	"ico",
	"iso",
	"jpeg",
	"jpg",
	"js",
	"less",
	"m4a",
	"m4v",
	"mov",
	"mp3",
	"mp4",
	"mpeg",
	"mpg",
	"pdf",
	"png",
	"ppt",
	"psd",
	"rar",
	"rss",
	"svg",
	"swf",
	"tif",
	"torrent",
	"ttf",
	"txt",
	"wav",
	"wmv",
	"woff",
	"xls",
	"xml",
	"zip",
}
