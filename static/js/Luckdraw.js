var winners = new Array()

Array.prototype.indexOf = function(val) {
    for (var i = 0; i < this.length; i++) {
        if (this[i] == val) return i;
    }
    return -1;
};


Array.prototype.remove = function(val) {
    var index = this.indexOf(val);
    if (index > -1) {
        this.splice(index, 1);
    }
};

$("#winners p").each(function(){
    winners.push($(this).html());
});

var xinm_new = new Array()
xinm_new = xinm
var phone_new = new Array()
phone_new = phone

var nametxt = $('.slot');
var phonetxt = $('.name');
var runing = true;
var trigger = true;
var num = 0;
var Lotterynumber
var Count
var pcount


$(function () {
	nametxt.css('background-image','url('+xinm_new[0]+')');
	phonetxt.html(phone_new[0]);
});




// 开始停止
function start(val) {
	if(!val || isNaN(Number(val)) || Number(val) < 1)
	{
		alert("输入错误");
		return false
	}
    for (i=0;i<winners.length;i++){
        temp = phone_new.indexOf(winners[i])
        if (temp > -1)
        {
            phone_new.remove(winners[i])
            xinm_new.remove(xinm_new[temp])
        }
    }
    pcount = xinm_new.length-1;//参加人数
    Lotterynumber = val
    Count = val
	if (runing) {
		if ( pcount <= Lotterynumber ) {
			alert("抽奖人数不足");
		}else{
			runing = false;
			$('#start').text('停止');
			startNum()
		}
	} else {
		$('#start').text('自动抽取中('+ Lotterynumber+')');
		zd();
	}
}

// 开始抽奖

function startLuck() {
	runing = false;
	$('#btntxt').removeClass('start').addClass('stop');
	startNum()
}

// 循环参加名单
function startNum() {
	num = Math.floor(Math.random() * pcount);
	nametxt.css('background-image','url('+xinm_new[num]+')');
	phonetxt.html(phone_new[num]);
	t = setTimeout(startNum, 0);
}

// 停止跳动
function stop() {
	pcount = xinm_new.length-1;
	clearInterval(t);
	t = 0;
}

// 打印中奖人

function zd() {
	if (trigger) {

		trigger = false;
		var i = 0;
		var last = Count;

		if ( pcount >= Lotterynumber ) {
			stopTime = window.setInterval(function () {
				if (runing) {
					runing = false;
					$('#btntxt').removeClass('start').addClass('stop');
					startNum();
				} else {
					runing = true;
					$('#btntxt').removeClass('stop').addClass('start');
					stop();

					i++;
					Lotterynumber--;

					$('#start').text('自动抽取中('+ Lotterynumber+')');

					if ( i == Count ) {
						console.log("抽奖结束");
						window.clearInterval(stopTime);
						$('#start').text("开始");
						Lotterynumber = Count;
						trigger = true;
					};

					$('.luck-user-list').prepend("<li><div class='portrait' style='background-image:url("+xinm_new[num]+")'></div><div class='luckuserName'>"+phone_new[num]+"</div></li>");
                    // if (i == last){
                    //     $('.luck-user-list').prepend("<li>__________________________</li>"); //----------------------------------
                    // }
					//将已中奖者从数组中"删除",防止二次中奖
					xinm_new.splice($.inArray(xinm_new[num], xinm_new), 1);
					phone_new.splice($.inArray(phone_new[num], phone_new), 1);

				}
			},1000);
		};
	}
}

