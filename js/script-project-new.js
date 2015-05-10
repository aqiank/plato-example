var NewProjectRules =
{
        title: {
                identifier: "title",
                rules: [
                        {
                                type: "empty",
                                prompt: "Please enter project title"
                        }
                ]
        },
        content: {
                identifier: "content",
                rules: [
                        {
                                type: "empty",
                                prompt: "Please enter project description"
                        }
                ]
        },
	startDate: {
                identifier: "start-date",
                rules: [
                        {
                                type: "empty",
                                prompt: "Please enter start date"
                        }
                ]
	},
	endDate: {
                identifier: "end-date",
                rules: [
                        {
                                type: "empty",
                                prompt: "Please enter end date"
                        }
                ]
	}
};

$(document).ready(function() {
	setImageInputPreview("#image-input", "#image-preview");

	$("#form-project-new")
		.form(NewProjectRules);
});
