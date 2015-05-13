var SignInRules =
{
        email: {
                identifier: "email",
                rules: [
                        {
                                type: "length[" + EmailLength + "]",
                                prompt: "Email must be at least " + EmailLength + " characters long"
                        }
                ]
        },
        password: {
                identifier: "password",
                rules: [
                        {
                                type: "length[" + PasswordLength + "]",
                                prompt: "Email must be at least " + PasswordLength + " characters long"
                        }
                ]
        },
};

var MessageRules =
{
        text: {
                identifier: "message-text",
                rules: [
                        {
                                type: "empty",
                                prompt: "Please enter a message"
                        }
                ]
        }
};

$(document).ready(function() {
        $("#form-sign-up")
		.form(SignInRules)
		.modal({duration: 300})
                .modal("attach events", ".sign-up");

        $("#form-sign-in")
		.form(SignInRules)
		.modal({duration: 300})
                .modal("attach events", ".sign-in");

        $("#form-message")
		.form(MessageRules)
		.modal({duration: 300})
                .modal("attach events", ".message");

        $(".dropdown").dropdown();

        // range input labels
        $(".range-wrapper").each(function(idx, el) {
                var input = $(el).children("input");
                var label = $(el).children("label");
                input.on("change input", function(e) {
                        label.text(e.target.value + " " + label.data("name"));
                });
        });

        // datepickers
        $(".datepicker").pickadate();

	// tabs
	$(".menu .item").tab();

	// progress
	$(".progress").progress();

	// parallax
	$(".parallax").parallax();

	// owl carousel
	$(".owl-carousel").owlCarousel({
		items: 3,
		autoHeight: true
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

