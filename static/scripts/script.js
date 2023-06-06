document.addEventListener("DOMContentLoaded", function() {
    // Проверка, авторизован ли пользователь
    let authenticated = "{{ .Authenticated }}";
    if (authenticated === "true") {
        // Показать выбор категории
        document.getElementById("login-form").style.display = "none";
        document.getElementById("request-list").style.display = "block";

        // Автоматическое отправление формы
        document.getElementById("category-form").submit();
    }
});

$(document).ready(function() {
    // Обработчик события изменения значения поля "Объект"
    $('#object').change(function() {
        // Получаем выбранное значение объекта
        let selectedObject = $(this).val();

        // Вызываем функцию для обновления списка подразделений
        updateDepartments(selectedObject);
    });

    // Функция для обновления списка подразделений
    function updateDepartments(selectedObject) {
        // Определяем список подразделений для каждого объекта
        let departments = {
            'Норильск': ['НМЗ 010', 'НМЗ ВЗИС', 'НМЗ Норникель  ПНР', 'НМЗ Норникель Вахтовый поселок', 'НМЗ Норникель ГПП-83', 'НМЗ Норникель Косвенный', 'НМЗ НСК'],
            'Мурманск': ['Мурманск BC1', 'Мурманск GWP5A', 'Мурманск Косвенный'],
            'ГПЗ': ['ГПЗ Амурский АВР Метрология', 'ГПЗ Амурский ВЗИС', 'ГПЗ Амурский Метрология Ф3', 'ГПЗ Амурский ПНР КИП Ф2', 'ГПЗ Амурский ПНР КИП Ф3', 'ГПЗ Амурский ПНР КИП Ф4',
                'ГПЗ Амурский ПНР-ПС Ф2', 'ГПЗ Амурский ПНР-ПС Ф3', 'ГПЗ Амурский ПНР-ПС Ф4', 'ГПЗ Амурский Техперевооружение', 'ГПЗ Амурский Фаза 2', 'ГПЗ Амурский Фаза 3', 'ГПЗ Амурский Фаза 4', 'ГПЗ Амурский Фаза 5'
                , 'ГПЗ Косвенный'],
            'Иркутск': ['ИНК Иркутский завод полимеров']
        };
        // Очищаем текущий список подразделений
        $('#department').empty();

        // Получаем список подразделений для выбранного объекта
        let selectedDepartments = departments[selectedObject];

        // Добавляем опции в список подразделений
        for (let i = 0; i < selectedDepartments.length; i++) {
            let department = selectedDepartments[i];
            $('#department').append('<option value="' + department + '">' + department + '</option>');
        }
    }
});
