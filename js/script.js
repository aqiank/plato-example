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

