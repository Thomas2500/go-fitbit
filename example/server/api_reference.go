package main

import (
	"net/http"
	"strconv"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<head>
	<meta charset="UTF-8" />
	<title>Internal Fitbit data backend</title>
	<link rel="stylesheet" href="/style.css" type="text/css" />
	<link rel="stylesheet" href="https://thomas.bella.network/css/fonts.css" type="text/css" />
</head>
<body>
	<div class="content">
		<h1>Internal Fitbit data backend</h1>
		<a href="/login">&gt;&gt; Fitbit Log In &lt;&lt;</a><br>
		<br>
		Quota: `

	rateLimit := fca.GetRatelimit()
	htmlIndex += strconv.Itoa(rateLimit.RateLimitUsed) + "/" + strconv.Itoa(rateLimit.RateLimitAvailable) + " requests until next reset<br>Next reset: <code>" + rateLimit.RateLimitReset.String() + "</code>"

	htmlIndex += `<br><br><hr>
		<h2>API reference:</h2>
		<ul class="apilist">
			<li>
				<a href="/profile">/profile</a><br>
				Get profile data of user
			</li>
			<li>
				<a href="/devices">/devices</a><br>
				Get a list of devices used by the user.
			</li>
			<li>
				<a href="/food/log">/food/log</a><br>
				Get user food logs. Default is today, can be extended with GET parameters with <code>date=2018-12-15</code> or <code>start=...&end=...</code>
			</li>
			<li>
				<a href="/food/goal">/food/goal</a><br>
				Get user goals for food
			</li>
			<li>
				<a href="/water/log">/water/log</a><br>
				Get user water logs. Default is today, can be extended with GET parameters with <code>date=2018-12-15</code> or <code>start=...&end=...</code>
			</li>
			<li>
				<a href="/water/goal">/water/goal</a><br>
				Get user goals for water
			</li>
			<li>
				<a href="/heart/day">/heart/day</a><br>
				Get the summary of heart rate of a user. Default is data for only today. GET parameters <code>date=</code> and <code>start=...&end=</code> can be given. <code>period=</code> can be <code>1d</code>, <code>7d</code>, <code>30d</code>, <code>1w</code>, <code>1m</code>.
			</li>
			<li>
				<a href="/heart/intraday">/heart/intraday</a><br>
				Get very detailed data about users heart rate. GET parameters <code>date=</code> and <code>start=...&end=</code> can be given. With <code>date=</code> also <code>time-start=00:00&time-end=08:00</code> can be used. <code>detail=1min</code> handles the resolution and can be <code>1sec</code> or <code>1min</code>. Default is today with 1min detail.
			</li>
			<li>
				<a href="/body/weight">/body/weight</a><br>
				Get user weight logs. Default is today. GET parameters <code>date=</code>, <code>start=...&end=...</code> (maximum 31 days) and <code>period=</code> is avaiable. Period can contain <code>1d</code>, <code>7d</code>, <code>30d</code>, <code>1w</code> or <code>1m</code>.
			</li>
			<li>
				<a href="/sleep/log">/sleep/log</a><br>
				Get user sleep logs. Default is today. GET parameters <code>date=</code> or <code>start=...&end=...</code>.
			</li>
			<li>
				<a href="/activities/summary">/activities/summary</a><br>
				Get user activity summary. Default is today. GET parameter available is <code>date=</code> .
			</li>
		</ul>
		<footer>
			Documentation available on GitHub at <a href="https://github.com/Thomas2500/go-fitbit">https://github.com/Thomas2500/go-fitbit</a>
		</footer>
	</div>
</body>
</html>`
	w.Write([]byte(htmlIndex))
}

func handleStyleFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")

	w.Write([]byte(`html, body {
	font-family: "Roboto", sans-serif;
	background-color: #FAFAFA !important;
	width: 100%;
	overflow-x: hidden;
}
a {
	text-decoration: none;
	color: blue;
}
footer {
	font-size: 80%;
	color: gray;
	border-top: 1px solid #d6d6d6;
	font-style: italic;
	padding: 15px 0px 5px;
	text-align: center;
}
footer a {
	color: darkgray;
	text-decoration: underline;
}
.content {
	max-width: 800px;
	margin-left: auto;
	margin-right: auto;
}
.apilist li {
	line-height: 26px;
}
.apilist a {
	background-color: #e0fbff;
	text-decoration: none;
	font-size: 85%;
	padding: 2px 15px;
	border: 1px solid #c9f8ff;
	border-radius: 3px;
	color: #0a5fff;
}
code {
	background-color: #efeded;
	border: 1px solid #e5e5e5;
	padding: 3px;
	border-radius: 3px;
}`))

}
