<!doctype html>
<html class="no-js" lang="">
<head>
    <meta charset="utf-8">
    <meta http-equiv="x-ua-compatible" content="ie=edge">
    <title>sShare-simple</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link rel="shortcut icon" href="img/logo.png">

    <link rel="stylesheet" href="css/normalize.css">
    <link rel="stylesheet" href="css/bootstrap.min.css">
    <link rel="stylesheet" href="css/custom.css">
</head>
<body>
<!--[if lte IE 9]>
<p class="browserupgrade">You are using an <strong>outdated</strong> browser. Please <a href="https://browsehappy.com/">upgrade
    your browser</a> to improve your experience and security.</p>
<![endif]-->

<div class="container" style="padding-top: 3rem;">
    <h1>sShare</h1>
    <footer class="blockquote-footer">生活不止眼前的苟，还有身后的苟。</footer>

    <div class="card text-center" style="margin-top: 20px;">
        <div class="card-header">
            获取代理账号 <span class="badge badge-info">服务器状态：<span id="status">0</span> / {{.Limit}}</span>
        </div>
        <div class="card-body">
            <h4 class="card-title">点击验证码，然后确认</h4>
            <div class="card-text" style="padding: 20px;" id="captcha">此处应有验证码</div>
            <button href="#" class="btn btn-primary" id="get">获取</button>
        </div>
        <div class="card-footer text-muted" id="show">
            代理信息显示在这里
        </div>
    </div>
    <p class="fixed-bottom text-dark text-right container">Powered by sShare</p>
</div>

<input type="hidden" id="vtoken">
<script src="js/vendor/modernizr-3.5.0.min.js"></script>
<script src="js/vendor/jquery-3.2.1.min.js"></script>
<script src="js/vendor/bootstrap.bundle.min.js"></script>

<script>
    $("#get")[0].onclick = function (ev) {
        var token = $("#vtoken")[0].innerHTML;
        if (token === "") {
            alert("请完成验证码。");
            return
        }
        $.post("api/new",
                {
                    "token": token
                },
                function (data) {
                    if (data.Status === "ACCEPT") {
                        $("#show")[0].innerHTML = "服务器地址：{{.Location}}" +
                                "<br>端口：" + data.Port +
                                "<br>密码：" + data.Pass +
                                "<br>其他信息";
                    } else {
                        $("#show")[0].innerHTML = "服务器已满，请稍候重试。";
                        setTimeout(function () {
                            window.location = window.location;
                        }, 2000);
                    }
                }, "json"
        ).fail(
                function () {
                    alert("获取账号失败，请稍候重试。");
                }
        );
    };

    function loadServerStatus() {
        $.post("api/count", function (data) {
            $("#status")[0].innerHTML = data
        }).fail(
                function () {
                    alert("无法连接到服务器");
                }
        );
    }

    loadServerStatus();
    setInterval(loadServerStatus, 2000);

    var saveToken = function (token) {
        $("#vtoken")[0].innerHTML = token;
    };
</script>

{{if eq .Captcha.Name "recaptcha"}}
<script src="https://www.recaptcha.net/recaptcha/api.js?onload=onloadCallback&render=explicit" async defer></script>
<script>
    var onloadCallback = function () {
        grecaptcha.render("captcha", {
            'sitekey': '{{.Captcha.SiteID}}',
            "callback": saveToken
        });

    };
</script>
{{else if eq .Captcha.Name "coinhive"}}
<script>
    $("#captcha")[0].append("<div class=\"coinhive-captcha\" \n" +
            "\t\tdata-hashes=\"{{.Captcha.Extra}}\" \n" +
            "\t\tdata-key=\"{{.Captcha.SiteID}}\"\n" +
            "\t\tdata-callback=\"saveToken\"\n" +
            "\t>\n" +
            "\t\t<em>Loading Captcha...<br>\n" +
            "\t\tIf it doesn't load, please disable Adblock!</em>\n" +
            "\t</div>")
</script>
{{else if eq .Captcha.Name "ppoi"}}
<script>
    $("#captcha")[0].append("<div class=\"projectpoi-captcha\" \n" +
            "\t\tdata-hashes=\"{{.Captcha.Extra}}\" \n" +
            "\t\tdata-key=\"{{.Captcha.SiteID}}\"\n" +
            "\t\tdata-callback=\"saveToken\"\n" +
            "\t>\n" +
            "\t\t<em>Loading Captcha...<br>\n" +
            "\t\tIf it doesn't load, please disable Adblock!</em>\n" +
            "\t</div>")
</script>
{{end}}
</body>
</html>
