#include <ncurses.h>
#include <cmath>
#include <iostream>

const int width = 128;
const int height = 30;

void drawGraph(const int data[], int dataSize) {
    initscr();
    start_color();
    init_pair(1, COLOR_GREEN, COLOR_BLACK);
    init_pair(2, COLOR_CYAN, COLOR_BLACK);

    clear();

    attron(COLOR_PAIR(1));
    for (int i = 0; i < dataSize - 1; ++i) {
        int x1 = i;
        int y1 = height - data[i] % height;
        int x2 = i + 1;
        int y2 = height - data[i + 1] % height;

        for (int j = x1; j <= x2; ++j) {
            int y = y1 + (y2 - y1) * (j - x1) / (x2 - x1);
            mvprintw(y, j, "*");
        }
    }
    attroff(COLOR_PAIR(1));

    attron(COLOR_PAIR(2));
    for (int i = 0; i < dataSize; ++i) {
        int x = i;
        int y = height - data[i] % height;
        mvprintw(y, x, "*");
    }
    attroff(COLOR_PAIR(2));

    refresh();
    getch();
    endwin();
}

int main() {
    int data[width];
    for (int i = 0; i < width; ++i) {
        data[i] = height / 2 + static_cast<int>(height / 2 * sin(i * 0.1));
    }

    std::cout << data << std::endl;

    drawGraph(data, width);

    return 0;
}

