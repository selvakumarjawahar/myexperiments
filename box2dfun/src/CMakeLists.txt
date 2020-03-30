add_executable(box2dfun main.cpp)

target_include_directories(box2dfun
        PRIVATE C:/Users/selva/Projects/3rdparty/box2d/include
        PRIVATE C:/Users/selva/Projects/3rdparty/SFML/SFML-2.5.1/include
        )

target_link_directories(box2dfun
        PRIVATE C:/Users/selva/Projects/3rdparty/box2d/build/src/Debug
        PRIVATE C:/Users/selva/Projects/3rdparty/SFML/SFML-2.5.1/cmake-build-release/lib/Debug
        )

target_link_libraries(box2dfun
        PRIVATE box2d.lib
        PRIVATE sfml-audio-d
        PRIVATE sfml-graphics-d
        PRIVATE sfml-main-d
        PRIVATE sfml-network-d
        PRIVATE sfml-system-d
        PRIVATE sfml-window-d
        )