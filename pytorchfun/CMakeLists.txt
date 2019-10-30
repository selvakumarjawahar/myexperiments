cmake_minimum_required(VERSION 3.15)
project(pytorchfun)

find_package(Torch REQUIRED)

set(CMAKE_CXX_STANDARD 11)

add_executable(pytorchfun main.cpp)
target_link_libraries(pytorchfun "${TORCH_LIBRARIES}")