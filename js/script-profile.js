$(document).ready(function() {
	$("#fullname").editable("/profile", {
		type: "textarea",
		cancel: "Cancel",
		submit: "Ok",
		indicator: "Saving..",
		submitdata: {what: "fullname"}
	});

	$("#description").editable("/profile", {
		type: "textarea",
		cancel: "Cancel",
		submit: "Ok",
		indicator: "Saving..",
		submitdata: {what: "description"}
	});

	$("#profession").editable("/profile", {
		data: "{'Art Director': 'Art Director', 'Creative Designer': 'Creative Designer', 'Creative Technologist': 'Creative Technologist', 'Copywriter': 'Copywriter', 'Planner': 'Planner', 'Manager': 'Manager', 'Programmer': 'Programmer', 'Producer': 'Producer', 'selected': 'Creative Technologist'}",
		type: "select",
		cancel: "Cancel",
		submit: "Ok",
		indicator: "Saving..",
		submitdata: {what: "profession"}
	});

	setImageInputPreview("#avatar-input", "#avatar-preview", "/profile");
});
