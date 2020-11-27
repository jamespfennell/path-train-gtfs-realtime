package monitoring

const statusHtmlTemplate = `
<!doctype html>

<html lang="en">
<head>
    <meta charset="utf-8">
    <title>PATH Train GTFS Realtime feed</title>

<style>
body {
  font-family: Helvetica, sans-serif;
  line-height: 1.4em;
}

table{
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

#spacer {
  width: 0px;
  padding: 4px;
}
td.time {
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

table tr td.fail {
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

<body style="font-family: Helvetica, sans-serif; ">



	<table>
  <tr>
  <td></td>
    <td colspan="13" id="razza">Source API</td>
    <td id="spacer"></td>
    <td></td>
    <td id="spacer"></td>
    <td></td>
  </tr>
    <tr style="height: 120px; ">
      <td></td>
      <td class="station">9th&nbsp;St</td>
      <td class="station">14th&nbsp;St</td>
      <td class="station">23rd&nbsp;St</td>
      <td class="station">33rd&nbsp;St</td>
      <td class="station">Christopher&nbsp;St</td>
      <td class="station">Exchange&nbsp;Pl</td>
      <td class="station">Grove&nbsp;St</td>
      <td class="station">Harrison</td>
      <td class="station">Hoboken</td>
      <td class="station">Journal&nbsp;Sq</td>
      <td class="station">Newark</td>
      <td class="station">Newport</td>
      <td class="station">WTC</td>
    <td id="spacer"></td>
      <td style="width: 70px; ">GTFS Builder</td>
    <td id="spacer"></td>
      <td style="width: 70px; ">Total feed update latency</td>
    </tr>

{{range .}}
		<tr>
      <td class="time">{{ .TimeDescription }}</td>
			<td class="success">{{ .StationDescription 10 }}</td>
			<td class="success">{{ .StationDescription 11 }}</td>
			<td class="success">{{ .StationDescription 12 }}</td>
			<td class="success">{{ .StationDescription 13 }}</td>
			<td class="success">{{ .StationDescription 9 }}</td>
			<td class="success">{{ .StationDescription 5 }}</td>
			<td class="success">{{ .StationDescription 4 }}</td>
			<td class="success">{{ .StationDescription 2 }}</td>
			<td class="success">{{ .StationDescription 8 }}</td>
			<td class="success">{{ .StationDescription 3 }}</td>
			<td class="success">{{ .StationDescription 1 }}</td>
			<td class="success">{{ .StationDescription 7 }}</td>
			<td class="success">{{ .StationDescription 6 }}</td>
    <td id="spacer"></td>
			<td class="success">S</td>
    <td id="spacer"></td>
			<td class="success">S</td>
		</tr>
{{end}}
	</table>
  


</body>
</html>
`
