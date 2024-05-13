package mybigcamel

import (
	"strings"
)

// Copied from golint
var cil = []string{
	"ACL", "API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID",
	"HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS",
	"RPC", "SLA", "SMTP", "SQL", "SSH", "TCP", "TLS", "TTL", "UDP", "UI",
	"UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XMPP", "XSRF", "XSS",
}

var ssr *strings.Replacer
var sReplace *strings.Replacer

func init() {
	var cir []string
	var ufr []string
	for i := len(cil) - 1; i >= 0; i-- {
		initialism := cil[i]
		cir = append(cir, initialism, strings.Title(strings.ToLower(initialism)))
		ufr = append(ufr, strings.Title(strings.ToLower(initialism)), initialism)
	}

	ssr = strings.NewReplacer(cir...)
	sReplace = strings.NewReplacer(ufr...)
}
