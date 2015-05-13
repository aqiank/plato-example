$(document).ready(function() {
        $("#oi-quotes").owlCarousel({
		singleItem: true,
		autoPlay: 10000,
		pagination: false
        });

	$("#oi-recommended-projects").owlCarousel({
		singleItem: true,
		autoPlay: 10000
	});

	$(".owl-carousel").owlCarousel();

	$(".google.plus.button").click(function() {
		// signInCallback defined in step 6.
		auth2.grantOfflineAccess({"redirect_uri": "postmessage"}).then(signInCallback);
	});

	$(".facebook.button").click(function() {
		FB.login(function(response) {
			console.log(response);
			$.post("/login",
				{accessToken: response.accessToken},
				function(resp) {
					window.location = "/";
				}
			);
		}, {scope: "public_profile,email"});
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
	}
}
