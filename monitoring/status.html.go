package monitoring

const statusHtmlTemplate = `
<!doctype html>

<html lang="en">
<head>
    <meta charset="utf-8">
    <title>PATH Train GTFS Realtime feed</title>

<style>
html {

}
body {
  font-family: Helvetica, sans-serif;
  line-height: 1.4em;
}

h1, h2 {
  text-align: center;
}

table{
	margin: 0px auto;
  color: #666;
  text-align: center;
}

table td {
  padding: 0px 10px;
  margin: 0px;
}

#razza {
  color: #666;
  border-bottom: 0px solid #666;
  text-align: center;
}

.hover {
	cursor: help;
text-decoration:underline;
text-decoration-style: dotted;
}

#spacer {
  width: 0px;
  padding: 4px;
}
td.time {
	padding: 4px; 
  color: #666;
}
td.station {
  padding: 0px 10px 10px 10px;
  writing-mode: vertical-rl;
text-orientation: mixed;
text-align: center;
/*
transform: rotate(90deg);*/
  vertical-align: bottom;
  text-align: right;
  color: #666;
  width: 20px;
}
table tr td.success {
  height: 40px;
  background-color: #39aa56;
  vertical-align: middle;
  text-align: center;
  color: white;
}

table tr td.failure {
  background-color: #db4545;
  vertical-align: middle;
  text-align: center;
  color: white;
}

/*
 red #db4545;
*/
</style>
</head>

<body>

	<h1>PATH GTFS Realtime Status Page</h1>
	<h2>History</h2>
	<table>
  <tr>
  <td></td>
    <td colspan="13" id="razza"><h3>Source API</h3>
	Number of stop time updates retrieved or<br />
	F if the data retrieval for the station failed
    <td id="spacer"></td>
    <td></td>
    <td id="spacer"></td>
    <td></td>
  </tr>
    <tr style="height: 120px; ">
      <td></td>
{{range .StationNames }}
    <td class="station">{{.}}</td>
{{end}}
    <td id="spacer"></td>
      <td style="width: 70px; ">GTFS Builder</td>
    <td id="spacer"></td>
      <td style="width: 70px; ">Time elapsed since last successful feed update</td>
    </tr>

{{ $updates := .Updates }}
{{ $stationsIDs := .StationIDs }}
{{range $u := .Updates}}

		<tr>
      <td class="time">{{ $u.TimeDescription }}</td>
{{range $s := $stationsIDs}}

			<td class="{{ $u.StationClass $s }}">{{ $u.StationDescription $s }}</td>
{{end}}
    <td id="spacer"></td>
			<td class="{{ $u.BuilderClass }}">{{ $u.BuilderDescription }}</td>
    <td id="spacer"></td>
			<td class="{{ $u.LatencyClass }}">{{ $u.LatencyDescription }}</td>
		</tr>
{{end}}
	</table>
  


</body>
</html>
`
