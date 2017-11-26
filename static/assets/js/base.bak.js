
Date.prototype.Format = function (fmt) { //author: meizz
    var o = {
        "M+": this.getMonth() + 1, //月份
        "d+": this.getDate(), //日
        "H+": this.getHours(), //小时
        "m+": this.getMinutes(), //分
        "s+": this.getSeconds(), //秒
        "q+": Math.floor((this.getMonth() + 3) / 3), //季度
        "S": this.getMilliseconds() //毫秒
    };
    if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)
    if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return fmt;
}

function preview(id){
    article_local_cache(id);
    window.open("/dashboard/preview");
}

function article_local_cache(){
    var timestamp = new Date().Format("yyyy-MM-dd HH:mm:ss");
    var article_id = document.getElementById('article-id');
    var article_text = document.getElementById('article-editor');
    var article_title = document.getElementById('article-title');
    var article_en_name = document.getElementById('article-en_name');
    var article_tags = document.getElementById('article-tags');
    localStorage.setItem("article_id",article_id.value);
    localStorage.setItem("article_text",article_text.value);
    localStorage.setItem("article_title",article_title.value);
    localStorage.setItem("article_en_name",article_en_name.value);
    localStorage.setItem("article_tags",article_tags.value);
    localStorage.setItem("create_time",timestamp);

    document.getElementById("save-time").innerHTML = "上次保存："+timestamp;

}

function get_article(article_name){
    var article_id = document.getElementById('article-id');
    var article_text = document.getElementById('article-editor');
    var article_title = document.getElementById('article-title');
    var article_en_name = document.getElementById('article-en_name');
    var article_tags = document.getElementById('article-tags');
    var core = new XMLHttpRequest();
    core.onreadystatechange=function(){
        if (core.readyState==4 && core.status==200){
            var article_str = core.responseText;
            var article_json = JSON && JSON.parse(article_str) || eval('(' + article_str + ')');
            article_id.value = article_json.id;
            article_text.value = article_json.content;
            article_title.value = article_json.title;
            article_en_name.value = article_json.en_name;
            article_tags.value = article_json.tags;
            check_article_local_cache();
        }
    }
    core.open("GET", "/api/articles/"+article_name, true);
    core.setRequestHeader("Content-type","application/x-www-form-urlencoded");
    core.setRequestHeader("X-Requested-By","sycki");
    core.send();
}

function check_article_local_cache(){
	var article_id = localStorage.getItem("article_id");
	var article_title = localStorage.getItem("article_title");
	var create_time = localStorage.getItem("create_time");
	if (article_id !== null && article_id == document.getElementById('article-id').value){
		var r=confirm("本地缓存了上次未写完的文章！\n文章：["+article_title+"]\n时间：["+create_time+"]\n是否恢复内容？");
		if (r==true){
			article_local_load();
		}
	}
}

function write_submit(){
//    var data = new FormData(document.getElementById("article-editor-form"));
    var article_id = document.getElementById('article-id');
    var article_text = document.getElementById('article-editor');
    var article_title = document.getElementById('article-title');
    var article_en_name = document.getElementById('article-en_name');
    var article_tags = document.getElementById('article-tags');

    var values = "title=" + article_title.value;
    values += "&id=" + article_id.value;
    values += "&content=" + encodeURIComponent(article_text.value);
    values += "&en_name=" + article_en_name.value;
    values += "&tags=" + article_tags.value;

    var core = new XMLHttpRequest();
    core.onreadystatechange=function(){
        if (core.readyState==4 && core.status==200){
            var article_str = core.responseText;
            document.getElementById("message").innerHTML = article_str;
            if(article_str === "文章已提交！"){
                localStorage.clear();
            }
        }
    }
    core.open("POST", "/dashboard/writer", true);
    core.setRequestHeader("Content-type","application/x-www-form-urlencoded");
    core.send(values);
    return;
}

function article_local_load(){
    var article_id = document.getElementById('article-id');
    var article_text = document.getElementById('article-editor');
    var article_title = document.getElementById('article-title');
    var article_en_name = document.getElementById('article-en_name');
    var article_tags = document.getElementById('article-tags');
    article_id.value = localStorage.getItem("article_id");
    article_text.value = localStorage.getItem("article_text");
    article_title.value = localStorage.getItem("article_title");
    article_en_name.value = localStorage.getItem("article_en_name");
    article_tags.value = localStorage.getItem("article_tags");
}

function getElementsByClassName(n) {
    var classElements = [],allElements = document.getElementsByTagName('*');
    for (var i=0; i< allElements.length; i++ )
   {
       if (allElements[i].className == n ) {
           classElements[classElements.length] = allElements[i];
        }
   }
   return classElements;
}

