#!/usr/bin/env ruby
# virtual_cat.rb
# encoding: UTF-8

require 'json'
require 'date'
require 'fileutils'

COLORS = {
  green: "\e[92m",
  red: "\e[91m",
  yellow: "\e[93m",
  blue: "\e[94m",
  magenta: "\e[95m",
  reset: "\e[0m"
}

def colorize(text, color)
  "#{COLORS[color]}#{text}#{COLORS[:reset]}"
end

CAT_ASCII = <<~'CAT'
 /\_/\
( o.o )
 > ^ <
CAT

CONFIG_FILE = File.join(Dir.home, '.virtual_cat.json')

def load_state
  return { 'name' => 'Барсик', 'hunger' => 50, 'happiness' => 70,
           'energy' => 80, 'health' => 90, 'age' => 0,
           'last_update' => Time.now.iso8601 } unless File.exist?(CONFIG_FILE)
  JSON.parse(File.read(CONFIG_FILE))
end

def save_state(state)
  File.write(CONFIG_FILE, JSON.pretty_generate(state))
end

def now_iso
  Time.now.iso8601
end

def clamp(val, min=0, max=100)
  [[val, min].max, max].min
end

def get_emotion(state)
  h = state['happiness']
  return '😺' if h >= 80
  return '😸' if h >= 50
  return '😼' if h >= 30
  '😿'
end

def get_mood_text(state)
  h = state['happiness']
  return 'Мурлычет и трётся об ноги' if h >= 80
  return 'Играет с клубком ниток' if h >= 50
  return 'Сидит на подоконнике и смотрит в окно' if h >= 30
  'Лежит в углу и грустит'
end

def show_status(state)
  puts colorize(CAT_ASCII, :magenta)
  puts colorize("  Имя: #{state['name']}", :blue)
  puts colorize("  Возраст: #{state['age']} лет", :blue)
  puts "  #{get_emotion(state)}  #{get_mood_text(state)}"
  puts colorize("  Голод: #{state['hunger']}/100", state['hunger'] > 70 ? :yellow : :green)
  puts colorize("  Счастье: #{state['happiness']}/100", state['happiness'] > 50 ? :green : :red)
  puts colorize("  Энергия: #{state['energy']}/100", state['energy'] > 50 ? :green : :yellow)
  puts colorize("  Здоровье: #{state['health']}/100", state['health'] < 30 ? :red : :green)
end

def random_event(msg=nil)
  events = ['Мяу!', 'Котик трётся о ноги.', 'Принёс игрушку.',
            'Пытается поймать муху.', 'Свернулся клубком.']
  text = msg || events.sample
  puts colorize("  ✨ #{text}", :yellow)
end

def auto_mode
  puts colorize("🤖  Автоматический режим включён.", :magenta)
  loop do
    state = load_state
    state['hunger'] = clamp(state['hunger'] + rand(-5..15))
    state['happiness'] = clamp(state['happiness'] + rand(-5..10))
    state['energy'] = clamp(state['energy'] + rand(-5..10))
    state['health'] = clamp(state['health'] + rand(-3..5))
    save_state(state)
    system('clear') || system('cls')
    show_status(state)
    puts colorize("\nНажмите Ctrl+C для выхода", :yellow)
    sleep(10)
  end
rescue Interrupt
  puts colorize("\n👋  Автоматический режим завершён.", :blue)
end

def main
  if ARGV.empty?
    puts colorize("Usage: virtual_cat <status|feed|play|sleep|heal|rename|auto> [name]", :yellow)
    exit 1
  end
  action = ARGV[0]
  state = load_state

  case action
  when 'status'
    show_status(state)
  when 'feed'
    state['hunger'] = clamp(state['hunger'] - 30)
    state['happiness'] = clamp(state['happiness'] + 10)
    state['health'] = clamp(state['health'] + 5)
    puts colorize("🐟  Ням-ням! Котик поел.", :green)
    random_event("Мурлычет от удовольствия")
    save_state(state)
  when 'play'
    if state['energy'] < 20
      puts colorize("😿  Котик слишком устал для игр.", :red)
      return
    end
    state['happiness'] = clamp(state['happiness'] + 25)
    state['energy'] = clamp(state['energy'] - 20)
    state['hunger'] = clamp(state['hunger'] + 10)
    puts colorize("🧶  Игра с клубком! Котик доволен.", :green)
    random_event("Прыгает за лазерной указкой")
    save_state(state)
  when 'sleep'
    state['energy'] = clamp(state['energy'] + 40)
    state['hunger'] = clamp(state['hunger'] + 10)
    puts colorize("😴  Котик уснул. Сладких снов!", :blue)
    random_event("Во сне дёргает лапками")
    save_state(state)
  when 'heal'
    if state['energy'] < 20
      puts colorize("😿  Нет сил лечить котика.", :red)
      return
    end
    state['health'] = clamp(state['health'] + 30)
    state['energy'] = clamp(state['energy'] - 20)
    puts colorize("💊  Котик вылечен! Он благодарен.", :green)
    save_state(state)
  when 'rename'
    if ARGV.size < 2
      puts colorize("Укажите имя: rename <имя>", :red)
      return
    end
    state['name'] = ARGV[1]
    puts colorize("🐱  Котика теперь зовут #{state['name']}!", :blue)
    save_state(state)
  when 'auto'
    auto_mode
  else
    puts colorize("Неизвестное действие.", :red)
  end
end

main if __FILE__ == $0
