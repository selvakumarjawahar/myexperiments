add_executable(sfmlfun main.cpp)

target_include_directories(sfmlfun PRIVATE
        C:/Users/selva/Projects/3rdparty/SFML/SFML-2.5.1/include)

target_link_directories(sfmlfun PRIVATE
        C:/Users/selva/Projects/3rdparty/SFML/SFML-2.5.1/cmake-build-release/lib/Debug)

target_link_libraries(sfmlfun
        PRIVATE sfml-audio-d
        PRIVATE sfml-graphics-d
        PRIVATE sfml-main-d
        PRIVATE sfml-network-d
        PRIVATE sfml-system-d
        PRIVATE sfml-window-d
        )