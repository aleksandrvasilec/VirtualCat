// virtual_cat.java
import java.io.*;
import java.nio.file.*;
import java.time.*;
import java.time.format.*;
import java.util.*;
import com.google.gson.*;

public class virtual_cat {
    private static final String RESET = "\u001B[0m";
    private static final String GREEN = "\u001B[92m";
    private static final String RED = "\u001B[91m";
    private static final String YELLOW = "\u001B[93m";
    private static final String BLUE = "\u001B[94m";
    private static final String MAGENTA = "\u001B[95m";

    private static String colorize(String text, String color) {
        return color + text + RESET;
    }

    private static final String CAT_ASCII = 
        "\n /\\_/\\\n( o.o )\n > ^ <\n";

    private static class State {
        String name;
        int hunger;
        int happiness;
        int energy;
        int health;
        int age;
        String lastUpdate;
    }

    private static String configFile = System.getProperty("user.home") + "/.virtual_cat.json";

    private static State loadState() throws IOException {
        Path path = Paths.get(configFile);
        if (!Files.exists(path)) {
            State s = new State();
            s.name = "Барсик";
            s.hunger = 50;
            s.happiness = 70;
            s.energy = 80;
            s.health = 90;
            s.age = 0;
            s.lastUpdate = Instant.now().toString();
            return s;
        }
        String json = new String(Files.readAllBytes(path));
        Gson gson = new Gson();
        return gson.fromJson(json, State.class);
    }

    private static void saveState(State state) throws IOException {
        Gson gson = new GsonBuilder().setPrettyPrinting().create();
        String json = gson.toJson(state);
        Files.write(Paths.get(configFile), json.getBytes());
    }

    private static int clamp(int val, int min, int max) {
        return Math.max(min, Math.min(val, max));
    }

    private static String getEmotion(State state) {
        if (state.happiness >= 80) return "😺";
        if (state.happiness >= 50) return "😸";
        if (state.happiness >= 30) return "😼";
        return "😿";
    }

    private static String getMoodText(State state) {
        if (state.happiness >= 80) return "Мурлычет и трётся об ноги";
        if (state.happiness >= 50) return "Играет с клубком ниток";
        if (state.happiness >= 30) return "Сидит на подоконнике и смотрит в окно";
        return "Лежит в углу и грустит";
    }

    private static void showStatus(State state) {
        System.out.println(colorize(CAT_ASCII, MAGENTA));
        System.out.println(colorize("  Имя: " + state.name, BLUE));
        System.out.println(colorize("  Возраст: " + state.age + " лет", BLUE));
        System.out.println("  " + getEmotion(state) + "  " + getMoodText(state));
        System.out.println(colorize("  Голод: " + state.hunger + "/100",
                                    state.hunger > 70 ? YELLOW : GREEN));
        System.out.println(colorize("  Счастье: " + state.happiness + "/100",
                                    state.happiness > 50 ? GREEN : RED));
        System.out.println(colorize("  Энергия: " + state.energy + "/100",
                                    state.energy > 50 ? GREEN : YELLOW));
        System.out.println(colorize("  Здоровье: " + state.health + "/100",
                                    state.health < 30 ? RED : GREEN));
    }

    private static void randomEvent(String msg) {
        String[] events = {"Мяу!", "Котик трётся о ноги.", "Принёс игрушку.",
                           "Пытается поймать муху.", "Свернулся клубком."};
        String text = (msg == null || msg.isEmpty()) ? events[new Random().nextInt(events.length)] : msg;
        System.out.println(colorize("  ✨ " + text, YELLOW));
    }

    private static void autoMode() throws IOException, InterruptedException {
        System.out.println(colorize("🤖  Автоматический режим включён.", MAGENTA));
        Random rand = new Random();
        while (true) {
            State state = loadState();
            state.hunger = clamp(state.hunger + rand.nextInt(21) - 5, 0, 100);
            state.happiness = clamp(state.happiness + rand.nextInt(16) - 5, 0, 100);
            state.energy = clamp(state.energy + rand.nextInt(16) - 5, 0, 100);
            state.health = clamp(state.health + rand.nextInt(9) - 3, 0, 100);
            saveState(state);
            System.out.print("\033[H\033[2J");
            System.out.flush();
            showStatus(state);
            System.out.println(colorize("\nНажмите Ctrl+C для выхода", YELLOW));
            Thread.sleep(10000);
        }
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        if (args.length < 1) {
            System.out.println(colorize("Usage: virtual_cat <status|feed|play|sleep|heal|rename|auto> [name]", YELLOW));
            return;
        }
        String action = args[0];
        State state = loadState();

        switch (action) {
            case "status":
                showStatus(state);
                break;
            case "feed":
                state.hunger = clamp(state.hunger - 30, 0, 100);
                state.happiness = clamp(state.happiness + 10, 0, 100);
                state.health = clamp(state.health + 5, 0, 100);
                System.out.println(colorize("🐟  Ням-ням! Котик поел.", GREEN));
                randomEvent("Мурлычет от удовольствия");
                saveState(state);
                break;
            case "play":
                if (state.energy < 20) {
                    System.out.println(colorize("😿  Котик слишком устал для игр.", RED));
                    return;
                }
                state.happiness = clamp(state.happiness + 25, 0, 100);
                state.energy = clamp(state.energy - 20, 0, 100);
                state.hunger = clamp(state.hunger + 10, 0, 100);
                System.out.println(colorize("🧶  Игра с клубком! Котик доволен.", GREEN));
                randomEvent("Прыгает за лазерной указкой");
                saveState(state);
                break;
            case "sleep":
                state.energy = clamp(state.energy + 40, 0, 100);
                state.hunger = clamp(state.hunger + 10, 0, 100);
                System.out.println(colorize("😴  Котик уснул. Сладких снов!", BLUE));
                randomEvent("Во сне дёргает лапками");
                saveState(state);
                break;
            case "heal":
                if (state.energy < 20) {
                    System.out.println(colorize("😿  Нет сил лечить котика.", RED));
                    return;
                }
                state.health = clamp(state.health + 30, 0, 100);
                state.energy = clamp(state.energy - 20, 0, 100);
                System.out.println(colorize("💊  Котик вылечен! Он благодарен.", GREEN));
                saveState(state);
                break;
            case "rename":
                if (args.length < 2) {
                    System.out.println(colorize("Укажите имя: rename <имя>", RED));
                    return;
                }
                state.name = args[1];
                System.out.println(colorize("🐱  Котика теперь зовут " + state.name + "!", BLUE));
                saveState(state);
                break;
            case "auto":
                autoMode();
                break;
            default:
                System.out.println(colorize("Неизвестное действие.", RED));
        }
    }
}
