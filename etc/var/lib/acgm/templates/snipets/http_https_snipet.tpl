# Redirect HTTP requests to HTTPS
server {
    listen { .Port };
    server_name { .Domains };
    # Redirect all port { .Port } (HTTP) requests to port { .SslPort } (HTTPS).
    return 301 https://$host:{ .SslPort }$request_uri;
}