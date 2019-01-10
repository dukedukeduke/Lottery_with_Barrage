<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Lottery</title>
    <link rel="stylesheet" href="/static/css/bootstrap.min.css">
    <link rel="stylesheet" href="/static/css/bootstrap-theme.min.css">
    <link rel="stylesheet" href="/static/css/style.css">
    <link rel="stylesheet" href="/static/css/barrager.css">

    <script src="/static/js/jquery-2.2.1.min.js"></script>
    <script src="/static/js/bootstrap.js"></script>
    <script src="/static/js/header.js"></script>
    <script src="/static/js/jquery.barrager.js"></script>
</head>
<body>
    <div class='luck-back'>
        <span id="conn_status">connecting...</span>
        <div class="lucky_list_show">
            <div class="panel-group">
                <div class="panel panel-default">
                    <div class="panel-heading">
                        <h4 class="panel-title">
                            <a data-toggle="collapse" href="#collapse1"><span class="glyphicon glyphicon-plus"></span>Winners</a>
                        </h4>
                    </div>
                    <div id="collapse1" class="panel-collapse collapse">
                        <div class="panel-body" id="winners">
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="luck-content ce-pack-end">
            <div id="luckuser" class="slotMachine">
                <div class="slot">
                    <span class="name">姓名</span>
                </div>
            </div>
            <div class="luck-content-btn">
                <a id="start" class="start" onclick='start($("#count_set").val())'>开始</a>
                <span id="count_set_span"><input id="count_set" ></span>
            </div>
            <div class="luck-user">
                <div class="luck-user-title">
                    <span>中奖名单</span>
                </div>
                <ul class="luck-user-list" id="lucky_list"></ul>
                <div class="luck-user-btn">
                    <a onclick="saveFile()" id="save_lottery">中奖人</a>
                </div>
            </div>
        </div>
    </div>


    <script>
        var all_employee = new Array();
        var xinm = new Array();
        var phone = new Array();
        all_employee = {{.Employee}}
        winners = {{.Winner}}
        if (winners){
            for (i=0;i<winners.length;i++)
                {
                    var winner_item=$("<p>"+winners[i]+"</p>");
                    winner_item.appendTo($("#winners"));
                }
        }


        for (i=0;i<all_employee.length;i++)
        {
            xinm[i] = all_employee[i].Icon
            phone[i] = all_employee[i].Name
        }
        let sock = null;
        const wsuri = "ws://127.0.0.1:1234/bubble";
        var heart_beat = 0;
        var need_heartbeat = false;
        window.onload = conn_server();

        function conn_server(){
            console.log("onload successfully");
            sock = new WebSocket(wsuri);
            sock.onopen = function() {
                document.getElementById('conn_status').innerHTML="Online";
                $("#conn_status").css("color", "black")
                console.log("connected to " + wsuri);
                need_heartbeat = true;
                setInterval(show, 9000);
            }
            sock.onclose = function(e) {
                console.log("connection closed (" + e.code + ")");
                document.getElementById('conn_status').innerHTML="Offline";
                $("#conn_status").css("color", "white");
                $("#conn_status").attr("onclick", "conn_server()");
                need_heartbeat = false;

            }

            var msg_list = new Array();
            sock.onmessage = function(e) {
                // console.log("message received: " + e.data);
                // document.getElementById('recv_msg').innerHTML=e.data;
                msg_list.push(e.data);
                while (msg_list.length > 0)
                    {
                        msg = msg_list.pop();
                        console.log(msg);
                        ico_list = ["/static/img/icon/blueE.ico", "/static/img/icon/blueN.ico", "/static/img/icon/bubble.ico", "/static/img/icon/cheers.ico",
                        "/static/img/icon/colonA.ico", "/static/img/icon/fiveX.ico", "/static/img/icon/keyboardE.ico", "/static/img/icon/party.ico",
                        "/static/img/icon/redA.ico", "/static/img/icon/redE.ico", "/static/img/icon/redN.ico", "/static/img/icon/X storage.ico", "/static/img/icon/blueX.ico"]
                        icon = ico_list[Math.floor(Math.random()*12)]
                        console.log(icon)
                        if (msg.indexOf("2018") > -1 || msg.indexOf("新年") > -1)
                            {
                                console.log("match")
                                $(".barrage").css("width", "500px");
                                $(".barrage_box").css("height", "100px");
                                $(".barrage_box div.p a").css("font-size", "25px");
                                $(".barrage_box div.p").css("margin", "9% 0px");
                                console.log(e.data.length);
                                if (e.data.length > 16)
                                {
                                    console.log("too long")

                                    danmu(e.data.slice(0, 16), icon);
                                    danmu("程序员：太长被我截尾了", icon);
                                }
                                else{
                                    danmu(e.data, icon);
                                }

                            }
                            else
                            {
                                console.log("not match")
                                $(".barrage").css("width", "500px");
                                $(".barrage_box").css("height", "40px");
                                $(".barrage_box div.p a").css("font-size", "14px");
                                $(".barrage_box div.p").css("margin", "0px 0px 0px 0px");
                                if (e.data.length > 16)
                                {
                                    console.log("too long")
                                    danmu(e.data.slice(0, 16), icon);
                                    danmu("程序员：太长被我截尾了", icon);
                                }
                                else{
                                    danmu(e.data, icon);
                                }
                            }
                    }

            }
        }

        var saveFile = function(event) {
            members = new Array();
            $("#lucky_list li").each(function(){
                members.push($(this).children().first().next().html());
                });
            $("#save_lottery").attr("href",'data:application/json;charset=utf-8;json,' + members);
            $("#save_lottery").attr("download","data.json");
            $.get("/save_lottery/?list="+members, function(result){
                status = result.split(":")[0]
                if (status == "success"){
                    tmp_list = result.split(":")[1].split(",")
                    $("#winners p").each(function(){
                        $(this).remove();
                        });
                    if (tmp_list){
                        for (i=0;i<tmp_list.length;i++)
                            {
                                var winner_item=$("<p>"+tmp_list[i]+"</p>");
                                winner_item.appendTo($("#winners"));
                            }
                    }
                    alert("Save successfully");
                }else{
                    alert("Save failed");
                }
            })
        };


        function show(){
            if (need_heartbeat === true)
            {
                heart_beat = heart_beat + 1;
                console.log("Heartbeat send", heart_beat);
                if (sock.readyState !== 1){
                    console.log("Connection Error");
                    if (document.getElementById('conn_status').innerHTML==="Online"){
                        document.getElementById('conn_status').innerHTML="Offline";
                    }
                }
                else{
                    sock.send("2018");
                }
            }
        }


        // 弹幕
        function danmu(msg, icon) {
            var item={
                img: icon, //图片
                info:msg, //文字
                href:'http://www.baidu.com', //链接
                close:true, //显示关闭按钮
                speed:10, //延迟,单位秒,默认6
                color:'#ffffff', //颜色,默认白色
                bottom: Math.random() * 800,
                old_ie_color:'#000000', //ie低版兼容色,不能与网页背景相同,默认黑色
            }
            $('body').barrager(item,0);

        }
    </script>
    <script src="/static/js/Luckdraw.js" type="text/javascript"></script>
</body>
</html>
