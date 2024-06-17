## This is a sample go project that implements 
- http server with tls encryption
- embedded static content from a folder (created by react build)

## Command to Create SSL Certificate
`$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./nginx.key -out ./nginx.crt`

### Commands to work with pfx certificate format
[Converting pfx to pem using openssl](https://stackoverflow.com/questions/15413646/converting-pfx-to-pem-using-openssl)
