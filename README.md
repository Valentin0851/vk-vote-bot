Бот для создания голосований в mattermost и управлении ими.
Создан в качестве ответа на тестовое задание на стажировку в компанию VK.
Автор - Оленев Артём
Установка:

Для уствновки склонируйте данный репозиторий:
git clone https://github.com/Valentin0851/vk-vote-bot/tree/master
Запуск:

    Замените в файле config.yaml константы на свои соответствующие значения
    Запустите tarantool и mattermost
    Соберите контейнер с ботом командой:
    docker compose build mattermost-bot

    Запустите бота командой:
    docker compose up -d mattermost-bot

Использование:

Используйте в своем канале mattermost следующие команды:

    Создание голосования:
    /vote create "<Question>" "<Option 1>" "<Option 2>" ... "<Option n>"
    После создания голосования в чат выведется ссобщение с id голосования, которое будет необходимо для последующей работы с ним.
    Голосование:
    /poll vote <poll_id> <option>
    получение результатов:
    /poll results <poll_id>
    Завершение голосования:
    /poll end <poll_id>
    Удаление голосования:
    /poll delete <poll_id>
