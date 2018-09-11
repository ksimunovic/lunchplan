package ksimunovic.lunchplan

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.jdbc.Predef._

object Conf {
	var users = System.getProperty("users", "30").toInt
	val baseUrl = System.getProperty("baseUrl", "https://docker:4430")
	var httpConf = http.baseURL(baseUrl)
	var duration = System.getProperty("duration", "240").toInt
}
