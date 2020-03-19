cmake_minimum_required(VERSION 3.15)
project(cppfun)

set(CMAKE_CXX_STANDARD 11)

add_executable(cppfun main.cpp CommandGenerator.h CommandGenerator.cpp StaticMapGenerator.h)

target_include_directories(cppfun
        PRIVATE /home/selva/Projects/thirdparty/boost_1_62_0)