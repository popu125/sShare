<!doctype html>
<html class="no-js" lang="">
<head>
    <meta charset="utf-8">
    <meta http-equiv="x-ua-compatible" content="ie=edge">
    <title>sShare-simple</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link rel="shortcut icon" href="img/logo.png">

    <link rel="stylesheet" href="https://cdn.bootcss.com/normalize/7.0.0/normalize.min.css">
    <link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/4.0.0-beta.2/css/bootstrap.min.css">
</head>
<body>
<!--[if lte IE 9]>
<p class="browserupgrade">You are using an <strong>outdated</strong> browser. Please <a href="https://browsehappy.com/">upgrade
    your browser</a> to improve your experience and security.</p>
<![endif]-->
<a href="https://github.com/popu125/sShare"><img style="position: absolute; top: 0; right: 0; border: 0;" src="https://camo.githubusercontent.com/38ef81f8aca64bb9a64448d0d70f1308ef5341ab/68747470733a2f2f73332e616d617a6f6e6177732e636f6d2f6769746875622f726962626f6e732f666f726b6d655f72696768745f6461726b626c75655f3132313632312e706e67" alt="Fork me on GitHub" data-canonical-src="https://s3.amazonaws.com/github/ribbons/forkme_right_darkblue_121621.png"></a>

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
            <button href="#" class="btn btn-primary" id="get" onclick="infoGet()">获取</button>
        </div>
        <div class="card-footer text-muted" id="show">
            代理信息显示在这里
        </div>
    </div>
    <p class="fixed-bottom text-dark text-right container">Powered by sShare</p>
</div>

<input type="hidden" id="vtoken">

<script src="https://cdn.bootcss.com/jquery/3.2.1/jquery.min.js"></script>
<script src="https://cdn.bootcss.com/bootstrap/4.0.0-beta.2/js/bootstrap.min.js"></script>

<script>
    function infoGet(ev) {
        var token = $("#vtoken")[0].value;
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
        $("#vtoken")[0].value = token;
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
{{else}}
<script>
	$("#vtoken")[0].value = "base";
</script>
{{end}}
</body>
</html>
