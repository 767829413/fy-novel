package definition

const (
	NovelTemp_EPUB = `    <h2>{{ .Title }}</h2>
	{{ .Content }}`

	NovelTemp_EPUB_OLD = `<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="zh-CN">
<head>
  <meta http-equiv="Content-Type" content="application/xhtml+xml; charset=utf-8" />
  <title>{{ .Title }}</title>
  <style type="text/css">
    p {
      text-indent: 2em;
      letter-spacing: 1px;
    }
  </style>
</head>
<body>
  <h2>{{ .Title }}</h2>
  <div>
    {{ .Content }}
  </div>
</body>
</html>
`
	NovelTemp_HTML = `<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <title>{{ .Title }}</title>
  <meta charset="UTF-8">
  <link href="https://cdn.staticfile.net/bootstrap/5.3.2/css/bootstrap.min.css" rel="stylesheet">
  <style type="text/css"/>
    body {
      max-width: 800px;
      margin: 100px auto;
      background: #111;
    }

    h1 {
      color: #939392;
    }

    p {
      text-indent: 2em;
      letter-spacing: 0.1em;
      color: #939392;
      font-size: 25px;
      margin: 40px 0;
    }

    .bottom-bar {
      margin-top: 100px;
    }
  </style>
</head>

<body>
  <h1>{{ .Title }}</h1>
  <div class="content">
    {{ .Content }}
  </div>
  <div class="bottom-bar d-flex justify-content-between">
    <button id="btn-pre" type="button" class="btn btn-primary btn-lg w-25" onclick="turnPage('pre')">上一页</button>
    <button id="btn-next" type="button" class="btn btn-primary btn-lg w-25" onclick="turnPage('next')">下一页</button>
  </div>
</body>

<script src="https://cdn.staticfile.net/bootstrap/5.3.2/js/bootstrap.min.js"></script>
<script type="text/javascript">
  // 当前页面 URL
  const url = decodeURI(location.href)
  const prefixUrl = url.substring(0, url.lastIndexOf('/') + 1)

  // 获取当前文件名
  let index = url.match(/(\d+)\_.html/)[1]
  console.log('当前页索引', index)

  if (index <= 1) {
    document.getElementById('btn-pre').disabled = true
  }

  function turnPage(action) {
    action === 'next' ? ++index : --index
    if (index < 1) {
      index = 1
      return
    }
    location.href = prefixUrl + index + '_.html'
  }

  // 监听键盘按下事件
  document.addEventListener('keyup', function (e) {
    // 检查按下的是左箭头键还是右箭头键
    switch (e.key) {
      case 'ArrowRight':
        turnPage('next')
        break
      case 'ArrowLeft':
        turnPage('pre')
        break
      default:
        break
    }
  })
</script>

</html>

`
)
