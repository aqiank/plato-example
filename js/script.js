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

$(document).ready(function() {
        $("#form-sign-up")
		.form(SignInRules)
		.modal({duration: 300})
                .modal("attach events", ".sign-up");

        $("#form-sign-in")
		.form(SignInRules)
		.modal({duration: 300})
                .modal("attach events", ".sign-in");

        $(".dropdown").dropdown();

        // handle range input labels
        $(".range-wrapper").each(function(idx, el) {
                var input = $(el).children("input");
                var label = $(el).children("label");
                input.on("change input", function(e) {
                        label.text(e.target.value + " " + label.data("name"));
                });
        });

        // handle datepickers
        $(".datepicker").pickadate();

	// handle tabs
	$(".menu .item").tab();
});

function setImageInputPreview(input, preview, uploadURL, success) {
	$(input).change(function(){
		if (this.files && this.files[0]) {
			var reader = new FileReader();
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

