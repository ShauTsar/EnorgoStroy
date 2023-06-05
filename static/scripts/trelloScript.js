function moveRight(taskId) {
    var task = document.getElementById(taskId);
    var currentColumn = task.parentNode;
    var nextColumn = currentColumn.nextElementSibling;

    if (nextColumn && nextColumn.classList.contains("column-body")) {
        nextColumn.appendChild(task);
    }
}

function deleteTask(taskId) {
    var task = document.getElementById(taskId);
    task.parentNode.removeChild(task);
}
// Объект, представляющий задачу
class Task {
    constructor(id, title, assignee, description, comment) {
        this.id = id;
        this.title = title;
        this.assignee = assignee;
        this.description = description;
        this.comment = comment;
    }
}

// Глобальный объект для хранения досок и задач
const boards = {
    'bit-stroitelstvo': {
        name: 'Бит.Строительство',
        columns: ['Обсуждение', 'В процессе написания', 'На Оценку', 'Оценка ТЗ', 'Согласование', 'Задачи в работе', 'Тестирование', 'Завершено'],
        tasks: {}
    },
    'zup-reg': {
        name: 'ЗУП (рег)',
        columns: ['Обсуждение', 'В процессе написания', 'На Оценку', 'Оценка ТЗ', 'Согласование', 'Задачи в работе', 'Тестирование', 'Завершено'],
        tasks: {}
    },
    'ohrana-truda': {
        name: 'Охрана труда',
        columns: ['Обсуждение', 'В процессе написания', 'На Оценку', 'Оценка ТЗ', 'Согласование', 'Задачи в работе', 'Тестирование', 'Завершено'],
        tasks: {}
    },
    'rabota': {
        name: 'Работа',
        columns: ['Обсуждение', 'В процессе написания', 'На Оценку', 'Оценка ТЗ', 'Согласование', 'Задачи в работе', 'Тестирование', 'Завершено'],
        tasks: {}
    }
};

function openBoard(boardName) {
    // Скрытие всех колонок
    const columnElements = document.getElementsByClassName('column');
    for (let i = 0; i < columnElements.length; i++) {
        columnElements[i].style.display = 'none';
    }

    // Отображение выбранной колонки
    const selectedColumn = document.getElementById(`column-${boardName}`);
    selectedColumn.style.display = 'block';
}


// Создание новой задачи
function createTask(columnId) {
    const modalTitle = document.getElementById('taskModalLabel');
    modalTitle.textContent = 'Добавить задачу';

    const taskTitle = document.getElementById('taskTitle');
    taskTitle.value = '';

    const taskAssignee = document.getElementById('taskAssignee');
    taskAssignee.value = '';

    const taskDescription = document.getElementById('taskDescription');
    taskDescription.value = '';

    const taskComment = document.getElementById('taskComment');
    taskComment.value = '';

    const saveButton = document.getElementsByClassName('btn-primary')[0];
    saveButton.onclick = function() {
        const id = Date.now().toString(); // Генерация уникального идентификатора задачи
        const title = taskTitle.value;
        const assignee = taskAssignee.value;
        const description = taskDescription.value;
        const comment = taskComment.value;

        const task = new Task(id, title, assignee, description, comment);
        boards[columnId].tasks[id] = task;

        renderTasks(columnId);
        $('#taskModal').modal('hide');
    };

    $('#taskModal').modal('show');
}

// Редактирование задачи
function editTask(taskId, columnId) {
    const modalTitle = document.getElementById('taskModalLabel');
    modalTitle.textContent = 'Редактировать задачу';

    const task = boards[columnId].tasks[taskId];

    const taskTitle = document.getElementById('taskTitle');
    taskTitle.value = task.title;

    const taskAssignee = document.getElementById('taskAssignee');
    taskAssignee.value = task.assignee;

    const taskDescription = document.getElementById('taskDescription');
    taskDescription.value = task.description;

    const taskComment = document.getElementById('taskComment');
    taskComment.value = task.comment;

    const saveButton = document.getElementsByClassName('btn-primary')[0];
    saveButton.onclick = function() {
        task.title = taskTitle.value;
        task.assignee = taskAssignee.value;
        task.description = taskDescription.value;
        task.comment = taskComment.value;

        renderTasks(columnId);
        $('#taskModal').modal('hide');
    };

    $('#taskModal').modal('show');
}

// Удаление задачи
function deleteTask(taskId, columnId) {
    delete boards[columnId].tasks[taskId];
    renderTasks(columnId);
}

// Отображение задач в колонке
function renderTasks(columnId) {
    const columnBody = document.getElementById(`column-body-${columnId}`);
    columnBody.innerHTML = '';

    const tasks = boards[columnId].tasks;

    for (const taskId in tasks) {
        const task = tasks[taskId];
        const taskElement = createTaskElement(taskId, task);
        columnBody.appendChild(taskElement);
    }
}

// Создание DOM-элемента для задачи
function createTaskElement(taskId, task) {
    let taskElement = document.getElementById(taskId);

    if (!taskElement) {
        taskElement = document.createElement('div');
        taskElement.classList.add('task');
        taskElement.id = taskId;
    }

    taskElement.innerHTML = `
    <div class="task-header">
      <span class="task-title">${task.title}</span>
      <div class="task-actions">
        <button type="button" class="btn btn-primary btn-sm" onclick="editTask('${taskId}', '${task.column}')">Редактировать</button>
        <button type="button" class="btn btn-danger btn-sm" onclick="deleteTask('${taskId}', '${task.column}')">Удалить</button>
      </div>
    </div>
    <div class="task-details">
      <p><strong>Исполнитель:</strong> ${task.assignee}</p>
      <p><strong>Описание:</strong> ${task.description}</p>
      <p><strong>Комментарий:</strong> ${task.comment}</p>
    </div>
  `;

    return taskElement;
}

// Инициализация досок и задач
function initializeBoards() {
    for (const boardName in boards) {
        const board = boards[boardName];
        const boardColumnContainer = document.getElementById(`${boardName}-columns`);

        for (const columnName of board.columns) {
            const columnElement = document.createElement('div');
            const columnId = boardName + '-' + columnName.toLowerCase().replace(/\s/g, '-');
            columnElement.classList.add('column');
            columnElement.id = columnId;
            columnElement.innerHTML = `
        <div class="column-header">${columnName}</div>
        <div class="column-body" id="column-body-${columnId}">
        </div>
      `;

            boardColumnContainer.appendChild(columnElement);
        }
    }
}

// Инициализация приложения
function initializeApp() {
    initializeBoards();
    renderTasks('bit-stroitelstvo-obsuzhdenie');
    renderTasks('zup-reg-obsuzhdenie');
    renderTasks('ohrana-truda-obsuzhdenie');
    renderTasks('rabota-obsuzhdenie');
}

// Вызов функции инициализации приложения после загрузки страницы
window.onload = initializeApp;


