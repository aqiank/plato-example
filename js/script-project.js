$(document).ready(function() {
	$(".task-edit-trigger").click(function(e) {
		$("#task-edit-title").val($(this).data("title"));
		$("#task-edit-description").val($(this).data("description"));
		$("#task-edit-start-date").val($(this).data("startdate"));
		$("#task-edit-end-date").val($(this).data("enddate"));
		$("#task-edit-is-milestone").prop("checked", $(this).data("ismilestone"));
		$("#task-edit-restricted").prop("checked", $(this).data("restricted"));
		$("#task-edit-done").prop("checked", $(this).data("done"));
		$("#task-edit-id").val($(this).data("id"));
	});
});
