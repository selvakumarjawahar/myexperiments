cmake_minimum_required(VERSION 3.15)
project(cppfun)

set(CMAKE_CXX_STANDARD 17)
find_package(Threads REQUIRED)
add_executable(shareptr
        main.cpp
)
target_link_libraries(shareptr
        PRIVATE Threads::Threads)