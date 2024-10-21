package ui

import "embed"

//go:embed html/* html/pages/* html/partials/*
var TemplateFiles embed.FS
