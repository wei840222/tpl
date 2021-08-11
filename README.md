# tpl

## How to use
1. prepare values.csv file
```csv
file,hosts,cert,crt,path
aaa,"www.example.com example.com","/web/cert/example.com.cert","/web/key/example.com.key","/a,/b"
aaa,"www.example2.com example2.com","/web/cert/example.com.cert","/web/key/example.com.key","/a"
bbb,"google.com,example.com","/abc","/abc","/a,/b"
```
2. prepare template.txt file
```
server {
    listen 443 ssl;
    server_name {{ .hosts }};
    ssl_certificate {{ .cert }};
    ssl_key {{ .crt }};

    {{- range splitList "," .path }}
    location {{ . }} {
        {{ sha1sum . }}
    }
    {{- end }}
}

```
3. exec
```bash
go run .
```

## Templete syntax
[https://pkg.go.dev/text/template#hdr-Examples](https://pkg.go.dev/text/template#hdr-Examples)

## Templete functions
[https://masterminds.github.io/sprig/](https://masterminds.github.io/sprig/)