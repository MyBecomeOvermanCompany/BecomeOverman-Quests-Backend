from pprint import pprint
from dotenv import load_dotenv
import requests
import os
import json
import sys

# Загружаем переменные из .env файла
load_dotenv()

api_key = os.getenv("API_KEY")
ai_model = "moonshotai/Kimi-K2-Thinking"


def ask_neural_network(user_message: str) -> str:
    """
    Отправляет запрос нейросети и возвращает ответ в формате JSON
    """
    url = "https://api.intelligence.io.solutions/api/v1/chat/completions"

    headers = {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + api_key
    }

    # Строгий промпт с JSON шаблоном
    system_prompt = """
    Ты помощник для генерации квестов в формате строгого JSON. 
    ВОЗВРАЩАЙ ТОЛЬКО JSON БЕЗ ЛЮБЫХ ДОПОЛНИТЕЛЬНЫХ ТЕКСТОВ И КОММЕНТАРИЕВ!

    Структура JSON должна быть такой:
    {
        "quest": {
            "title": "Название квеста",
            "description": "Описание квеста [Generated]",
            "category": "health/willpower/intelligence/creativity/social",
            "rarity": "common/rare/epic/legendary",
            "difficulty": 1-5,
            "price": 10-100,
            "tasks_count": 3-7,
            "reward_xp": 50-500,
            "reward_coin": 25-250,
            "time_limit_hours": 24-336
        },
        "tasks": [
            {
                "title": "Название задачи 1",
                "description": "Описание задачи 1",
                "difficulty": 1-3,
                "rarity": "common/rare/epic",
                "category": "health/willpower/intelligence/creativity/social",
                "base_xp_reward": 10-50,
                "base_coin_reward": 5-25,
                "task_order": 1
            }
        ]
    }

    Правила:
    - difficulty квеста должен быть средним от difficulty задач
    - price = reward_coin * 1.5 (округлить)
    - tasks_count должно соответствовать количеству задач в массиве
    - time_limit_hours: 24-168 (1-7 дней)
    - reward_xp = сумма base_xp_reward всех задач * 1.5
    - reward_coin = сумма base_coin_reward всех задач * 1.5
    """

    data = {
        "model": ai_model,
        "messages": [
            {
                "role": "system",
                "content": system_prompt,
            },
            {
                "role": "user",
                "content": user_message
            }
        ],
        "temperature": 0.7
    }

    try:
        response = requests.post(url, headers=headers, json=data)
        response.raise_for_status()
        data = response.json()

        text = data['choices'][0]['message']['content']

        # Очищаем ответ от возможных тегов thinking
        if '</think>\n\n' in text:
            prepared_text = text.split('</think>\n\n')[1]
        else:
            prepared_text = text

        # Убедимся что это валидный JSON
        json.loads(prepared_text)
        return prepared_text

    except Exception as e:
        return json.dumps({"error": f"AI service error: {str(e)}"})


# Пример использования
if __name__ == "__main__":
    if len(sys.argv) > 1:
        user_message = sys.argv[1]
        result = ask_neural_network(user_message)
        print(result)
    else:
        print('{"error": "No message provided"}')