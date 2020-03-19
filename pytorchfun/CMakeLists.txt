cmake_minimum_required(VERSION 3.15)
project(pytorchfun)

find_package(Torch REQUIRED)

set(CMAKE_CXX_STANDARD 11)

add_executable(pytorchfun main.cpp)
target_link_libraries(pytorchfun "${TORCH_LIBRARIES}")

if (MSVC)
    file(GLOB TORCH_DLLS "${TORCH_INSTALL_PREFIX}/lib/*.dll")
    add_custom_command(TARGET pytorchfun
            POST_BUILD
            COMMAND ${CMAKE_COMMAND} -E copy_if_different
            ${TORCH_DLLS}
            $<TARGET_FILE_DIR:pytorchfun>)
endif (MSVC)
