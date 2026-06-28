// virtual_cat.cs
using System;
using System.IO;
using System.Text.Json;
using System.Threading;
using System.Runtime.InteropServices;

class VirtualCat
{
    static string Colorize(string text, string color)
    {
        string col = color switch
        {
            "green" => "\x1b[92m",
            "red" => "\x1b[91m",
            "yellow" => "\x1b[93m",
            "blue" => "\x1b[94m",
            "magenta" => "\x1b[95m",
            _ => "\x1b[0m"
        };
        return col + text + "\x1b[0m";
    }

    const string CAT_ASCII = @"
 /\_/\
( o.o )
 > ^ <
";

    class State
    {
        public string Name { get; set; }
        public int Hunger { get; set; }
        public int Happiness { get; set; }
        public int Energy { get; set; }
        public int Health { get; set; }
        public int Age { get; set; }
        public string LastUpdate { get; set; }
    }

    static string ConfigFile => Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".virtual_cat.json");

    static State LoadState()
    {
        if (!File.Exists(ConfigFile))
            return new State { Name = "Барсик", Hunger = 50, Happiness = 70, Energy = 80, Health = 90, Age = 0, LastUpdate = DateTime.Now.ToString("o") };
        string json = File.ReadAllText(ConfigFile);
        return JsonSerializer.Deserialize<State>(json) ?? new State();
    }

    static void SaveState(State state)
    {
        string json = JsonSerializer.Serialize(state, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(ConfigFile, json);
    }

    static string NowISO() => DateTime.Now.ToString("o");

    static int Clamp(int val, int min, int max) => Math.Max(min, Math.Min(val, max));

    static string GetEmotion(State state)
    {
        if (state.Happiness >= 80) return "😺";
        if (state.Happiness >= 50) return "😸";
        if (state.Happiness >= 30) return "😼";
        return "😿";
    }

    static string GetMoodText(State state)
    {
        if (state.Happiness >= 80) return "Мурлычет и трётся об ноги";
        if (state.Happiness >= 50) return "Играет с клубком ниток";
        if (state.Happiness >= 30) return "Сидит на подоконнике и смотрит в окно";
        return "Лежит в углу и грустит";
    }

    static void ShowStatus(State state)
    {
        Console.WriteLine(Colorize(CAT_ASCII, "magenta"));
        Console.WriteLine(Colorize($"  Имя: {state.Name}", "blue"));
        Console.WriteLine(Colorize($"  Возраст: {state.Age} лет", "blue"));
        Console.WriteLine($"  {GetEmotion(state)}  {GetMoodText(state)}");
        Console.WriteLine(Colorize($"  Голод: {state.Hunger}/100", state.Hunger > 70 ? "yellow" : "green"));
        Console.WriteLine(Colorize($"  Счастье: {state.Happiness}/100", state.Happiness > 50 ? "green" : "red"));
        Console.WriteLine(Colorize($"  Энергия: {state.Energy}/100", state.Energy > 50 ? "green" : "yellow"));
        Console.WriteLine(Colorize($"  Здоровье: {state.Health}/100", state.Health < 30 ? "red" : "green"));
    }

    static void RandomEvent(string msg = "")
    {
        string[] events = { "Мяу!", "Котик трётся о ноги.", "Принёс игрушку.",
                           "Пытается поймать муху.", "Свернулся клубком." };
        string text = string.IsNullOrEmpty(msg) ? events[new Random().Next(events.Length)] : msg;
        Console.WriteLine(Colorize($"  ✨ {text}", "yellow"));
    }

    static void AutoMode()
    {
        Console.WriteLine(Colorize("🤖  Автоматический режим включён.", "magenta"));
        var rand = new Random();
        while (true)
        {
            var state = LoadState();
            state.Hunger = Clamp(state.Hunger + rand.Next(-5, 16), 0, 100);
            state.Happiness = Clamp(state.Happiness + rand.Next(-5, 11), 0, 100);
            state.Energy = Clamp(state.Energy + rand.Next(-5, 11), 0, 100);
            state.Health = Clamp(state.Health + rand.Next(-3, 6), 0, 100);
            SaveState(state);
            Console.Clear();
            ShowStatus(state);
            Console.WriteLine(Colorize("\nНажмите Ctrl+C для выхода", "yellow"));
            Thread.Sleep(10000);
        }
    }

    static void Main(string[] args)
    {
        if (args.Length < 1)
        {
            Console.WriteLine(Colorize("Usage: virtual_cat <status|feed|play|sleep|heal|rename|auto> [name]", "yellow"));
            return;
        }
        string action = args[0];
        var state = LoadState();

        switch (action)
        {
            case "status":
                ShowStatus(state);
                break;
            case "feed":
                state.Hunger = Clamp(state.Hunger - 30, 0, 100);
                state.Happiness = Clamp(state.Happiness + 10, 0, 100);
                state.Health = Clamp(state.Health + 5, 0, 100);
                Console.WriteLine(Colorize("🐟  Ням-ням! Котик поел.", "green"));
                RandomEvent("Мурлычет от удовольствия");
                SaveState(state);
                break;
            case "play":
                if (state.Energy < 20)
                {
                    Console.WriteLine(Colorize("😿  Котик слишком устал для игр.", "red"));
                    return;
                }
                state.Happiness = Clamp(state.Happiness + 25, 0, 100);
                state.Energy = Clamp(state.Energy - 20, 0, 100);
                state.Hunger = Clamp(state.Hunger + 10, 0, 100);
                Console.WriteLine(Colorize("🧶  Игра с клубком! Котик доволен.", "green"));
                RandomEvent("Прыгает за лазерной указкой");
                SaveState(state);
                break;
            case "sleep":
                state.Energy = Clamp(state.Energy + 40, 0, 100);
                state.Hunger = Clamp(state.Hunger + 10, 0, 100);
                Console.WriteLine(Colorize("😴  Котик уснул. Сладких снов!", "blue"));
                RandomEvent("Во сне дёргает лапками");
                SaveState(state);
                break;
            case "heal":
                if (state.Energy < 20)
                {
                    Console.WriteLine(Colorize("😿  Нет сил лечить котика.", "red"));
                    return;
                }
                state.Health = Clamp(state.Health + 30, 0, 100);
                state.Energy = Clamp(state.Energy - 20, 0, 100);
                Console.WriteLine(Colorize("💊  Котик вылечен! Он благодарен.", "green"));
                SaveState(state);
                break;
            case "rename":
                if (args.Length < 2)
                {
                    Console.WriteLine(Colorize("Укажите имя: rename <имя>", "red"));
                    return;
                }
                state.Name = args[1];
                Console.WriteLine(Colorize($"🐱  Котика теперь зовут {state.Name}!", "blue"));
                SaveState(state);
                break;
            case "auto":
                AutoMode();
                break;
            default:
                Console.WriteLine(Colorize("Неизвестное действие.", "red"));
                break;
        }
    }
}
