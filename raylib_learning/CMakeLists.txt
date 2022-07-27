cmake_minimum_required(VERSION 3.15)

project(raylib_learning
	VERSION 0.1
	LANGUAGES C CXX)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

set (SRCS
		examples/raylib_example.c)

find_package(PkgConfig REQUIRED)

pkg_check_modules(RAYLIB REQUIRED raylib)

add_executable(raylib_example
	${SRCS})

target_include_directories(raylib_example
		PRIVATE ${RAYLIB_INCLUDE_DIRS})

target_compile_options(raylib_example
		PRIVATE ${RAYLIB_CFLAGS})

target_link_libraries(raylib_example
		PRIVATE ${RAYLIB_LDFLAGS})


