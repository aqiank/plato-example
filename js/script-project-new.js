$(document).ready(function() {
	setImageInputPreview("#image-input", "#image-preview");

	$("#form-project-new").submit(function(e) {
		var form = $(this);

		$.ajax({
			url: form.attr("action"),
			method: "POST",
			dataType: "json",
			data: form.serialize(),
		}).done(function(resp) {
			Materialize.toast("Successfully created project!", 1000, "green");
			window.location = "/project/" + resp.data;
		}).fail(function(resp) {
			Materialize.toast(resp.responseText, 1000, "red");
		});

		e.preventDefault();
	});
});
