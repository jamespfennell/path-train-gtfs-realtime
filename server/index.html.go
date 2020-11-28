package server

const indexHTMLPage = `
<!doctype html>

<html lang="en">
<head>
    <meta charset="utf-8">
    <title>PATH Train GTFS Realtime</title>

<style>
body {
  font-family: Helvetica, sans-serif;
  line-height: 1.4em;
}

h1 {
  text-align: center;
}
ul {
	font-size: 1.2em;
}

li {
	margin: 8px; 
}

</style>
</head>

<body>
<div style="width: 600px; margin: 10px auto;">
	<h1>PATH Train GTFS Realtime</h1>
	<ul>
		<li><a href="./feed/">Data feed</a></li>
		<li><a href="./status/">Status page</a></li>
		<li><a href="./status/json/">JSON status blob</a></li>
		<li><a href="https://github.com/jamespfennell/path-train-gtfs-realtime/">Github repository</a></li>
	</ul>
</div>
</body>
</html>
`
