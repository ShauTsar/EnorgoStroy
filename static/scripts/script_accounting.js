const cards = document.querySelectorAll('.card');
cards.forEach(card => {
    card.addEventListener('mouseover', () => {
        card.classList.add('pulse-animation');
    });
    card.addEventListener('mouseout', () => {
        card.classList.remove('pulse-animation');
    });
});
$(document).ready(function() {
    $("#add-item-button").click(function() {
        $("#modal-form").modal("show");
    });
});
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