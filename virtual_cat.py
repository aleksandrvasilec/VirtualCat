# virtual_cat.py
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import os
import json
import time
import random
import argparse
from datetime import datetime, timedelta
from pathlib import Path

# ANSI colors
COLORS = {
    'green': '\033[92m',
    'red': '\033[91m',
    'yellow': '\033[93m',
    'blue': '\033[94m',
    'magenta': '\033[95m',
    'reset': '\033[0m'
}

def colorize(text, color):
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

CONFIG_FILE = Path.home() / '.virtual_cat.json'

# ASCII-арт котика
CAT_ASCII = r"""
 /\_/\
( o.o )
 > ^ <
"""

def get_default_state():
    return {
        'name': 'Барсик',
        'hunger': 50,
        'happiness': 70,
        'energy': 80,
        'health': 90,
        'age': 0,
        'last_update': datetime.now().isoformat()
    }

def load_state():
    if CONFIG_FILE.exists():
        with open(CONFIG_FILE, 'r') as f:
            return json.load(f)
    return get_default_state()

def save_state(state):
    with open(CONFIG_FILE, 'w') as f:
        json.dump(state, f, indent=2)

def update_age(state):
    last = datetime.fromisoformat(state['last_update'])
    now = datetime.now()
    delta = (now - last).total_seconds()
    # 1 день = 1 год кота (упрощённо)
    state['age'] += delta // 86400
    state['last_update'] = now.isoformat()
    return state

def clamp(value, min_val=0, max_val=100):
    return max(min_val, min(value, max_val))

def get_emotion(state):
    h = state['happiness']
    if h >= 80:
        return '😺'  # счастлив
    elif h >= 50:
        return '😸'  # игривый
    elif h >= 30:
        return '😼'  # задумчивый
    else:
        return '😿'  # грустный

def get_mood_text(state):
    h = state['happiness']
    if h >= 80:
        return "Мурлычет и трётся об ноги"
    elif h >= 50:
        return "Играет с клубком ниток"
    elif h >= 30:
        return "Сидит на подоконнике и смотрит в окно"
    else:
        return "Лежит в углу и грустит"

def show_status(state):
    state = update_age(state)
    save_state(state)
    name = state['name']
    hunger = state['hunger']
    happiness = state['happiness']
    energy = state['energy']
    health = state['health']
    age = state['age']

    print(colorize(CAT_ASCII, 'magenta'))
    print(colorize(f"  Имя: {name}", 'blue'))
    print(colorize(f"  Возраст: {age} лет", 'blue'))
    print(f"  {get_emotion(state)}  {get_mood_text(state)}")
    print(colorize(f"  Голод: {hunger}/100", 'yellow' if hunger > 70 else 'green'))
    print(colorize(f"  Счастье: {happiness}/100", 'green' if happiness > 50 else 'red'))
    print(colorize(f"  Энергия: {energy}/100", 'green' if energy > 50 else 'yellow'))
    print(colorize(f"  Здоровье: {health}/100", 'red' if health < 30 else 'green'))

def feed(state):
    state['hunger'] = clamp(state['hunger'] - 30)
    state['happiness'] = clamp(state['happiness'] + 10)
    state['health'] = clamp(state['health'] + 5)
    print(colorize("🐟  Ням-ням! Котик поел.", 'green'))
    show_random_event("Мурлычет от удовольствия")

def play(state):
    if state['energy'] < 20:
        print(colorize("😿  Котик слишком устал для игр.", 'red'))
        return
    state['happiness'] = clamp(state['happiness'] + 25)
    state['energy'] = clamp(state['energy'] - 20)
    state['hunger'] = clamp(state['hunger'] + 10)
    print(colorize("🧶  Игра с клубком! Котик доволен.", 'green'))
    show_random_event("Прыгает за лазерной указкой")

def sleep_cat(state):
    state['energy'] = clamp(state['energy'] + 40)
    state['hunger'] = clamp(state['hunger'] + 10)
    print(colorize("😴  Котик уснул. Сладких снов!", 'blue'))
    show_random_event("Во сне дёргает лапками")

def heal(state):
    if state['energy'] < 20:
        print(colorize("😿  Нет сил лечить котика.", 'red'))
        return
    state['health'] = clamp(state['health'] + 30)
    state['energy'] = clamp(state['energy'] - 20)
    print(colorize("💊  Котик вылечен! Он благодарен.", 'green'))

def rename(state, name):
    state['name'] = name
    print(colorize(f"🐱  Котика теперь зовут {name}!", 'blue'))

def show_random_event(message):
    events = [
        "Мяу!",
        "Котик трётся о ноги.",
        "Принёс игрушку.",
        "Пытается поймать муху.",
        "Свернулся клубком."
    ]
    print(colorize(f"  ✨ {message or random.choice(events)}", 'yellow'))

def auto_mode():
    state = load_state()
    print(colorize("🤖  Автоматический режим включён.", 'magenta'))
    try:
        while True:
            state = update_age(state)
            # Случайное изменение параметров
            state['hunger'] = clamp(state['hunger'] + random.randint(-5, 15))
            state['happiness'] = clamp(state['happiness'] + random.randint(-5, 10))
            state['energy'] = clamp(state['energy'] + random.randint(-5, 10))
            state['health'] = clamp(state['health'] + random.randint(-3, 5))
            save_state(state)
            os.system('clear' if os.name == 'posix' else 'cls')
            show_status(state)
            print(colorize("\nНажмите Ctrl+C для выхода", 'yellow'))
            time.sleep(10)
    except KeyboardInterrupt:
        print(colorize("\n👋  Автоматический режим завершён.", 'blue'))

def main():
    parser = argparse.ArgumentParser(description="Virtual Cat – виртуальный питомец")
    parser.add_argument('action', nargs='?', default='status',
                        choices=['status', 'feed', 'play', 'sleep', 'heal', 'rename', 'auto'],
                        help='Действие')
    parser.add_argument('name', nargs='?', help='Новое имя котика (для rename)')
    args = parser.parse_args()

    state = load_state()
    state = update_age(state)

    if args.action == 'status':
        show_status(state)
    elif args.action == 'feed':
        feed(state)
        save_state(state)
    elif args.action == 'play':
        play(state)
        save_state(state)
    elif args.action == 'sleep':
        sleep_cat(state)
        save_state(state)
    elif args.action == 'heal':
        heal(state)
        save_state(state)
    elif args.action == 'rename':
        if args.name:
            rename(state, args.name)
            save_state(state)
        else:
            print(colorize("Укажите имя: rename <имя>", 'red'))
    elif args.action == 'auto':
        auto_mode()
    else:
        print(colorize("Неизвестное действие.", 'red'))

if __name__ == '__main__':
    main()
