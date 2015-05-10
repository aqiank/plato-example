$(document).ready(function() {
        $("#oi-quotes").slick({
        	dots: false,
        	arrows: false,
        	infinite: true,
        	slidesToShow: 1,
        	adaptiveHeight: true,
        	autoplay: true,
        	autoplaySpeed: 10000,
        	fade: true,
        	cssEase: "linear",
        });

	$("#oi-recommended-projects").slick({
        	dots: false,
        	arrows: false,
        	infinite: true,
        	slidesToShow: 1,
        	adaptiveHeight: true,
        	autoplay: true,
        	autoplaySpeed: 5000,
        	fade: true,
        	cssEase: "linear",
	});

	$(".google.plus.button").click(function() {
		// signInCallback defined in step 6.
		auth2.grantOfflineAccess({"redirect_uri": "postmessage"}).then(signInCallback);
	});
});

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
