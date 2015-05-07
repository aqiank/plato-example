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

function logout() {
        $.post("/" + APIVersion + "/logout", function(resp) {
                if (resp.status == "ok") {
                        window.location = "/";
                }
                console.log(resp);
        }, "json");
}

$(document).ready(function() {
        $("#form-sign-in").form(SignInRules);
        $("#form-sign-up").form(SignInRules);
        $("#form-sign-in")
                .modal("attach events", ".sign-in", "show")
                .modal({
                        onApprove: function() {
                                return false;
                        }
                });
        $("#form-sign-up")
                .modal("attach events", ".sign-up", "show")
                .modal({
                        onApprove: function() {
                                return false;
                        }
                });

        $(".dropdown").dropdown();

        // handle range input labels
        $(".range-wrapper").each(function(idx, el) {
                var input = $(el).children("input");
                var label = $(el).children("label");
                input.on("change input", function(e) {
                        label.text(e.target.value + " " + label.data("name"));
                });
        });
});
