cmake_minimum_required(VERSION 3.10)
project(ConsoleGraph)

set(CMAKE_CXX_STANDARD 14)

find_package(Curses REQUIRED)
include_directories(${CURSES_INCLUDE_DIR})

add_executable(ConsoleGraph main.cpp)

target_link_libraries(ConsoleGraph ${CURSES_LIBRARIES})
