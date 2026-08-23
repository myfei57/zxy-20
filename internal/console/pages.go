package console

import "fmt"

const shell = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
body{font-family:system-ui,sans-serif;margin:0;background:#f4f6f8;color:#222}
nav{background:#123;color:#fff;padding:10px 18px;display:flex;gap:18px;align-items:center}
nav a{color:#cfe3ff;text-decoration:none;font-size:14px}
nav a:hover{color:#fff}
main{max-width:1080px;margin:22px auto;padding:0 16px}
h1{font-size:20px}
table{width:100%%;border-collapse:collapse;background:#fff;box-shadow:0 1px 3px rgba(0,0,0,.12)}
th,td{padding:8px 10px;border-bottom:1px solid #e3e8ee;text-align:left;font-size:13px}
th{background:#eef2f6}
button{background:#1a5fb4;color:#fff;border:0;padding:6px 12px;border-radius:4px;cursor:pointer}
input{padding:5px 8px;border:1px solid #c7d0d9;border-radius:4px}
.error{color:#b00020}
</style>
</head>
<body>
<nav><strong>BMS 控制台</strong><a href="/rooms">房间</a><a href="/plans">计划</a><a href="/meters">计量</a><a href="/audit">审计</a></nav>
<main><h1>%s</h1><div id="app">加载中…</div></main>
<script>%s</script>
</body>
</html>`

var pages = map[string]string{
	"rooms": renderPage("房间管理", `
async function load(){const res=await fetch('/api/rooms');const data=await res.json();
let html='<table><tr><th>房间</th><th>设定温度</th><th>模式</th><th>绑定计划</th><th>设备状态</th></tr>';
data.forEach(v=>{html+='<tr><td>'+(v.room.name||v.room.id)+'</td><td>'+v.room.setpoint+'</td><td>'+v.room.mode+'</td><td>'+(v.room.bound_plan_id||'-')+'</td><td>'+v.device_state+'</td></tr>';});
html+='</table>';document.getElementById('app').innerHTML=html;}
load();`),
	"plans": renderPage("运行计划", `
async function load(){const res=await fetch('/api/plans');const data=await res.json();
let html='<table><tr><th>计划</th><th>楼宇</th><th>版本</th><th>状态</th><th>下发游标</th><th>操作</th></tr>';
data.forEach(p=>{html+='<tr><td>'+(p.name||p.id)+'</td><td>'+p.building_id+'</td><td>'+p.version+'</td><td>'+(p.active?'生效':'草稿')+'</td><td>'+p.cursor+'</td><td><button onclick="go(\''+p.id+'\')">切换</button></td></tr>';});
html+='</table>';document.getElementById('app').innerHTML=html;}
async function go(id){await fetch('/api/plans/'+id+'/switch',{method:'POST'});load();}
load();`),
	"meters": renderPage("能耗计量", `
async function load(){const res=await fetch('/api/meters');const data=await res.json();
let html='<table><tr><th>表计</th><th>房间</th><th>游标</th><th>最新读数</th><th>操作</th></tr>';
data.forEach(m=>{html+='<tr><td>'+(m.name||m.id)+'</td><td>'+m.room_id+'</td><td>'+m.cursor+'</td><td>'+m.last_value+'</td><td><button onclick="read(\''+m.id+'\')">采集</button></td></tr>';});
html+='</table>';document.getElementById('app').innerHTML=html;}
async function read(id){await fetch('/api/meters/'+id+'/read',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({value:Math.random()*10+5})});load();}
load();`),
	"audit": renderPage("审计日志", `
async function load(){const res=await fetch('/api/audit');const data=await res.json();
let html='<table><tr><th>时间</th><th>类型</th><th>对象</th><th>结果</th><th>说明</th></tr>';
data.forEach(e=>{html+='<tr><td>'+e.occurred_at+'</td><td>'+e.kind+'</td><td>'+e.entity_id+'</td><td>'+e.result+'</td><td>'+(e.message||'')+'</td></tr>';});
html+='</table>';document.getElementById('app').innerHTML=html;}
load();`),
}

func renderPage(title, script string) string {
	return fmt.Sprintf(shell, title, title, script)
}
