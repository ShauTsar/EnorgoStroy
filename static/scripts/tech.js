$(function() {
    $('input[name="deadline"]').daterangepicker({
        singleDatePicker: true,
        showDropdowns: true,
        autoUpdateInput: false,
        locale: {
            format: 'DD.MM.YYYY',
            cancelLabel: 'Очистить'
        }
    });

    $('input[name="employmentDate"]').daterangepicker({
        singleDatePicker: true,
        showDropdowns: true,
        autoUpdateInput: false,
        locale: {
            format: 'DD.MM.YYYY',
            cancelLabel: 'Очистить'
        }
    });

    $('input[name="deadline"]').on('apply.daterangepicker', function(ev, picker) {
        $(this).val(picker.startDate.format('DD.MM.YYYY'));
    });

    $('input[name="employmentDate"]').on('apply.daterangepicker', function(ev, picker) {
        $(this).val(picker.startDate.format('DD.MM.YYYY'));
    });

    $('input[name="deadline"]').on('cancel.daterangepicker', function(ev, picker) {
        $(this).val('');
    });

    $('input[name="employmentDate"]').on('cancel.daterangepicker', function(ev, picker) {
        $(this).val('');
    });

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

