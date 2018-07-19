function bin2hex (s) {
  // From: http://phpjs.org/functions
  // +   original by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
  // +   bugfixed by: Onno Marsman
  // +   bugfixed by: Linuxworld
  // +   improved by: ntoniazzi (http://phpjs.org/functions/bin2hex:361#comment_177616)
  // *     example 1: bin2hex('Kev');
  // *     returns 1: '4b6576'
  // *     example 2: bin2hex(String.fromCharCode(0x00));
  // *     returns 2: '00'

  var i, l, o = "", n;

  s += "";

  for (i = 0, l = s.length; i < l; i++) {
    n = s.charCodeAt(i).toString(16)
    o += n.length < 2 ? "0" + n : n;
  }

  return o;
};

function reshow(){
    var canvas = document.createElement('canvas');
    var ctx = canvas.getContext('2d');
    var txt = 'blog.giepin.com';
    ctx.textBaseline = "top";
    ctx.font = "14px 'Arial'";
    ctx.textBaseline = "tencent";
    ctx.fillStyle = "#f60";
    ctx.fillRect(125,1,62,20);
    ctx.fillStyle = "#069";
    ctx.fillText(txt, 2, 15);
    ctx.fillStyle = "rgba(102, 104, 0, 0.7)";
    ctx.fillText(txt, 4, 17);

    var b64 = canvas.toDataURL().replace("data:image/png;base64,","");
    var bin = atob(b64);
    var crc = bin2hex(bin.slice(-16,-12));
    return crc;
    //localStorage.setItem("article_show",crc);
};

function apitor(method, url, callback){
    core.open(method,url, true);
    core.onreadystatechange=callback
	core.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	// "Content-Type: application/json"
	core.setRequestHeader("x-requested-by", "mknote");
    core.send("name=" + reshow_context);
}

//当按下赞同按钮时，改变它们的颜色
function article_like(artID){
    var never = 'rgb(235, 243, 251)';
    var already = 'rgb(57, 150, 242)';

    var up = document.getElementById("but_"+artID);
    var up_value = document.getElementById("value_"+artID);
    var up_color = up.style.backgroundColor;
	
	if (up_color == 'rgb(57, 150, 242)'){
		return
	}
	
	apitor("POST","/v1/like/" + artID, function(){
        if (core.readyState==4 && core.status==200){
            up.style.backgroundColor=already;
            up_value.innerHTML = parseInt(up_value.innerHTML)+1;
        }
    });
}

function visit(){
    artID = window.location.pathname
	apitor("POST", "/v1/visit" + window.location.pathname, function(){});
}

//function generate_index(){
//	apitor("GET", "/v1/index", function(){
//	    if (core.readyState==4 && core.status==200){
//	        var str = core.responseText;
//			var json = JSON && JSON.parse(str) || eval('(' + str + ')');
//			var sections = new Vue({
//			  el: '#sections',
//			  data: {
//			    tags: json
//			  }
//			});
//	    }
//	});
//}

function click_tag(tag){
	var list = document.getElementById("list_"+tag);
	var cli_stat = document.getElementById("cli_"+tag);
	if(list.style.display == 'inline-block'){
		list.style.display = "none";
		cli_stat.className = "cli_close";
	}else{
		list.style.display = "inline-block";
		cli_stat.className = "cli_open";
	}
}

function but_menu(){
	var list = document.getElementById("sections");
	var but = document.getElementById("but_menu");
	var t = but.textContent;
	if(t == ">"){
		list.style.left = "0px";
		but.style.left = "190px";
		but.textContent = "<";
	}else{
		list.style.left = "-300px";
		but.style.left = "10px";
		but.textContent = ">";
	}
}

var reshow_context = reshow();
var core = new XMLHttpRequest();


window.onload=function(){
	setTimeout("visit()",50000);
}
