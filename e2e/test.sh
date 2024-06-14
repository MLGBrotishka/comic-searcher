#!/bin/bash

# Переменные окружения
SERVER_URL="http://localhost:8080"
LOGIN_USER="testuser"
LOGIN_PASSWORD="secure_debug"
SEARCH_KEYWORDS="apple,doctor"
EXPECTED_COMIC="apple a day"

if ! command -v jq &> /dev/null
then
    echo "jq could not be found, please install."
    exit 1
fi

# Запуск сервера
echo "Starting server..."
MAKE_PID=$!
make run &
BACK_PID=$!

echo "Ожидание 2 секунд..."
sleep 2
# Логин и получение токена
echo "Logging in..."
TOKEN=$(curl -s -X POST "$SERVER_URL/login" \
             -H "Content-Type: application/json" \
             -d "{\"login\": \"$LOGIN_USER\", \"password\": \"$LOGIN_PASSWORD\"}" \
             | jq -r '.access_token')

if [ -z "$TOKEN" ]; then
    echo "Failed to get token"
    exit 1
fi

# Обновление базы данных
echo "Updating DB..."
curl -s -X POST "$SERVER_URL/update" \
     -H "Authorization:$TOKEN"

# Поиск комиксов
echo "Searching for comics..."
SEARCH_RESULT=$(curl -s -X GET "$SERVER_URL/pics?search=$SEARCH_KEYWORDS" \
                     -H "Authorization:$TOKEN")
echo "$SEARCH_RESULT"
# Проверка результата
COMICS_FOUND=$(echo "$SEARCH_RESULT" | jq -r '.urls[]')

if [[ "$COMICS_FOUND" =~ "$EXPECTED_COMIC" ]]; then
    echo "Expected comic '$EXPECTED_COMIC' not found in search results:"
    echo "$SEARCH_RESULT"
    exit 1
else
    echo "Test passed: Expected comic '$EXPECTED_COMIC' found."
fi

# Остановка сервера
echo "Stopping server..."
echo "Пока остановка вручную"
kill -2 $BACK_PID
wait $BACK_PID
