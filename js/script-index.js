$(document).ready(function() {
        $("#oi-quotes").owlCarousel({
		items: 1,
		autoplay: true,
		autoplayTimeout: 10000,
		dots: false,
        });

	$("#oi-recommended-projects").owlCarousel({
		items: 1,
		autoplay: true,
		autoplayTimeout: 10000,
	});

	$(".google").click(function(e) {
		// signInCallback defined in step 6.
		auth2.grantOfflineAccess({"redirect_uri": "postmessage"}).then(signInCallback);
		e.preventDefault();
	});

	$(".facebook").click(function(e) {
		FB.login(function(response) {
			console.log(response);
			if (response.status == "connected") {
				$.post("/login",
					{
						accessToken: response.authResponse.accessToken,
						loginFrom: "facebook"
					},
					function(resp) {
						window.location = "/";
					}
				);
			} else {
				console.log("Failed to sign into Facebook");
			}
		}, {scope: "public_profile,email"});
		e.preventDefault();
	});
});


// Google Sign In callback
function signInCallback(authResult) {
	console.log("Result: " + authResult);
	if (authResult["code"]) {
		// Hide the sign-in button now that the user is authorized, for example:
		$("#signinButton").attr("style", "display: none");

		// Send the code to the server
		$.post("/login",
			{authCode: authResult["code"], loginFrom: "google"},
			function(result) {
				// Handle or verify the server response.
				window.location = "/";
			}
		);
	} else {
		// There was an error.
		console.log("Failed to sign into Google");
	}
}
