// virtual_cat.cpp
#include <iostream>
#include <fstream>
#include <string>
#include <cstdlib>
#include <ctime>
#include <json/json.h> // sudo apt-get install libjsoncpp-dev

using namespace std;

const string RESET = "\033[0m";
const string GREEN = "\033[92m";
const string RED = "\033[91m";
const string YELLOW = "\033[93m";
const string BLUE = "\033[94m";
const string MAGENTA = "\033[95m";

string colorize(const string& text, const string& color) {
    return color + text + RESET;
}

const string CAT_ASCII = R"(
 /\_/\
( o.o )
 > ^ <
)";

string getHomeDir() {
    const char* home = getenv("HOME");
    if (!home) home = getenv("USERPROFILE");
    return string(home);
}

string getConfigFile() {
    return getHomeDir() + "/.virtual_cat.json";
}

Json::Value loadState() {
    ifstream f(getConfigFile());
    Json::Value root;
    if (!f) {
        root["name"] = "Барсик";
        root["hunger"] = 50;
        root["happiness"] = 70;
        root["energy"] = 80;
        root["health"] = 90;
        root["age"] = 0;
        root["last_update"] = "";
        return root;
    }
    f >> root;
    return root;
}

void saveState(const Json::Value& state) {
    ofstream f(getConfigFile());
    f << state.toStyledString();
}

string nowISO() {
    time_t t = time(nullptr);
    char buf[64];
    strftime(buf, sizeof(buf), "%Y-%m-%dT%H:%M:%S", localtime(&t));
    return string(buf);
}

int clamp(int value, int min=0, int max=100) {
    return max(min, min(value, max));
}

string getEmotion(const Json::Value& state) {
    int h = state["happiness"].asInt();
    if (h >= 80) return "😺";
    else if (h >= 50) return "😸";
    else if (h >= 30) return "😼";
    else return "😿";
}

string getMoodText(const Json::Value& state) {
    int h = state["happiness"].asInt();
    if (h >= 80) return "Мурлычет и трётся об ноги";
    else if (h >= 50) return "Играет с клубком ниток";
    else if (h >= 30) return "Сидит на подоконнике и смотрит в окно";
    else return "Лежит в углу и грустит";
}

void showStatus(Json::Value& state) {
    cout << colorize(CAT_ASCII, MAGENTA) << endl;
    cout << colorize("  Имя: " + state["name"].asString(), BLUE) << endl;
    cout << colorize("  Возраст: " + to_string(state["age"].asInt()) + " лет", BLUE) << endl;
    cout << "  " << getEmotion(state) << "  " << getMoodText(state) << endl;
    cout << colorize("  Голод: " + to_string(state["hunger"].asInt()) + "/100",
                     state["hunger"].asInt() > 70 ? YELLOW : GREEN) << endl;
    cout << colorize("  Счастье: " + to_string(state["happiness"].asInt()) + "/100",
                     state["happiness"].asInt() > 50 ? GREEN : RED) << endl;
    cout << colorize("  Энергия: " + to_string(state["energy"].asInt()) + "/100",
                     state["energy"].asInt() > 50 ? GREEN : YELLOW) << endl;
    cout << colorize("  Здоровье: " + to_string(state["health"].asInt()) + "/100",
                     state["health"].asInt() < 30 ? RED : GREEN) << endl;
}

void randomEvent(const string& msg = "") {
    const char* events[] = {"Мяу!", "Котик трётся о ноги.", "Принёс игрушку.",
                            "Пытается поймать муху.", "Свернулся клубком."};
    cout << colorize("  ✨ " + (msg.empty() ? events[rand() % 5] : msg), YELLOW) << endl;
}

int main(int argc, char* argv[]) {
    srand(time(nullptr));
    if (argc < 2) {
        cout << colorize("Usage: virtual_cat <status|feed|play|sleep|heal|rename|auto> [name]", YELLOW) << endl;
        return 1;
    }
    string action = argv[1];
    Json::Value state = loadState();

    if (action == "status") {
        showStatus(state);
    } else if (action == "feed") {
        state["hunger"] = clamp(state["hunger"].asInt() - 30);
        state["happiness"] = clamp(state["happiness"].asInt() + 10);
        state["health"] = clamp(state["health"].asInt() + 5);
        cout << colorize("🐟  Ням-ням! Котик поел.", GREEN) << endl;
        randomEvent("Мурлычет от удовольствия");
        saveState(state);
    } else if (action == "play") {
        if (state["energy"].asInt() < 20) {
            cout << colorize("😿  Котик слишком устал для игр.", RED) << endl;
            return 0;
        }
        state["happiness"] = clamp(state["happiness"].asInt() + 25);
        state["energy"] = clamp(state["energy"].asInt() - 20);
        state["hunger"] = clamp(state["hunger"].asInt() + 10);
        cout << colorize("🧶  Игра с клубком! Котик доволен.", GREEN) << endl;
        randomEvent("Прыгает за лазерной указкой");
        saveState(state);
    } else if (action == "sleep") {
        state["energy"] = clamp(state["energy"].asInt() + 40);
        state["hunger"] = clamp(state["hunger"].asInt() + 10);
        cout << colorize("😴  Котик уснул. Сладких снов!", BLUE) << endl;
        randomEvent("Во сне дёргает лапками");
        saveState(state);
    } else if (action == "heal") {
        if (state["energy"].asInt() < 20) {
            cout << colorize("😿  Нет сил лечить котика.", RED) << endl;
            return 0;
        }
        state["health"] = clamp(state["health"].asInt() + 30);
        state["energy"] = clamp(state["energy"].asInt() - 20);
        cout << colorize("💊  Котик вылечен! Он благодарен.", GREEN) << endl;
        saveState(state);
    } else if (action == "rename") {
        if (argc < 3) {
            cout << colorize("Укажите имя: rename <имя>", RED) << endl;
            return 0;
        }
        state["name"] = argv[2];
        cout << colorize("🐱  Котика теперь зовут " + state["name"].asString() + "!", BLUE) << endl;
        saveState(state);
    } else if (action == "auto") {
        cout << colorize("🤖  Автоматический режим включён.", MAGENTA) << endl;
        while (true) {
            state["hunger"] = clamp(state["hunger"].asInt() + (rand() % 21 - 5));
            state["happiness"] = clamp(state["happiness"].asInt() + (rand() % 16 - 5));
            state["energy"] = clamp(state["energy"].asInt() + (rand() % 16 - 5));
            state["health"] = clamp(state["health"].asInt() + (rand() % 9 - 3));
            saveState(state);
            system("clear");
            showStatus(state);
            cout << colorize("\nНажмите Ctrl+C для выхода", YELLOW) << endl;
            sleep(10);
        }
    } else {
        cout << colorize("Неизвестное действие.", RED) << endl;
    }
    return 0;
}
