cmake_minimum_required(VERSION 3.18)
project(learn-gstreamer C CXX)

find_package(PkgConfig REQUIRED) 

pkg_check_modules(gstreamer  REQUIRED IMPORTED_TARGET gstreamer-1.0)

add_executable(learn-gstreamer
        src/main.cpp)


set_target_properties(learn-gstreamer PROPERTIES
        CXX_STANDARD 17
        CXX_STANDARD_REQUIRED YES
        CXX_EXTENSIONS YES
        )

target_link_libraries(learn-gstreamer
        PRIVATE PkgConfig::gstreamer
)
