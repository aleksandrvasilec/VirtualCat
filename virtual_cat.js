// virtual_cat.js
#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const os = require('os');
const readline = require('readline');

const COLORS = {
    green: '\x1b[92m',
    red: '\x1b[91m',
    yellow: '\x1b[93m',
    blue: '\x1b[94m',
    magenta: '\x1b[95m',
    reset: '\x1b[0m'
};

function colorize(text, color) {
    return COLORS[color] + text + COLORS.reset;
}

const CAT_ASCII = `
 /\\_/\\
( o.o )
 > ^ <
`;

const CONFIG_FILE = path.join(os.homedir(), '.virtual_cat.json');

function loadState() {
    try {
        const data = fs.readFileSync(CONFIG_FILE, 'utf8');
        return JSON.parse(data);
    } catch (err) {
        return {
            name: 'Барсик',
            hunger: 50,
            happiness: 70,
            energy: 80,
            health: 90,
            age: 0,
            lastUpdate: new Date().toISOString()
        };
    }
}

function saveState(state) {
    fs.writeFileSync(CONFIG_FILE, JSON.stringify(state, null, 2));
}

function nowISO() {
    return new Date().toISOString();
}

function clamp(val, min, max) {
    return Math.max(min, Math.min(val, max));
}

function getEmotion(state) {
    if (state.happiness >= 80) return '😺';
    if (state.happiness >= 50) return '😸';
    if (state.happiness >= 30) return '😼';
    return '😿';
}

function getMoodText(state) {
    if (state.happiness >= 80) return 'Мурлычет и трётся об ноги';
    if (state.happiness >= 50) return 'Играет с клубком ниток';
    if (state.happiness >= 30) return 'Сидит на подоконнике и смотрит в окно';
    return 'Лежит в углу и грустит';
}

function showStatus(state) {
    console.log(colorize(CAT_ASCII, 'magenta'));
    console.log(colorize(`  Имя: ${state.name}`, 'blue'));
    console.log(colorize(`  Возраст: ${state.age} лет`, 'blue'));
    console.log(`  ${getEmotion(state)}  ${getMoodText(state)}`);
    console.log(colorize(`  Голод: ${state.hunger}/100`, state.hunger > 70 ? 'yellow' : 'green'));
    console.log(colorize(`  Счастье: ${state.happiness}/100`, state.happiness > 50 ? 'green' : 'red'));
    console.log(colorize(`  Энергия: ${state.energy}/100`, state.energy > 50 ? 'green' : 'yellow'));
    console.log(colorize(`  Здоровье: ${state.health}/100`, state.health < 30 ? 'red' : 'green'));
}

function randomEvent(msg) {
    const events = ['Мяу!', 'Котик трётся о ноги.', 'Принёс игрушку.',
        'Пытается поймать муху.', 'Свернулся клубком.'];
    const text = msg || events[Math.floor(Math.random() * events.length)];
    console.log(colorize(`  ✨ ${text}`, 'yellow'));
}

function autoMode() {
    console.log(colorize('🤖  Автоматический режим включён.', 'magenta'));
    const interval = setInterval(() => {
        let state = loadState();
        state.hunger = clamp(state.hunger + Math.floor(Math.random() * 21) - 5, 0, 100);
        state.happiness = clamp(state.happiness + Math.floor(Math.random() * 16) - 5, 0, 100);
        state.energy = clamp(state.energy + Math.floor(Math.random() * 16) - 5, 0, 100);
        state.health = clamp(state.health + Math.floor(Math.random() * 9) - 3, 0, 100);
        saveState(state);
        console.clear();
        showStatus(state);
        console.log(colorize('\nНажмите Ctrl+C для выхода', 'yellow'));
    }, 10000);
    process.on('SIGINT', () => {
        clearInterval(interval);
        console.log(colorize('\n👋  Автоматический режим завершён.', 'blue'));
        process.exit(0);
    });
}

function main() {
    const args = process.argv.slice(2);
    if (args.length < 1) {
        console.log(colorize('Usage: node virtual_cat.js <status|feed|play|sleep|heal|rename|auto> [name]', 'yellow'));
        process.exit(1);
    }
    const action = args[0];
    const state = loadState();

    switch (action) {
        case 'status':
            showStatus(state);
            break;
        case 'feed':
            state.hunger = clamp(state.hunger - 30, 0, 100);
            state.happiness = clamp(state.happiness + 10, 0, 100);
            state.health = clamp(state.health + 5, 0, 100);
            console.log(colorize('🐟  Ням-ням! Котик поел.', 'green'));
            randomEvent('Мурлычет от удовольствия');
            saveState(state);
            break;
        case 'play':
            if (state.energy < 20) {
                console.log(colorize('😿  Котик слишком устал для игр.', 'red'));
                return;
            }
            state.happiness = clamp(state.happiness + 25, 0, 100);
            state.energy = clamp(state.energy - 20, 0, 100);
            state.hunger = clamp(state.hunger + 10, 0, 100);
            console.log(colorize('🧶  Игра с клубком! Котик доволен.', 'green'));
            randomEvent('Прыгает за лазерной указкой');
            saveState(state);
            break;
        case 'sleep':
            state.energy = clamp(state.energy + 40, 0, 100);
            state.hunger = clamp(state.hunger + 10, 0, 100);
            console.log(colorize('😴  Котик уснул. Сладких снов!', 'blue'));
            randomEvent('Во сне дёргает лапками');
            saveState(state);
            break;
        case 'heal':
            if (state.energy < 20) {
                console.log(colorize('😿  Нет сил лечить котика.', 'red'));
                return;
            }
            state.health = clamp(state.health + 30, 0, 100);
            state.energy = clamp(state.energy - 20, 0, 100);
            console.log(colorize('💊  Котик вылечен! Он благодарен.', 'green'));
            saveState(state);
            break;
        case 'rename':
            if (args.length < 2) {
                console.log(colorize('Укажите имя: rename <имя>', 'red'));
                return;
            }
            state.name = args[1];
            console.log(colorize(`🐱  Котика теперь зовут ${state.name}!`, 'blue'));
            saveState(state);
            break;
        case 'auto':
            autoMode();
            break;
        default:
            console.log(colorize('Неизвестное действие.', 'red'));
    }
}

main();