//解析id指代的元素内的markdown内容
function parser_article(id){
    var article = document.getElementById(id);
    if(article == null)
        return;
    var markdown_text = article.innerHTML.replace(/&gt;/g, ">").replace(/&lt;/g, "<").replace(/&amp;/g, "&");
    var html_text = marked(markdown_text);
    html_text = html_text.replace(/(<h[1-6]{1}.+>)(.+)(<\/h[1-6]{1}>)/g, "<a class='anchor-ids' id='$2'>$1$2$3</a>");
    //html_text = html_text.replace(/(<h[1-6]{1}.+>)(.+)(<\/h[1-6]{1}>)/g, "<a class='idoffset'>id</a><a class='anchor-ids' id='$2'>$1$2$3</a>");
    article.innerHTML = html_text;

    var sections = document.getElementById("sections");
    if(sections == null)
        return;
    var redClassElements = getElementsByClassName('anchor-ids');
    var lists = sections.innerHTML;
    var oldindex = 0;
    for (var i=0; i<redClassElements.length; i++) {
        var ht = redClassElements[i].textContent;
        var hh = redClassElements[i].innerHTML;
        var index = parseInt(hh.substring(2,3));
        if(index > oldindex){
            lists += "<ul><a href='#"+ redClassElements[i].id +"'><li>" + ht + "</li></a>";
            oldindex = index;
        }else if(index == oldindex){
            lists += "<a href='#"+ redClassElements[i].id +"'><li>" + ht + "</li></a>";
            oldindex = index;
        }else if(index < oldindex){
            lists += "</ul><a href='#"+ redClassElements[i].id +"'><li>" + ht + "</li></a>";
            oldindex = index;
        }
    }
    lists += "</ul>";
    sections.innerHTML = lists;
}

function login_data_check(){
    var input_user = document.getElementById("input-user");
    var input_pass = document.getElementById("input-pass");
    var input_code = document.getElementById("input-code");
    var is_pass = true;

    var input_user_r = /^\w{5,20}$/;
    if(input_user_r.test(input_user.value)){
        input_user.style.borderColor="#afa";
    }else{
        input_user.style.borderColor="#faa";
        is_pass = false;
    }

    var input_pass_r = /^\w{5,20}$/;
    if(input_pass_r.test(input_pass.value)){
        input_pass.style.borderColor="#afa";
    }else{
        input_pass.style.borderColor="#faa";
        is_pass = false;
    }

    var input_code_r = /^\w{4,5}$/;
    if(input_code_r.test(input_code.value)){
        input_code.style.borderColor="#afa";
    }else{
        input_code.style.borderColor="#faa";
        is_pass = false;
    }

    return is_pass;
}

//当按下赞同按钮时，改变它们的颜色
function article_desc_thumb_up(buttonid){
    var never = 'rgb(235, 243, 251)';
    var already = 'rgb(57, 150, 242)';

    var up = document.getElementById(buttonid+"_thumb_up");
    var up_value = document.getElementById(buttonid+"_thumb_up_value");
    var up_color = up.style.backgroundColor;

    core.onreadystatechange=function(){
        if (core.readyState==4 && core.status==200 && core.responseText=='success'){
            if (up_color == never || up_color == ''){
                up.style.backgroundColor=already;
                up_value.innerHTML = parseInt(up_value.innerHTML)+1;
            }
            else{
                up.style.backgroundColor=never;
                up_value.innerHTML = parseInt(up_value.innerHTML)-1;
            }
        }
    }

    if (up_color == never || up_color == ''){
        core.open("GET","/article/"+buttonid+"/like/"+reshow(),true);
        core.send();
    }
    else{
        core.open("GET","/article/"+buttonid+"/unlike/"+reshow(),true);
        core.send();
    }

}

//当按下反对按钮时，改变它们的颜色
function article_desc_thumb_down(buttonid){
	var up = document.getElementById(buttonid+"_thumb_down");
	var down = document.getElementById(buttonid+"_thumb_up");

	var up_color = up.style.backgroundColor;
	var down_color = down.style.backgroundColor;

	if (up_color == 'rgb(235, 243, 251)' || up_color == ''){
		up.style.backgroundColor="#3996f2";
		down.style.backgroundColor="#ebf3fb";
	}else{
		up.style.backgroundColor="#ebf3fb";
	}
	
}

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
}

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
}

var core = new XMLHttpRequest();
var reshow_context = reshow();

function init(){
    if(!window.article_id || window.article_id == null){
        core.open("POST", "/reshow", true);
        core.setRequestHeader("Content-type","application/x-www-form-urlencoded");
        core.send("name=" + reshow_context + "&url=" + location.pathname);
    }else{
        core.open("POST", "/article/viewer", true);
        core.setRequestHeader("Content-type","application/x-www-form-urlencoded");
        core.send("name=" + reshow_context + "&articleId="+article_id);
    }

}

window.onload=function(){
    parser_article("article");
	setTimeout("init()",50000);
}