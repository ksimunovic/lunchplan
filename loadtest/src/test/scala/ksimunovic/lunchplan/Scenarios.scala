package ksimunovic.lunchplan

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.jdbc.Predef._
import scala.concurrent.duration._

object Scenarios {

    val rampUpTimeSecs = 60

	/*
	 *	HTTP scenarios
     */
	 
	val headers_10 = Map("Content-Type" -> """application/json""")


	// Browse
	val browse_guids = csv("accounts.csv").circular
	val scn_Browse = scenario("ApiLogins")
      .during(Conf.duration) {
		feed(browse_guids)
		.exec(
          http("ApiLogin")
			.post("")
			.body(StringBody("""{ "email": "probniEmail", "password": "dugackaSigurnaSifra" }""")).asJSON
            .headers(headers_10)
            .check(status.is(200))
          )
        .pause(1)
    }
}
