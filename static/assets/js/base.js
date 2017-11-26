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

var reshow_context = reshow();
var core = new XMLHttpRequest();
core.open("GET", "/api/v1" + window.location.pathname, true);
core.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
core.setRequestHeader("x-requested-by", "mknote");

core.onreadystatechange=function(){
    if (core.readyState==4 && core.status==200){
        var art_str = core.responseText;
		var art_json = JSON && JSON.parse(art_str) || eval('(' + art_str + ')');
		console.log(art_json.Content);
		var html_text = marked(art_json.Content);
		console.log(html_text);
		document.getElementById("article").innerHTML = html_text;
		var vm = new Vue({
		  el: '#desc',
		  data: {
		    art: art_json
		  }
		})
    }
};

core.send("name=" + reshow_context);


//当按下赞同按钮时，改变它们的颜色
function article_like(artID){
    var never = 'rgb(235, 243, 251)';
    var already = 'rgb(57, 150, 242)';

    var up = document.getElementById("but_"+artID);
    var up_value = document.getElementById("value_"+artID);
    var up_color = up.style.backgroundColor;

    core.onreadystatechange=function(){
        if (core.readyState==4 && core.status==200){
            up.style.backgroundColor=already;
            up_value.innerHTML = parseInt(up_value.innerHTML)+1;
        }
    }

    core.open("PUT","/api/v1/like" + artID, true);
    core.send("name=" + reshow_context);

}
