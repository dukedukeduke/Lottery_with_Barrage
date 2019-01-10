<html>
<head>
<meta charset="utf-8">
<title>Message</title>
<link rel="stylesheet" href="/static/css/bootstrap.min.css">
<link rel="stylesheet" href="/static/css/bootstrap-theme.min.css">
<link rel="stylesheet" href="/static/css/barrager.css">

<script src="/static/js/jquery-2.2.1.min.js"></script>
<script src="/static/js/jquery.barrager.js"></script>
</head>
<body style="background-image: url(/static/img/back.jpg); background-size: cover; background-position: top center; position: absolute; height: 100%; width: 100%; color: #fff;">
<div class="container">
    <div style="border: 1px;border-radius:10%;position: absolute; top: 20%; bottom:20%; left:20%; right:20%; background: rgba(0,0,0,.6); padding: 20px; border-radius: 5px; padding-right:260px;">
        <form action="/message/" method="post">
            <div>
                <div>
                    <input type="text" name="message" id="message" style="position: absolute;left:20%;top: 20%;width:40%;color:black;height:10%;border-radius:5px;font-size:30px">
                </div>
                <div>
                    <input type="submit" value="send" style="position:absolute;top:20%;left:63%;background: #f29807; width:100px; text-decoration: none; display: inline-block; color: #fff; text-align: center; margin: 0 10px; cursor: pointer;height:10%;border-radius:5px;">
                </div>
                <div style="border:1px;position:absolute;top:40%;left:10%;text-align: center; margin: 0 10px; cursor: pointer;width:80%;height:50%" id="danmu_message">

                </div>
            </div>
        </form>
    </div>
</div>
<script>
    window.setInterval(function(){if ($("#message").val()){danmu($("#message").val());
    }},3000);

    // 弹幕
    function danmu(msg) {
        var item={
            info:msg, //文字
            speed:10, //延迟,单位秒,默认6
            color:'#ffffff', //颜色,默认白色
            bottom: Math.random() * $(window).height()*0.5,
            old_ie_color:'#000000', //ie低版兼容色,不能与网页背景相同,默认黑色
        }
        $('body').barrager(item, 0);
    }
</script>
</body>
</html>