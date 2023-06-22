$(function() {
    $('input[name="employment"]').change(function() {
        if (this.checked) {
            $("#employmentDate").show();
            $("#deadline").hide();
            $("#deadlineName").hide();
        } else {
            $("#employmentDate").hide();
            $("#deadline").show();
            $("#deadlineName").show();
        }
    });
});
function handleLocationChange() {
    var locationSelect = document.getElementById("location");
    var objectField = document.getElementById("objectField");

    if (locationSelect.value === "object") {
        objectField.style.display = "block";
    } else {
        objectField.style.display = "none";
    }
}

