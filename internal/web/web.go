package web

import "embed"

//go:embed dist/index.html
var IndexHtml embed.FS

//go:embed dist/assets/*
var Assets embed.FS

//go:embed dist/favicon.ico
var Favicon embed.FS
