package handler

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func initIndex() error {
	f := filepath.Join(AppRoot(), "index.html")
	_, err := os.Stat(f)
	if err == nil {
		return nil
	}
	fi, err := os.OpenFile(f, os.O_CREATE|os.O_TRUNC|os.O_APPEND|syscall.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_ = fi.Truncate(0)
	_, err = fi.WriteString(strings.ReplaceAll(htmlTemplate, "__WEB_PATH__", AppRoot()))
	if err != nil {
		return err
	}
	return nil
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>fs-server</title>
    <style type="text/css">
        html, body, .c {
            height: 100%;
        }

        .c {
            display: flex;
            justify-content: center;
            align-items: center;
            flex-direction: column;

            font-family: "Microsoft JhengHei UI";
            text-align: center;
        }

        .c .title {
            font-size: 86px;
            margin-bottom: 20px;
        }

        .c .info {
            color: #6e6e6e;
            font-size: 38px;
        }

        .c .info span {
            border-bottom: 2px solid #a1a1a1;
            padding: 2px;
            cursor: copy;
        }
    </style>
</head>
<body>

<div class="c">
    <div class="title">你好, Hello</div>
    <div class="info">file://<span>__WEB_PATH__</span></div>
</div>

</body>
</html>

`
