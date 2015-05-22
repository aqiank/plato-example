$(document).ready(function() {
	// modal
	$(".modal-trigger").leanModal();

	// knob
	$(".knob").knob();

	// owl carousel
	$(".owl-carousel").owlCarousel({
		items: 3,
		autoHeight: true
	});

	// select
	$("select").material_select();

	// sidenav
	$(".button-collapse").sideNav();

	// dropdown
	$(".dropdown-button").dropdown({
		inDuration: 300,
		outDuration: 225,
		constrain_width: false,
		hover: true,
		gutter: 0,
		belowOrigin: true,
	});

	// datepicker
	$(".datepicker").pickadate({
		selectMonths: true,
		selectYears: 5,
	});

	// projects
	$(".oi-projects").owlCarousel({
		items: 3,
		autoplay: true,
		autoplayTimeout: 10000,
		responsive: {
			0: {
				items: 1
			},
			800: {
				items: 2
			},
			1000: {
				items: 3
			}
		}
	});

	$("#modal-sign-in").submit(function(e) {
		var form = $(this);
		$.ajax({
			url: form.attr("action"),
			method: "POST",
			data: form.serialize(),
			dataType: "json",
		}).done(function(resp) {
			Materialize.toast("Successfully logged in!", 1000, "green");
			window.location = "/";
		}).fail(function(resp) {
			Materialize.toast(resp.responseText, 1000, "red");
		});
		e.preventDefault();
	});

	$("#modal-sign-up").submit(function(e) {
		var form = $(this);
		$.ajax({
			url: form.attr("action"),
			method: "POST",
			data: form.serialize(),
			dataType: "json",
		}).done(function(resp) {
			Materialize.toast("Sent a verification code to your email!", 1000, "green");
			window.location = "/";
		}).fail(function(resp) {
			Materialize.toast(resp.responseText, 1000, "red");
		});
		e.preventDefault();
	});
});

function setImageInputPreview(input, preview, uploadURL, success) {
	$(input).change(function(){
		if (this.files && this.files[0]) { var reader = new FileReader();
			var file = this.files[0];
			reader.onload = function (e) {
				$(preview).attr("src", e.target.result);
				if (typeof uploadURL === "string") {
					var data = new FormData();
					data.append("avatar", file);
					data.append("what", "avatar");
					$.ajax({
						url: uploadURL,
						data: data,
						cache: false,
						contentType: false,
						processData: false,
						type: "POST",
						success: function(data) {
							if (typeof success !== "undefined") {
								success(data);
							}
						}
					});
				}
			}
			reader.readAsDataURL(file);
		}
	});
}

