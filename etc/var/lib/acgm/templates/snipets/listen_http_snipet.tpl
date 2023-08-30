
listen {{if .SslPort}} ssl {.SslPort} {{ else }} {.Port} {{ end }} {.Domains};
listen [::]:{{if .SslPort}} ssl {.SslPort} {{ else }} {.Port} {{ end }} {.Domains};

