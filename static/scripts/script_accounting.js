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
