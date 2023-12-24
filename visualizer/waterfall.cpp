#include <ncurses.h>
#include <cmath>
#include <iostream>
#include <fstream>

const int num_colors = 18;

int* scan_flow(int size) {
    int* array = new int[size];

    for (int i = 0; i < size; ++i) {
        std::cin >> array[i];
    }

    return array;
}

bool is_end_of_flow(int a) {
    return a < 0;
}

void shift_and_insert(int** array2D, int* newValues, int M, int N) {

    delete(array2D[M-1]);

    for (int i = M - 1; i > 0; --i) {
        array2D[i] = array2D[i - 1];
    }

    array2D[0] = newValues;
}

void print_buffer(int* buffer[], int M, int N) {
    for (int i = 0; i < M; ++i) {
        for (int j = 0; j < N; ++j) {
            std::cout << buffer[i][j] << " ";
        }
        std::cout << "\n";
    }
}

//Оттенки синего: Темный синий: 18 | Средний синий: 19 | Светлый синий: 20 | Голубой: 23 | Светло-голубой: 39
// Зеленовато-голубой: 44 | Бирюзовый: 45 | Бирюзово-зеленый: 50 | Зеленовато-бирюзовый: 51

//Оттенки красного: Светлый красный: 196 | Оранжево-красный: 202 | Оранжевый: 208 | Темный оранжевый: 214
// Красновато-коричневый: 160 | Темно-красный: 124 | Красный: 160 | Темный красный: 196 | Красно-коричневый: 88

void init_color_pair() {

    int pair_number = 1;
    int blueShades[] = {18, 19, 20, 23, 39, 44, 45, 50, 51};
    int redShades[] = {196, 202, 208, 214, 160, 124, 160, 196, 88};

    for (const auto& bc: blueShades) {
        init_pair(pair_number, bc, bc);
        pair_number++;
    }

    for (const auto& rs: redShades) {
        init_pair(pair_number, rs, rs);
        pair_number++;
    }
}

int get_color_index(int value, int max_value) {
    int step = max_value / (num_colors - 1);

    int color_index = value / step + 1;

    if (color_index < 1) {
        color_index = 1;
    } else if (color_index > num_colors) {
        color_index = num_colors;
    }

    return color_index;
}

void draw_buffer(int* waterfall_buffer[], int waterfall_size, int flow_size, int max_amplitude) {
    clear();

    for (int i = 0; i < waterfall_size; ++i) {
        for (int j = 0; j < flow_size; ++j) {
            int value = waterfall_buffer[i][j];

            int colorIndex = get_color_index(value, max_amplitude);

            attron(COLOR_PAIR(colorIndex));
            mvprintw(i, j, "#");
            attroff(COLOR_PAIR(colorIndex));
        }
    }

    refresh();
}

void waterfall_buffer(int flow_size, int waterfall_size, int max_amplitude, int period, WINDOW* win) {
    int* waterfall_buffer[waterfall_size];
    for (int i = 1; i < waterfall_size; i++) {
        waterfall_buffer[i] = new int[flow_size]{};
    }

    std::fstream fifoStream("./visfifo", std::ios::in);

    waterfall_buffer[0] = scan_flow(flow_size);

    while (true) {
        int* flow = scan_flow(flow_size);
        if (is_end_of_flow(flow[0])) {
            break;
        }
        shift_and_insert(waterfall_buffer, flow, waterfall_size, flow_size);

        // real buffer
        // print_buffer(waterfall_buffer, waterfall_size, flow_size);

        // ncurses buffer
        draw_buffer(waterfall_buffer, waterfall_size, flow_size, max_amplitude);

        fifoStream.close();
        fifoStream.open("./visfifo");
    }

    endwin();

    for (int i = 0; i < waterfall_size; i++) {
        delete(waterfall_buffer[i]);
    }
}

void waterfall() {
    int flow_size, waterfall_size, max_amplitude, period;

    std::cout << "waterfall start" << std::endl << "write schema: " << std::endl;
    std::cout << "flow size | " << "waterfall size | " << "max amplitude | " << "period | " << std::endl;
    std::fstream fifoStream("./visfifo", std::ios::in);
    std::cin >> flow_size >> waterfall_size >> max_amplitude >> period;
    fifoStream.close();
    fifoStream.open("./visfifo");

    std::cout << "your schema:" << "\nfs: " << flow_size << "\nws: " << waterfall_size <<
    "\nma: " << max_amplitude << "\np: " << period << std::endl;

    initscr();
    start_color();

    init_color_pair();

    WINDOW* win = newwin(flow_size, waterfall_size, 0, 0);

    waterfall_buffer(flow_size, waterfall_size, max_amplitude, period, win);

    endwin();
}

